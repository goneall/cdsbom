//
// Copyright (c) Jeff Mendoza <jlm@jlm.name>
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
//

// Package enhance enhances sbom documents with ClearlyDefined license
// information
package enhance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/guacsec/sw-id-core/coordinates"
	"github.com/protobom/protobom/pkg/sbom"

	"github.com/jeffmendoza/cdsbom/pkg/cd"
)

// Do modifies the License and LicenseConcluded fields of the Nodes in the
// provided protobom Document with results from ClearlyDefined Warnings and
// updates are printed to stdout. TODO: Update to use a provided io.Writer or
// logger, also to use provided http client/transport and context.
func Do(s *sbom.Document) error {
	coords := coordList(s)
	defs, err := getDefs(coords)
	if err != nil {
		return err
	}
	updateLicenses(s, defs)
	return nil
}

// CoordList takes an SBOM document and returns a slice of all ClearlyDefined
// Coordinates found in that document.
func coordList(s *sbom.Document) []string {
	nodes := s.GetNodeList().GetNodes()
	coords := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if p := node.GetIdentifiers()[int32(sbom.SoftwareIdentifierType_PURL)]; p != "" {
			if c, err := coordinates.ConvertPurlToCoordinate(p); err == nil {
				coords = append(coords, c.ToString())
			} else {
				fmt.Printf("Coordinate conversion not supported for: %q\n", p)
			}
		}
	}
	return coords
}

func getDefs(coords []string) (map[string]*cd.Definition, error) {
	allDefs := make(map[string]*cd.Definition)
	chunkSize := 100
	for i := 0; i < len(coords); i += chunkSize {
		end := i + chunkSize
		if end > len(coords) {
			end = len(coords)
		}
		defs, err := getDefsFromService(coords[i:end])
		if err != nil {
			return nil, err
		}
		for k, v := range defs {
			allDefs[k] = v
		}
	}
	return allDefs, nil
}

func getDefsFromService(coords []string) (map[string]*cd.Definition, error) {
	cs, err := json.Marshal(coords)
	if err != nil {
		return nil, fmt.Errorf("error marshalling coordinates: %w", err)
	}
	rsp, err := http.Post("https://api.clearlydefined.io/definitions", "application/json", bytes.NewBuffer(cs))
	if err != nil {
		return nil, fmt.Errorf("error querying ClearlyDefined: %w", err)
	}
	if rsp.StatusCode != http.StatusOK {
		fmt.Println(string(cs))
		return nil, fmt.Errorf("error querying ClearlyDefined: %v", rsp.Status)
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var defs map[string]*cd.Definition
	if err := json.Unmarshal(body, &defs); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return defs, nil
}

func updateLicenses(s *sbom.Document, defs map[string]*cd.Definition) {
	for _, node := range s.GetNodeList().GetNodes() {
		updateNode(node, defs)
	}
}

func updateNode(n *sbom.Node, defs map[string]*cd.Definition) {
	p := n.GetIdentifiers()[int32(sbom.SoftwareIdentifierType_PURL)]
	if p == "" {
		return
	}
	c, err := coordinates.ConvertPurlToCoordinate(p)
	if err != nil {
		return
	}
	d, ok := defs[c.ToString()]
	if !ok {
		return
	}
	if len(d.Described.Tools) == 0 {
		return
	}
	old := strings.Join(n.GetLicenses(), " AND ")
	new := d.Licensed.Declared
	if old != new {
		fmt.Printf("Update Declared License\n")
		fmt.Printf("Name: %v\tVersion: %v\n", n.GetName(), n.GetVersion())
		fmt.Printf("\t\t\t\tSBOM License: %q\tCD License: %q\n", old, new)
		n.Licenses = []string{new}
	}

	oldDisc := n.GetLicenseConcluded()
	newDisc := strings.Join(d.Licensed.Facets.Core.Discovered.Expressions, " AND ")
	if oldDisc != newDisc {
		fmt.Printf("Update Discovered License\n")
		fmt.Printf("Name: %v\tVersion: %v\n", n.GetName(), n.GetVersion())
		fmt.Printf("\t\t\t\tSBOM License: %q\tCD License: %q\n", oldDisc, newDisc)
		n.LicenseConcluded = newDisc
	}
}
