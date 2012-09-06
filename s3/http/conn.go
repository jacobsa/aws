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

// A connection to a particular server over a particular protocol (HTTP or
// HTTPS).
type Conn interface {
	// Call the server with the supplied request, returning a response if and
	// only if a response was received from the server. (That is, a 500 error
	// from the server will be returned here as a response with a nil error).
	SendRequest(r *Request) (*Response, error)
}
