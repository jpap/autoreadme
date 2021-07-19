// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"
)

var gopaths = build.Default.SrcDirs()

type Doc struct {
	Name       string
	ImportPath string
	Synopsis   string
	Doc        string
	Title      string
	RepoPath   string
	IsLibrary  bool
	Bugs       []string
	Commands   []string // import path of any cmd/* subpackages
	HasTravis  bool     // true when a `.travis.yml` file is in the package dir
	Examples   map[string]Example
}

type Example struct {
	Name   string
	Code   string
	Output string // the expected output, if not ""
}

// Map returns the receiver as a map for use with a template.
func (d Doc) Map() map[string]interface{} {
	return map[string]interface{}{
		"Name":       d.Name,
		"ImportPath": d.ImportPath,
		"Synopsis":   d.Synopsis,
		"Doc":        d.Doc,
		"Today":      time.Now().Format("2006.01.02"),
		"Title":      d.Title,
		"RepoPath":   d.RepoPath,
		"Bugs":       d.Bugs,
		"Library":    d.IsLibrary,
		"Commands":   d.Commands,
		"Travis":     d.HasTravis,
		"Examples":   d.Examples,
	}
}

func importPath(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for _, p := range gopaths {
		p, err = filepath.Rel(p, abs)
		if err != nil {
			return "", err
		}
		//not a relative path, therefore the correct one
		if len(p) > 0 && p[0] != '.' {
			return filepath.ToSlash(p), nil
		}
	}
	return "", fmt.Errorf("Not in go root or element of $GOPATH: %s", dir)
}

func NewDoc(dir string) (d Doc, err error) {
	var pkg *packages.Package
	pkg, err = loadPackage(dir,
		packages.NeedName|
			packages.NeedFiles|
			packages.NeedSyntax|
			packages.NeedTypes|
			packages.NeedDeps,
	)
	if err != nil {
		return
	}

	d.ImportPath = pkg.PkgPath

	// Manually construct the *ast.Package... ast.NewPackage gets a bit ambitious
	// with resolving identifiers that it just doesn't work with the pre-parsed
	// []*ast.File provided by x/tools/go/packages.
	astPkg := &ast.Package{
		Name:  pkg.Name,
		Files: make(map[string]*ast.File),
	}
	for i, f := range pkg.Syntax {
		// Filenames in *packages.Package are a pain... and they don't matter here,
		// so let's make them up. ;-)
		filename := fmt.Sprintf("file-%d.go", i)
		astPkg.Files[filename] = f
	}

	// Parse package docs
	docPkg := doc.New(astPkg, pkg.PkgPath, 0)
	d.Doc = packageDocString(docPkg)
	d.Synopsis = doc.Synopsis(docPkg.Doc)

	// Render examples
	d.Examples = make(map[string]Example)
	for _, f := range pkg.Syntax {
		for _, ex := range doc.Examples(f) {
			d.Examples[ex.Name] = renderExample(ex)
		}
	}

	// Render bugs
	for _, bug := range docPkg.Notes["BUG"] {
		d.Bugs = append(d.Bugs, bug.Body)
	}

	name := pkg.Name
	if name == "main" {
		// main package: get the package name from the import path
		var dir string
		dir, err = goPackagesDir(pkg)
		if err != nil {
			return
		}
		name = filepath.Base(dir)
	}

	// Strip the first path component of the import path to derive the repo path.
	pathelms := strings.Split(pkg.PkgPath, "/")[1:]
	d.RepoPath = path.Join(pathelms...)

	d.Name = name
	d.Title = name
	if len(*Title) > 0 {
		d.Title = *Title
	}

	if pkg.Name == "main" {
		d.Commands = append(d.Commands, pkg.PkgPath)
	} else {
		d.IsLibrary = true
	}

	// Look for additional cmd/* subdirs in the package, that correspond to main packages.
	if dirs, err := filepath.Glob(filepath.Join(dir, "cmd/*")); err == nil {
		for _, dir := range dirs {
			cmdPkg, err := loadPackage(dir, packages.NeedName)
			if err != nil {
				continue // ignore it
			}
			if cmdPkg.Name == "main" {
				d.Commands = append(d.Commands, cmdPkg.PkgPath)
			}
		}
	}

	d.HasTravis = hasTravisConfig(dir)

	return
}

func hasTravisConfig(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, ".travis.yml")); err == nil {
		return true
	}
	return false
}
