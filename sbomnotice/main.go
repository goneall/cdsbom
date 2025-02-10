//
// Copyright (c) Jeff Mendoza <jlm@jlm.name>
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// SPDX-License-Identifier: MIT
//

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/protobom/protobom/pkg/reader"
	"github.com/protobom/protobom/pkg/sbom"

	"github.com/jeffmendoza/cdsbom/pkg/enhance"
)

func main() {
	inFile, outFile := flags()

	document := read(inFile)

	notice, err := enhance.Notice(context.Background(), document)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	write(notice, outFile)
	fmt.Println("Complete")
}

// flags sets up and parses flags. Return values are input file and output file
// respecively.
func flags() (string, string) {
	o := flag.String("out", "NOTICE", "Name of output file")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\tThis program generates a NOTICE file from an SBOM\n")
		fmt.Printf("%s [options] <in-SBOM-file>\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	r := flag.Args()
	if len(r) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	i := r[0]

	return i, *o
}

// read reads in the sbom document and also returns the format.
func read(i string) *sbom.Document {
	reader := reader.New()
	d, err := reader.ParseFile(i)
	if err != nil {
		fmt.Printf("Error reading input SBOM: %v\n", err)
		os.Exit(1)
	}
	return d
}

// write writes the sbom document to a file with a format
func write(notice string, of string) {
	err := os.WriteFile(of, []byte(notice), 0666)
	if err != nil {
		fmt.Printf("Error writing outout: %v\n", err)
		os.Exit(1)
	}
}
