//
// Copyright (c) Jeff Mendoza <jlm@jlm.name>
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// SPDX-License-Identifier: MIT
//

package enhance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/protobom/protobom/pkg/sbom"
)

type NoticeReq struct {
	Coordinates []string `json:"coordinates"`
}

type NoticeRsp struct {
	Content string `json:"content"`
	Summary struct {
		Total    int `json:"total"`
		Warnings struct {
			NoDefinition []string `json:"noDefinition"`
			NoLicense    []string `json:"noLicense"`
			NoCopyright  []string `json:"noCopyright"`
		} `json:"warnings"`
	} `json:"summary"`
}

// Notice takes an SBOM document and queries ClearlyDefined for a NOTICE file
// for all the recognized components in the SBOM.
func Notice(ctx context.Context, s *sbom.Document) (string, error) {
	c := coordList(s)
	return request(ctx, c)
}

// request gets the NOTICE file for the coords from ClearlyDefined
// TODO: generalize chunking from getDefsto work with that and this.
func request(ctx context.Context, coords []string) (string, error) {
	cs, err := json.Marshal(NoticeReq{coords})
	if err != nil {
		return "", fmt.Errorf("error marshaling coordinates: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.clearlydefined.io/notices", bytes.NewBuffer(cs))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting NOTICE file: %w", err)
	}
	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting NOTICE file: %v", rsp.Status)
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	var nr NoticeRsp
	if err := json.Unmarshal(body, &nr); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}
	if len(nr.Summary.Warnings.NoDefinition) > 0 {
		fmt.Printf("Warning, no definition: %v\n", nr.Summary.Warnings.NoDefinition)
	}
	if len(nr.Summary.Warnings.NoLicense) > 0 {
		fmt.Printf("Warning, no license: %v\n", nr.Summary.Warnings.NoLicense)
	}
	if len(nr.Summary.Warnings.NoCopyright) > 0 {
		fmt.Printf("Warning, no copyright: %v\n", nr.Summary.Warnings.NoCopyright)
	}
	return nr.Content, nil
}
