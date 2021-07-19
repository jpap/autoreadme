// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var templateString = `<!-- DO NOT EDIT. -->
<!-- Automatically generated with https://go.jpap.org/godoc-readme-gen -->

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

func getTemplate(dir string) (*template.Template, error) {
	if len(*flagTemplate) == 0 {
		// Use the built-in template
		return builtinTemplate, nil
	}

	path := *flagTemplate
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		// File does not exist.  If it's the default name, use the built-in
		// template, otherwise return an error.
		if *flagTemplate == defaultTemplateFile {
			return builtinTemplate, nil
		}
		return nil, fmt.Errorf("failed to open template file: %s", *flagTemplate)
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return template.New("README").Parse(string(bs))
}
