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
	"net/url"
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
func NewHttpConn(endpoint *url.URL) (HttpConn, error)
