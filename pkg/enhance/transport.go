//
// Copyright (c) Jeff Mendoza <jlm@jlm.name>
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// SPDX-License-Identifier: MIT
//

package enhance

import (
	"net/http"

	"golang.org/x/time/rate"
)

type transport struct {
	Wrapped http.RoundTripper
	RL      *rate.Limiter
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.RL.Wait(r.Context()); err != nil {
		return nil, err
	}
	return t.Wrapped.RoundTrip(r)
}
