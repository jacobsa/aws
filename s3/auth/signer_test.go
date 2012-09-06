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
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3/http"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestSigner(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type SignerTest struct{}

func init() { RegisterTestSuite(&SignerTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *SignerTest) CallsFunction() {
	// Function
	var stsArg *http.Request
	sts := func(r *http.Request) (string, error) { stsArg = r; return "", nil }

	// Signer
	signer, err := newSigner(sts, &aws.AccessKey{})
	AssertEq(nil, err)

	// Call
	req := &http.Request{Headers: make(map[string]string)}
	signer.Sign(req)

	ExpectEq(req, stsArg)
}

func (t *SignerTest) FunctionReturnsError() {
	// Function
	sts := func(r *http.Request) (string, error) { return "", errors.New("taco") }

	// Signer
	signer, err := newSigner(sts, &aws.AccessKey{})
	AssertEq(nil, err)

	// Call
	err = signer.Sign(&http.Request{Headers: make(map[string]string)})

	ExpectThat(err, Error(HasSubstr("Sign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *SignerTest) FunctionReturnsString() {
	// Function
	sts := func(r *http.Request) (string, error) { return "taco", nil }

	// Signer
	key := &aws.AccessKey{Id: "queso", Secret: "burrito"}
	signer, err := newSigner(sts, key)
	AssertEq(nil, err)

	// Expected output
	h := hmac.New(sha1.New, []byte("burrito"))
	_, err = h.Write([]byte("taco"))
	AssertEq(nil, err)

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	_, err = encoder.Write(h.Sum(nil))
	AssertEq(nil, err)
	AssertEq(nil, encoder.Close())

	expected := "AWS queso:" + buf.String()

	// Call
	req := &http.Request{
		Headers: map[string]string{
			"foo": "bar",
		},
	}

	err = signer.Sign(req)
	AssertEq(nil, err)

	ExpectEq("bar", req.Headers["foo"])
	ExpectEq(expected, req.Headers["Authorization"])
}

func (t *SignerTest) GoldenTests() {
	type testCase struct {
		stringToSign        string
		expectedHeaderValue string
	}

	// Golden tests taken from Amazon doc examples.
	key := &aws.AccessKey{
		Id:     "AKIAIOSFODNN7EXAMPLE",
		Secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	cases := []testCase{
		testCase{
			"GET\n\n\nTue, 27 Mar 2007 19:36:42 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:bWq2s1WEIj+Ydj0vQ697zp+IXMU=",
		},
		testCase{
			"PUT\n\nimage/jpeg\nTue, 27 Mar 2007 21:15:45 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:MyyxeRY7whkBe+bq8fHCL/2kKUg=",
		},
		testCase{
			"GET\n\n\nTue, 27 Mar 2007 19:42:41 +0000\n/johnsmith/",
			"AWS AKIAIOSFODNN7EXAMPLE:htDYFYduRNen8P9ZfE/s9SuKy0U=",
		},
		testCase{
			"GET\n\n\nTue, 27 Mar 2007 19:44:46 +0000\n/johnsmith/?acl",
			"AWS AKIAIOSFODNN7EXAMPLE:c2WLPFtWHVgbEmeEG93a4cG37dM=",
		},
		testCase{
			"DELETE\n\n\nx-amz-date:Tue, 27 Mar 2007 21:20:26 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:9b2sXq0KfxsxHtdZkzx/9Ngqyh8=",
		},
	}

	for i, c := range cases {
		// Function
		sts := func(r *http.Request) (string, error) { return c.stringToSign, nil }

		// Signer
		signer, err := newSigner(sts, key)
		AssertEq(nil, err)

		// Call
		req := &http.Request{Headers: make(map[string]string)}
		err = signer.Sign(req)
		AssertEq(nil, err)

		ExpectEq(c.expectedHeaderValue, req.Headers["Authorization"], "Case %d: %v", i, c)
	}
}
