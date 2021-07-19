// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Force         = flag.Bool("f", false, "Run even if README.md exists, overwriting original")
	Recursive     = flag.Bool("r", false, "Run in all subdirectories containing Go code")
	PrintTemplate = flag.Bool("print-template", false, "write the built in template to stdout and exit")
	Template      = flag.String("template", "", "specify a file to use as template, overrides built in template and .README.template.md")
	Title         = flag.String("title", "", "title of the README.md")
	Defs          defFlag
)

func getOrCreateReadmeFile(dir string) (*os.File, error) {
	nm := filepath.Join(dir, "README.md")
	if !*Force {
		_, err := os.Stat(nm)
		if err == nil {
			return nil, fmt.Errorf("README.md already exists at %s. Use -f to overwrite", dir)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return os.Create(nm)
}

// goDirList returns a list of directories that contain .go files from the
// rootPath directory.
func goDirList(rootPath string) (ds []string, err error) {
	if !*Recursive {
		// Do not recurse...
		return []string{rootPath}, nil
	}

	var follow func(string) error
	follow = func(dir string) error {
		fis, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		hasgo := false
		for _, fi := range fis {
			nm := fi.Name()
			if fi.IsDir() {
				if nm[0] != '.' {
					if err := follow(nm); err != nil {
						return err
					}
				}
			} else if strings.HasSuffix(nm, ".go") {
				hasgo = true
			}
		}
		if hasgo {
			ds = append(ds, dir)
		}
		return nil
	}

	if err = follow(rootPath); err != nil {
		return nil, err
	}

	return ds, nil
}

func main() {
	log.SetFlags(0)

	flag.Var(&Defs, "def", "Template define having the form: name=value")
	flag.Usage = func() {
		log.Printf("Usage of %s: %s [directory]", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *PrintTemplate {
		fmt.Print(templateString)
		return
	}

	rootPath := "."

	if args := flag.Args(); len(args) > 1 {
		log.Fatalln("Too many arguments: ", args[1:])
		flag.Usage()
	} else if len(args) == 1 {
		rootPath = args[0]
	}

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Fatalln(err)
	}

	dirList, err := goDirList(rootPath)
	if err != nil {
		log.Fatalln(err)
	}

	issuedWarning := false
	emitWarning := func(dir string, err error) {
		issuedWarning = true
		log.Printf("Could not create README.md for %q: %v\n", dir, err)
	}

	for _, dir := range dirList {
		tmpl, err := getTemplate(dir)
		if err != nil {
			if *Template != "" && *Recursive {
				//don't want to show same error message for all dirs so just bail now
				log.Fatalln(err)
			}
			emitWarning(dir, err)
			continue
		}

		doc, err := NewDoc(dir)
		if err != nil {
			emitWarning(dir, err)
			continue
		}

		f, err := getOrCreateReadmeFile(dir)
		if err != nil {
			emitWarning(dir, err)
			continue
		}

		// Convert the doc to a map, so we can add additional fields
		docm := doc.Map()
		for _, d := range Defs {
			// Lowercase define
			docm[d.Name] = d.Value
			// Uppercase define
			d = d.UpperClone()
			docm[d.Name] = d.Value
		}

		// Execute the template
		if err = tmpl.Execute(f, docm); err != nil {
			emitWarning(dir, err)
		}
	}

	if issuedWarning {
		os.Exit(1)
	}
}
