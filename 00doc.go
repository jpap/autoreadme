// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

// Automatically generate a README.md for your Go project.
//
// This tool creates a github-formatted README.md using the same format as
// godoc.  It includes the package summary and generates badges for pkg.go.dev
// and travis-ci.
//
// This is a fork of James Frasche's project found at
// https://github.com/jimmyfrasche/autoreadme.
//
// Heuristics
//
// Go code in the current directory is analyzed by default, using the template
// file `.README.template.md` in the same directory, or a default template
// otherwise.  The Markdown output is always written to `README.md`, which is
// never overwritten unless the `-f` flag is given.
//
// A specific directory can be provided on the command line, and a custom
// template specified using the `-template` flag.
//
// Multiple directories can be recursed automatically using the `-r` flag.
//
// Examples
//
// Create a README.md for the directory a/b/c
//  godoc-readme-gen a/b/c
//
// Overwrite the README.md in the current directory
//  godoc-readme-gen -f
//
// Generate for the current directory and all subdirectories containing Go code
//  godoc-readme-gen -r
//
// Copy the built-in template to a file for the creation of a new template.
//  godoc-readme-gen -print-template > .README.template.md
//
// Generate using a custom template.
//  godoc-readme-gen -template path/to/readme.template
//
// Template Variables
//
// The following variables are available in custom templates:
//
// `.Name` Package name.
//
// `.Title` The -title flag value, or package name if not provided.
//
// `.Doc` Package-level documentation.
//
// `.Synopsis` The first sentence from the .Doc variable.
//
// `.ImportPath` Package import path.
//
// `.RepoPath` The import path without the first path component. For example,
// the import github.com/golang/go is represented as "golang/go".  This is
// typically the path within the repo of the package.
//
// `.Bugs` A []string of all bugs as per godoc.
//
// `.Commands` A []string of import paths of all main packages.  In addition to
// the directory provided to the tool, we also check cmd/* directories for
// additional main packages.
//
// `.Library` True if the package is not a main package.
//
// `.Today` The current date in YYYY.MM.DD format.
//
// `.Travis` True if there is a `.travis.yml` file in the package directory.
//
// `.Examples` a map of Example with all examples from `*_test.go` files. These
// can be used to include selective examples into the README.  The Example
// struct has the following fields:
//   .Name    Name of the example
//   .Code    Rendered example code similar to godoc
//   .Output  Example output, if any
//
package main

// To install: `go install go.jpap.org/godoc-readme-gen`
//
//go:generate godoc-readme-gen -f -title "GoDoc README Markdown Generator"
