// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// goPackagesDir returns the directory of pkg.
//
// It relies on the fact that GoFiles are absolute paths.
// See https://github.com/golang/go/issues/38445#issuecomment-732806052
func goPackagesDir(pkg *packages.Package) (string, error) {
	if len(pkg.GoFiles) == 0 {
		return "", errors.New("package has no Go files")
	}
	return filepath.Dir(pkg.GoFiles[0]), nil
}

// loadPackage loads the package in the given filesystem directory.
func loadPackage(dir string, m packages.LoadMode) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: m,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load package in dir %q: %w", dir, err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		// Package failed to parse
		os.Exit(1)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("could not find package: %s", dir)
	}
	return pkgs[0], nil
}
