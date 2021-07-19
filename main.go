// Copyright 2013 James Frasche. All rights reserved.
// Copyright 2021 John Papandriopoulos.
// Use of this code is governed by a BSD-License found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const defaultTemplateFile = ".README.template.md"

var (
	flagForce         = flag.Bool("f", false, "Run even if README.md exists, overwriting original")
	flagPrintTemplate = flag.Bool("print-template", false, "Print the built in template to stdout and exit")
	flagTemplate      = flag.String("template", defaultTemplateFile, "Template to use, or builtin if does not exist")
	flagTitle         = flag.String("title", "", "Title of the README.md")
	flagDefs          defFlag
)

func getOrCreateReadmeFile(dir string) (*os.File, error) {
	nm := filepath.Join(dir, "README.md")
	if !*flagForce {
		_, err := os.Stat(nm)
		if err == nil {
			return nil, fmt.Errorf("README.md already exists at %s. Use -f to overwrite", dir)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return os.Create(nm)
}

func main() {
	log.SetFlags(0)

	flag.Var(&flagDefs, "def", "Template define having the form: name=value")
	flag.Usage = func() {
		log.Printf("Usage of %s: %s [directory]", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagPrintTemplate {
		fmt.Print(templateString)
		return
	}

	dir := "."

	if args := flag.Args(); len(args) > 1 {
		log.Fatalln("Too many arguments: ", args[1:])
		flag.Usage()
	} else if len(args) == 1 {
		dir = args[0]
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := NewDoc(dir)
	if err != nil {
		log.Fatalf("Failed to load package in dir %q: %v\n", dir, err)
	}

	f, err := getOrCreateReadmeFile(dir)
	if err != nil {
		log.Fatalf("Failed to create README.md for %q: %v\n", dir, err)
	}

	// Convert the doc to a map, so we can add additional fields
	docm := doc.Map()
	for _, d := range flagDefs {
		// Lowercase define
		docm[d.Name] = d.Value
		// Uppercase define
		d = d.UpperClone()
		docm[d.Name] = d.Value
	}

	// Execute the template
	tmpl, err := getTemplate(dir)
	if err != nil {
		log.Fatalf("Failed to load template: %v\n", err)
	}

	if err = tmpl.Execute(f, docm); err != nil {
		log.Fatalf("Failed to execute template: %v\n", err)
	}
}
