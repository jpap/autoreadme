// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"strings"
	"text/template"
)

var templateString = `
# {{.Title}}
{{- if .Library}} [![GoDoc](https://pkg.go.dev/badge/{{.Import}}.svg)](https://pkg.go.dev/{{.Import}}){{end}}
{{- if .Travis}} [![Build Status](https://travis-ci.org/{{.RepoPath}}.png?branch=master)](https://travis-ci.org/{{.RepoPath}}){{end}}

{{if .Command -}}
# Install

$CODEBLOCKshell
go install "{{.Import}}"
$CODEBLOCK
{{end -}}

{{if .Library -}}
# Import

$CODEBLOCKgo
import "{{.Import}}"
$CODEBLOCK
{{end -}}

# Overview

{{.Doc}}

{{if .Bugs -}}
# Bugs

{{range .Bugs}}* {{.}}{{end}}
{{end}}
`

var tmpl *template.Template

func init() {
	// Backticks aren't allowed in a string literal...
	templateString = strings.ReplaceAll(templateString, "$CODEBLOCK", "```")
	tmpl = template.Must(template.New("").Parse(templateString))
}
