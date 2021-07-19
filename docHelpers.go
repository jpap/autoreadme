// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/format"
	"go/token"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/lexers"
)

// A par is the start-end line numbers of a paragraph.
type par struct {
	Index int
	Len   int
}

// A section is a list of paragraphs.
type section []par

var regexpNumberedItem = regexp.MustCompile(`^[0-9]+\.`)

func packageDocString(pkg *doc.Package) string {
	var lines []string
	push := func(ss ...string) {
		lines = append(lines, ss...)
	}
	nl := func() {
		push("\n")
	}

	// Map all of the paragraph lines, so that we can detect bullets and lists.
	var sections []section  // All sections
	var currSection section // The "current" section

	endPar := func() {
		sections = append(sections, currSection)
		currSection = make(section, 0)
	}

	for _, b := range blocks(pkg.Doc) {
		ls := b.lines
		switch b.op {
		case opPara:
			p := par{len(lines), len(ls)}
			currSection = append(currSection, p)

			push(ls...)
			nl()
		case opHead:
			endPar()
			push("## ")
			push(ls...)
			nl()
		case opPre:
			// Detect language... and hope the lexer name is the same as the lingust
			// name.  GitHub Markdown uses linguist:
			// https://github.com/github/linguist/blob/master/lib/linguist/languages.yml
			codeBlock := strings.Join(ls, "\n") + "\n"

			// It turns out using alecthomas/chroma uses a far-too-basic heuristic for
			// detecting Go, and so we add a more comprehensive heuristic here
			// instead.  Given that we're mostly going to find Go code here, we lean
			// towards detecting it over not.
			//
			// We then defer to the generic analyzer following.
			if strings.Contains(codeBlock, "package ") ||
				strings.Contains(codeBlock, "func(") ||
				strings.Contains(codeBlock, " := ") ||
				strings.Contains(codeBlock, "fmt.") ||
				strings.Contains(codeBlock, "var ") {
				push("```go")
			} else if lexer := lexers.Analyse(codeBlock); lexer != nil {
				push("```" + strings.ToLower(lexer.Config().Name))
			} else {
				push("```")
			}

			nl()
			push(ls...)
			push("```")
			nl()
			nl()
		}
	}
	endPar()

	// Detect and indent paragraphs within bullets/lists.  We do not support
	// nested lists, because the GoDoc format is too ambiguous for that use-case,
	// Unless numbering is given as 1.1, 1.1.1, 1.1.2, etc., which is quite
	// complex, and non-standard for bulleted lists: *, *.*, *.*, *.*.*, etc.
	//
	// If we detect a "^[0-9]+\." or "* " paragraph prefix, then assume all of the
	// subsequent paragraphs in the section are part of the same bullet/list,
	// until we find another bullet/list item.
	//
	// The only ambiguity here is a regular paragraph following a bullet/list
	// block, so consider a non-standard "section separator" paragraph of "...",
	// which we will elide.
	for _, sec := range sections {
		indent := false
		for _, par := range sec {
			leadin := lines[par.Index]

			switch {
			case par.Len == 1 && leadin == "...\n":
				// Found a section separator: remove it.
				indent = false
				lines[par.Index] = "\n"

			case regexpNumberedItem.MatchString(leadin):
				// Found a numbered list: indent subsequent pars.
				indent = true

			case strings.HasPrefix(leadin, "* "):
				// Found a NEW bulleted list: indent subsequent pars.
				indent = true

			default:
				if indent {
					for i := 0; i < par.Len; i++ {
						lines[par.Index+i] = "    " + lines[par.Index+i]
					}
				}
			}

		}
	}

	return strings.Join(lines, "")
}

func renderExample(ex *doc.Example) Example {
	e := Example{
		Name: strings.Replace(ex.Name, "_", " ", -1),
	}

	c := &bytes.Buffer{}
	format.Node(c, token.NewFileSet(), ex.Code)
	e.Code = fmt.Sprintf("Code:\n\n```go\n%s\n```\n", c.String())

	if ex.Output != "" {
		e.Output = fmt.Sprintf("Output:\n\n```\n%s\n```\n", ex.Output)
	}

	return e
}
