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
	"errors"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/exp/sdb/conn"
	"github.com/jacobsa/aws/exp/sdb/conn/mock"
	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestConn(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type ConnTest struct {
	key      aws.AccessKey
	httpConn mock_conn.MockHttpConn
	signer   mock_conn.MockSigner

	c conn.Conn
}

func init() { RegisterTestSuite(&ConnTest{}) }

func (t *ConnTest) SetUp(i *TestInfo) {
	var err error

	t.key = aws.AccessKey{Id: "some_id", Secret: "some_secret"}
	t.httpConn = mock_conn.NewMockHttpConn(i.MockController, "httpConn")
	t.signer = mock_conn.NewMockSigner(i.MockController, "signer")

	t.c, err = conn.NewConn(t.key, t.httpConn, t.signer)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *ConnTest) CallsSigner() {
	req := conn.Request{
		"foo": "bar",
	}

	// Signer
	var signArg conn.Request
	ExpectCall(t.signer, "SignRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r conn.Request) error {
		signArg = r
		return errors.New("")
	}))

	// Call
	t.c.SendRequest(req)

	AssertNe(nil, signArg)
	AssertNe(req, signArg)

	ExpectEq("bar", signArg["foo"])
	ExpectEq("TODO", signArg["Timestamp"])
	ExpectEq("2", signArg["SignatureVersion"])
	ExpectEq("HmacSHA1", signArg["SignatureMethod"])
	ExpectEq(t.key.Id, signArg["AWSAccessKeyId"])
}

func (t *ConnTest) SignerReturnsError() {
	req := conn.Request{}

	// Signer
	ExpectCall(t.signer, "SignRequest")(Any()).
		WillOnce(oglemock.Return(errors.New("taco")))

	// Call
	_, err := t.c.SendRequest(req)

	ExpectThat(err, Error(HasSubstr("SignRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ConnTest) CallsHttpConn() {
	req := conn.Request{
		"foo": "bar",
	}

	// Signer
	ExpectCall(t.signer, "SignRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r conn.Request) error {
		// Add a parameter.
		r["baz"] = "qux"

		return nil
	}))

	// HTTP conn
	var sendArg conn.Request
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r conn.Request) (*conn.HttpResponse, error) {
		sendArg = r
		return nil, errors.New("")
	}))

	// Call
	t.c.SendRequest(req)

	AssertNe(nil, sendArg)
	AssertNe(req, sendArg)

	ExpectEq("bar", sendArg["foo"])
	ExpectEq("qux", sendArg["baz"])
	ExpectEq("TODO", sendArg["Timestamp"])
	ExpectEq("2", sendArg["SignatureVersion"])
	ExpectEq("HmacSHA1", sendArg["SignatureMethod"])
	ExpectEq(t.key.Id, sendArg["AWSAccessKeyId"])
}

func (t *ConnTest) HttpConnReturnsError() {
	req := conn.Request{}

	// Signer
	ExpectCall(t.signer, "SignRequest")(Any()).
		WillOnce(oglemock.Return(nil))

	// HTTP conn
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	_, err := t.c.SendRequest(req)

	ExpectThat(err, Error(HasSubstr("SendRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ConnTest) ServerReturnsError() {
	req := conn.Request{}

	// Signer
	ExpectCall(t.signer, "SignRequest")(Any()).
		WillOnce(oglemock.Return(nil))

	// HTTP conn
	httpResp := &conn.HttpResponse{StatusCode: 500, Body: []byte("taco")}
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(httpResp, nil))

	// Call
	_, err := t.c.SendRequest(req)

	ExpectThat(err, Error(HasSubstr("server")))
	ExpectThat(err, Error(HasSubstr("500")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ConnTest) ServerSaysOkay() {
	ExpectEq("TODO", "")
}
