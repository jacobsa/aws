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
	"log"
	"net"
	"net/url"
	"syscall"
)

////////////////////////////////////////////////////////////////////////
// Public
////////////////////////////////////////////////////////////////////////

// Return a connection that wraps the supplied one, retrying a few times when
// it returns certain errors that S3 has been known to return transiently.
//
// Exposed only for testing; do not use directly. NewConn incorporates this
// functionality for you.
func NewRetryingConn(wrapped Conn) (c Conn, err error) {
	c = &retryingConn{wrapped}
	return
}

////////////////////////////////////////////////////////////////////////
// Implementation
////////////////////////////////////////////////////////////////////////

type retryingConn struct {
	wrapped Conn
}

func shouldRetry(err error) bool {
	// If there's no error, stop.
	if err == nil {
		return false
	}

	// Look for "broken pipe" errors. S3 seems to close keep-alive connections
	// that have been in use for awhile (on the order of 20-30 minutes). Perhaps
	// it's a machine being restarted on their end?
	if httpErr, ok := err.(*Error); ok {
		// EPIPE errors show up (or showed up at one time) as net.OpErrors.
		if opErr, ok := httpErr.OriginalErr.(*net.OpError); ok {
			if errno, ok := opErr.Err.(syscall.Errno); ok {
				if errno == syscall.EPIPE {
					log.Println("EPIPE; retrying.")
					return true
				}
			}
		}

		// Another class of errors that show up is url.Errors with the error string
		// "EOF.
		if urlErr, ok := httpErr.OriginalErr.(*url.Error); ok {
			if urlErr.Err.Error() == "EOF" {
					log.Println("EOF; retrying.")
					return true
			}
		}
	}

	return false
}

func (c *retryingConn) SendRequest(req *Request) (resp *Response, err error) {
	const maxTries = 3

	for i := 0; i < maxTries; i++ {
		resp, err = c.wrapped.SendRequest(req)
		if !shouldRetry(err) {
			break
		}
	}

	return
}
