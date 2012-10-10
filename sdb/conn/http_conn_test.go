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

package conn_test

import (
	"github.com/jacobsa/aws/exp/sdb/conn"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHttpConn(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type localHandler struct {
	// Input seen.
	req     *http.Request
	reqBody []byte

	// To be returned.
	statusCode int
	body       []byte
}

func (h *localHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	// Record the request.
	if h.req != nil {
		panic("Called twice.")
	}

	h.req = r

	// Record the body.
	if h.reqBody, err = ioutil.ReadAll(r.Body); err != nil {
		panic(err)
	}

	// Write out the response.
	w.WriteHeader(h.statusCode)
	if _, err := w.Write(h.body); err != nil {
		panic(err)
	}
}

type HttpConnTest struct {
	handler  localHandler
	server   *httptest.Server
	endpoint *url.URL
}

func init() { RegisterTestSuite(&HttpConnTest{}) }

func (t *HttpConnTest) SetUp(i *TestInfo) {
	t.server = httptest.NewServer(&t.handler)

	var err error
	t.endpoint, err = url.Parse(t.server.URL)
	AssertEq(nil, err)
}

func (t *HttpConnTest) TearDown() {
	t.server.Close()
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *HttpConnTest) InvalidScheme() {
	// Connection
	_, err := conn.NewHttpConn(&url.URL{Scheme: "taco", Host: "localhost"})

	ExpectThat(err, Error(HasSubstr("scheme")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *HttpConnTest) UnknownHost() {
	// Connection
	c, err := conn.NewHttpConn(&url.URL{Scheme: "http", Host: "foo.sidofhdksjhf"})
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	_, err = c.SendRequest(req)

	ExpectThat(err, Error(HasSubstr("foo.sidofhdksjhf")))
	ExpectThat(err, Error(HasSubstr("no such host")))
}

func (t *HttpConnTest) BasicHttpInfo() {
	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	_, err = c.SendRequest(req)
	AssertEq(nil, err)

	AssertNe(nil, t.handler.req)
	sysReq := t.handler.req

	ExpectEq("POST", sysReq.Method)
	ExpectEq("/", sysReq.URL.Path)
	ExpectEq(t.endpoint.Host, sysReq.Host)

	ExpectThat(
		sysReq.Header["Content-Type"],
		ElementsAre("application/x-www-form-urlencoded; charset=utf-8"))
}

func (t *HttpConnTest) RequestContainsNoParameters() {
	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	_, err = c.SendRequest(req)
	AssertEq(nil, err)

	AssertNe(nil, t.handler.reqBody)
	ExpectEq(0, len(t.handler.reqBody), "Body: %s", t.handler.reqBody)
}

func (t *HttpConnTest) RequestContainsOneParameter() {
	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{
		"taco": "burrito",
	}

	// Call
	_, err = c.SendRequest(req)
	AssertEq(nil, err)

	AssertNe(nil, t.handler.reqBody)
	ExpectEq(
		"taco=burrito",
		string(t.handler.reqBody))
}

func (t *HttpConnTest) RequestContainsMultipleParameters() {
	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{
		"taco":      "burrito",
		"enchilada": "queso",
		"nachos":    "carnitas",
	}

	// Call
	_, err = c.SendRequest(req)
	AssertEq(nil, err)
	AssertNe(nil, t.handler.reqBody)

	components := strings.Split(string(t.handler.reqBody), "&")
	AssertThat(components, ElementsAre(Any(), Any(), Any()))
	ExpectThat(components, Contains("taco=burrito"))
	ExpectThat(components, Contains("enchilada=queso"))
	ExpectThat(components, Contains("nachos=carnitas"))
}

func (t *HttpConnTest) ParametersNeedEscaping() {
	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{
		"타코":    "burrito",
		"b&az=": "qu?%x",
		"a+b":   "c d",
	}

	// Call
	_, err = c.SendRequest(req)
	AssertEq(nil, err)
	AssertNe(nil, t.handler.reqBody)

	components := strings.Split(string(t.handler.reqBody), "&")
	AssertThat(components, ElementsAre(Any(), Any(), Any()))
	ExpectThat(components, Contains("%ED%83%80%EC%BD%94=burrito"))
	ExpectThat(components, Contains("b%26az%3D=qu%3F%25x"))
	ExpectThat(components, Contains("a%2Bb=c%20d"))
}

func (t *HttpConnTest) ReturnsStatusCode() {
	// Handler
	t.handler.statusCode = 123

	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	resp, err := c.SendRequest(req)
	AssertEq(nil, err)

	ExpectEq(123, resp.StatusCode)
}

func (t *HttpConnTest) ReturnsBody() {
	// Handler
	t.handler.body = []byte{0xde, 0xad, 0x00, 0xbe, 0xef}

	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	resp, err := c.SendRequest(req)
	AssertEq(nil, err)

	ExpectThat(resp.Body, DeepEquals(t.handler.body))
}

func (t *HttpConnTest) ServerReturnsEmptyBody() {
	// Handler
	t.handler.body = []byte{}

	// Connection
	c, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)

	// Request
	req := conn.Request{}

	// Call
	resp, err := c.SendRequest(req)
	AssertEq(nil, err)

	ExpectThat(resp.Body, ElementsAre())
}

func (t *HttpConnTest) HttpsAllowed() {
	t.endpoint.Scheme = "https"

	// Connection
	_, err := conn.NewHttpConn(t.endpoint)
	AssertEq(nil, err)
}
