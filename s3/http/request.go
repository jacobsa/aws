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

// An HTTP request to S3.
type Request struct {
	// The HTTP verb; e.g. "PUT" or "GET".
	Verb string

	// The path of the HTTP request URI, including the bucket name.
	//
	// For example:
	//     /mybucket/foo/bar/baz.jpg
	//
	Path string

	// HTTP headers to be included in the request.
	Headers map[string]string
}
