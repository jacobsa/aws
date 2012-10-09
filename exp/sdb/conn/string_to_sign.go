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
	"strings"
)

// Given a request, return the canonicalized string that should be signed.
//
// Reference:
//     http://goo.gl/sRr8w
func computeStringToSign(req Request, host string) (string, error) {
	// We always use HTTP POST.
	verb := "POST"

	// The host header must be made lower-case.
	host = strings.ToLower(host)

	// The request URI is always /.
	requestUri := "/"

	// Canonicalize the query string, in this case the POST body.
	queryString := assemblePostBody(req)

	parts := []string{verb, host, requestUri, queryString}
	return strings.Join(parts, "\n"), nil
}
