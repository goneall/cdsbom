//
// Copyright (c) Jeff Mendoza <jlm@jlm.name>
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// SPDX-License-Identifier: MIT
//

package main

import (
	"encoding/json"
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

	coords := enhance.CoordList(document)
	bts, err := json.Marshal(coords)
	if err != nil {
		fmt.Printf("Error marshaling coordinates: %v\n", err)
		os.Exit(1)
	}

	write(bts, outFile)
	fmt.Println("Complete")
}

// flags sets up and parses flags. Return values are input file and output file
// respecively.
func flags() (string, string) {
	o := flag.String("out", "coords.json", "Name of output file")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\tThis program generates a list of ClearlyDefined Coordinates from an SBOM\n")
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

// write writes the document to a file
func write(data []byte, of string) {
	err := os.WriteFile(of, data, 0666)
	if err != nil {
		fmt.Printf("Error writing outout: %v\n", err)
		os.Exit(1)
	}
}
