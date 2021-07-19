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
	"strings"
)

func packageDocString(pkg *doc.Package) string {
	var acc []string
	push := func(ss ...string) {
		acc = append(acc, ss...)
	}
	nl := func() {
		push("\n")
	}
	for _, b := range blocks(pkg.Doc) {
		ls := b.lines
		switch b.op {
		case opPara:
			push(ls...)
			nl()
		case opHead:
			push("## ")
			push(ls...)
			nl()
		case opPre:
			push("```")
			nl()
			push(ls...)
			push("```")
			nl()
			nl()
		}
	}
	return strings.Join(acc, "")
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
