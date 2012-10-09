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

package conn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpResponse struct {
	// The HTTP status code, e.g. 200 or 404.
	StatusCode int

	// The response body. This is the empty slice if the body was empty.
	Body []byte
}

// A connection to a host over HTTP.
type HttpConn interface {
	// Send the supplied request to the service.
	SendRequest(req Request) (resp *HttpResponse, err error)
}

// Return a connection to the supplied endpoint, based on its scheme and host
// fields.
func NewHttpConn(endpoint *url.URL) (HttpConn, error) {
	switch endpoint.Scheme {
	case "http", "https":
	default:
		return nil, fmt.Errorf("Unsupported scheme: %s", endpoint.Scheme)
	}

	return &httpConn{endpoint}, nil
}

type httpConn struct {
	endpoint *url.URL
}

func (c *httpConn) SendRequest(req Request) (resp *HttpResponse, err error) {
	// Create an appropriate URL.
	u := url.URL{
		Scheme: c.endpoint.Scheme,
		Host:   c.endpoint.Host,
		Path:   "/",
	}

	urlStr := u.String()

	// Create an appropriate body.
	bodyVals := url.Values{}
	for key, val := range req {
		bodyVals.Set(key, val)
	}

	body := bodyVals.Encode()

	// Exception: Amazon's escaping rules disagree with Go about how to
	// URL-encode a space; they require %20 rather than '+'. Fix this up, noting
	// that actual plus characters were already percent-encoded.
	//
	// Reference:
	//     http://goo.gl/0aD5S
	body = strings.Replace(body, "+", "%20", -1)

	// Create a request to the system HTTP library.
	sysReq, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	// Set required headers.
	//
	// More info:
	//     http://goo.gl/0aD5S
	sysReq.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded; charset=utf-8")

	sysReq.Header.Set(
		"Host",
		u.Host)

	// Call the system HTTP library.
	sysResp, err := http.DefaultClient.Do(sysReq)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %v", err)
	}

	// Convert the response.
	resp = &HttpResponse{
		StatusCode: sysResp.StatusCode,
	}

	if resp.Body, err = ioutil.ReadAll(sysResp.Body); err != nil {
		return nil, fmt.Errorf("Reading body: %v", err)
	}

	sysResp.Body.Close()

	return resp, nil
}
