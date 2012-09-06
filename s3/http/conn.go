// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// A connection to a particular server over a particular protocol (HTTP or
// HTTPS).
type Conn interface {
	// Call the server with the supplied request, returning a response if and
	// only if a response was received from the server. (That is, a 500 error
	// from the server will be returned here as a response with a nil error).
	SendRequest(r *Request) (*Response, error)
}

func NewConn(host, scheme string) (Conn, error) {
	switch scheme {
	case "http", "https":
	default:
		return nil, fmt.Errorf("Unsupported scheme: %s", scheme)
	}

	return &conn{host, scheme}, nil
}

type conn struct {
	host string
	scheme string
}

func (c *conn) SendRequest(r *Request) (*Response, error) {
	// Create an appropriate URL.
	url := url.URL{
		Scheme: c.scheme,
		Host: c.host,
		Path: r.Path,
	}

	urlStr := url.String()

	// Call the system HTTP library.
	_, err := http.NewRequest(r.Verb, urlStr, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	return nil, fmt.Errorf("TODO")
}
