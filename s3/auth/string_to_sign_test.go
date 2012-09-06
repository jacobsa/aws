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
	"github.com/jacobsa/aws/s3/http"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestStringToSign(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type StringToSignTest struct {
}

func init() { RegisterTestSuite(&StringToSignTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *StringToSignTest) MissingDateHeader() {
	// Request
	req := &http.Request{
		Verb:    "PUT",
		Path:    "/foo/bar/baz",
		Headers: map[string]string{},
	}

	// Call
	_, err := stringToSign(req)

	ExpectThat(err, Error(HasSubstr("Date")))
}

func (t *StringToSignTest) MinimalRequest() {
	// Request
	req := &http.Request{
		Verb: "PUT",
		Path: "/foo/bar/baz",
		Headers: map[string]string{
			"Date": "Tue, 27 Mar 2007 19:36:42 +0000",
		},
	}

	// Call
	s, err := stringToSign(req)
	AssertEq(nil, err)

	ExpectThat(
		s,
		Equals(
			"PUT\n"+
				"\n"+ // Content-MD5
				"\n"+ // Content-Type
				"Tue, 27 Mar 2007 19:36:42 +0000\n"+
				"/foo/bar/baz"))
}

func (t *StringToSignTest) IncludesContentMd5() {
	// Request
	req := &http.Request{
		Verb: "PUT",
		Path: "/foo/bar/baz",
		Headers: map[string]string{
			"Date":        "some_date",
			"Content-MD5": "deadbeeffeedface",
		},
	}

	// Call
	s, err := stringToSign(req)
	AssertEq(nil, err)

	ExpectThat(
		s,
		Equals(
			"PUT\n"+
				"deadbeeffeedface\n"+
				"\n"+ // Content-Type
				"some_date\n"+
				"/foo/bar/baz"))
}

func (t *StringToSignTest) IncludesContentType() {
	// Request
	req := &http.Request{
		Verb: "PUT",
		Path: "/foo/bar/baz",
		Headers: map[string]string{
			"Date":         "some_date",
			"Content-Type": "blah/foo",
		},
	}

	// Call
	s, err := stringToSign(req)
	AssertEq(nil, err)

	ExpectThat(
		s,
		Equals(
			"PUT\n"+
				"\n"+ // Content-MD5
				"blah/foo\n"+
				"some_date\n"+
				"/foo/bar/baz"))
}

func (t *StringToSignTest) ComplicatedRequest() {
	// Request
	req := &http.Request{
		Verb: "PUT",
		Path: "/foo/bar/baz",
		Headers: map[string]string{
			"Date":         "some_date",
			"Content-MD5":  "deadbeeffeedface",
			"Content-Type": "blah/foo",
		},
	}

	// Call
	s, err := stringToSign(req)
	AssertEq(nil, err)

	ExpectThat(
		s,
		Equals(
			"PUT\n"+
				"deadbeeffeedface\n"+
				"blah/foo\n"+
				"some_date\n"+
				"/foo/bar/baz"))
}
