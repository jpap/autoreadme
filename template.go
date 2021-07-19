// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var templateString = `
# {{.Title}}
{{- if .Library}} [![GoDoc](https://pkg.go.dev/badge/{{.ImportPath}}.svg)](https://pkg.go.dev/{{.ImportPath}}){{end}}
{{- if .Travis}} [![Build Status](https://travis-ci.org/{{.RepoPath}}.png?branch=master)](https://travis-ci.org/{{.RepoPath}}){{end}}

{{if .Commands -}}
# Install

$CODEBLOCKshell
{{range $cmdPath := .Commands -}}
go install {{$cmdPath}}
{{end -}}
$CODEBLOCK
{{end -}}

{{if .Library -}}
# Import

$CODEBLOCKgo
import "{{.ImportPath}}"
$CODEBLOCK
{{end -}}

# Overview

{{.Doc}}

{{if .Bugs -}}
# Bugs

{{range .Bugs}}* {{.}}{{end}}
{{end}}
`

var builtinTemplate *template.Template

func init() {
	// Backticks aren't allowed in a string literal...
	templateString = strings.ReplaceAll(templateString, "$CODEBLOCK", "```")
	builtinTemplate = template.Must(template.New("").Parse(templateString))
}

// only read and parse specified template once if -template and -r specified
var singletonTemplate *template.Template

func getTemplate(dir string) (*template.Template, error) {
	if *Template != "" {
		if singletonTemplate != nil {
			return singletonTemplate, nil
		}
		bs, err := ioutil.ReadFile(*Template)
		if err != nil {
			return nil, err
		}
		singletonTemplate, err = template.New("").Parse(string(bs))
		if err != nil {
			return nil, err
		}
		return singletonTemplate, nil
	}

	bs, err := ioutil.ReadFile(filepath.Join(dir, ".README.template.md"))
	//local template file found, prefer
	if err == nil {
		return template.New("").Parse(string(bs))
	}
	//the file was found but something else happened
	if !os.IsNotExist(err) {
		return nil, err
	}

	//the file was not found and no template specified, so use default
	return builtinTemplate, nil
}
