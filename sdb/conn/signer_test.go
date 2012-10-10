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
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/jacobsa/aws"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestSigner(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type SignerTest struct {
}

func init() { RegisterTestSuite(&SignerTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *SignerTest) CallsFunction() {
	// Function
	var reqArg Request
	var hostArg string

	sts := func(r Request, host string) (string, error) {
		reqArg = r
		hostArg = host
		return "", nil
	}

	// Signer
	signer := newSigner(aws.AccessKey{}, "some_host", sts)

	// Call
	req := Request{"foo": "bar"}
	signer.SignRequest(req)

	ExpectEq(req, reqArg)
	ExpectEq("some_host", hostArg)
}

func (t *SignerTest) FunctionReturnsError() {
	// Function
	sts := func(r Request, h string) (string, error) {
		return "", errors.New("taco")
	}

	// Signer
	signer := newSigner(aws.AccessKey{}, "", sts)

	// Call
	req := Request{}
	err := signer.SignRequest(req)

	ExpectThat(err, Error(HasSubstr("computeStringToSign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *SignerTest) FunctionReturnsString() {
	// Function
	sts := func(r Request, h string) (string, error) {
		return "taco", nil
	}

	// Signer
	key := aws.AccessKey{Id: "queso", Secret: "burrito"}
	signer := newSigner(key, "", sts)

	// Expected output
	h := hmac.New(sha1.New, []byte("burrito"))
	_, err := h.Write([]byte("taco"))
	AssertEq(nil, err)

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	_, err = encoder.Write(h.Sum(nil))
	AssertEq(nil, err)
	AssertEq(nil, encoder.Close())

	expected := buf.String()

	// Call
	req := Request{
		"foo": "bar",
	}

	err = signer.SignRequest(req)
	AssertEq(nil, err)

	ExpectEq("bar", req["foo"])
	ExpectEq(expected, req["Signature"])
}
