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

package http_test

import (
	"github.com/jacobsa/aws/s3/http"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	sys_http "net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestConn(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type localHandler struct {
	// To be returned.
	statusCode int
	body []byte
}

func (h *localHandler) ServeHTTP(w sys_http.ResponseWriter, r *sys_http.Request) {
	w.WriteHeader(h.statusCode)
	if _, err := w.Write(h.body); err != nil {
		panic(err)
	}
}

type ConnTest struct {
	handler localHandler
	server *httptest.Server
	endpoint *url.URL
}

func init() { RegisterTestSuite(&ConnTest{}) }

func (t *ConnTest) SetUp(i *TestInfo) {
	t.server = httptest.NewServer(&t.handler)

	var err error
	t.endpoint, err = url.Parse(t.server.URL)
	AssertEq(nil, err)
}

func (t *ConnTest) TearDown() {
	t.server.Close()
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *ConnTest) InvalidScheme() {
	// Connection
	_, err := http.NewConn(&url.URL{Scheme: "taco", Host: "localhost"})

	ExpectThat(err, Error(HasSubstr("scheme")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ConnTest) UnknownHost() {
	// Connection
	conn, err := http.NewConn(&url.URL{Scheme: "http", Host: "foo.sidofhdksjhf"})
	AssertEq(nil, err)

	// Request
	req := &http.Request{
		Verb: "GET",
		Path: "/foo",
		Headers: map[string]string{},
	}

	// Call
	_, err = conn.SendRequest(req)

	ExpectThat(err, Error(HasSubstr("foo.sidofhdksjhf")))
	ExpectThat(err, Error(HasSubstr("no such host")))
}

func (t *ConnTest) InvalidVerb() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) PassesOnRequestInfo() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ReturnsStatusCode() {
	// Handler
	t.handler.statusCode = 123

	// Connection
	conn, err := http.NewConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := &http.Request{
		Verb: "GET",
		Path: "/",
		Headers: map[string]string{},
	}

	// Call
	resp, err := conn.SendRequest(req)
	AssertEq(nil, err)

	ExpectEq(123, resp.StatusCode)
}

func (t *ConnTest) ReturnsBody() {
	// Handler
	t.handler.body = []byte{0xde, 0xad, 0x00, 0xbe, 0xef}

	// Connection
	conn, err := http.NewConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := &http.Request{
		Verb: "GET",
		Path: "/",
		Headers: map[string]string{},
	}

	// Call
	resp, err := conn.SendRequest(req)
	AssertEq(nil, err)

	ExpectThat(resp.Body, DeepEquals(t.handler.body))
}

func (t *ConnTest) ServerReturnsEmptyBody() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) HttpsWorksProperly() {
	ExpectEq("TODO", "")
}
