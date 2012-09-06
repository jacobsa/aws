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

package auth

import (
	"errors"
	"fmt"
	"github.com/jacobsa/aws/s3/http"
)

// Given an HTTP request, return the string that should be signed for that
// request. The request must include a `Date` header.
//
// See here for more info:
//
//     http://goo.gl/Z8DiC
//
func stringToSign(r *http.Request) (string, error) {
	// Grab the HTTP headers specifically called out by the signing algorithm.
	date, ok := r.Headers["Date"]
	if !ok {
		return "", errors.New("stringToSign requires a Date header.")
	}

	contentMd5 := r.Headers["Content-MD5"]
	contentType := r.Headers["Content-Type"]

	// We don't yet support CanonicalizedAmzHeaders.
	canonicalizedAmzHeaders := ""

	// We currently only support simple path-style requests.
	canonicalizedResource := r.Path

	// Put everything together.
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s%s",
		r.Verb,
		contentMd5,
		contentType,
		date,
		canonicalizedAmzHeaders,
		canonicalizedResource), nil
}
