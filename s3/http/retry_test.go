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
	"errors"
	"github.com/jacobsa/aws/s3/http"
	"github.com/jacobsa/aws/s3/http/mock"
	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
	"net"
	"syscall"
	"testing"
)

func TestRetry(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type RetryingConnTest struct {
	wrapped mock_http.MockConn
	conn http.Conn

	req *http.Request
	resp *http.Response
	err error
}

func init() { RegisterTestSuite(&RetryingConnTest{}) }

func (t *RetryingConnTest) SetUp(i *TestInfo) {
	var err error

	t.wrapped = mock_http.NewMockConn(i.MockController, "wrapped")
	t.conn, err = http.NewRetryingConn(t.wrapped)
	AssertEq(nil, err)
}

func (t *RetryingConnTest) call() {
	t.resp, t.err = t.conn.SendRequest(t.req)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *RetryingConnTest) CallsWrapped() {
	t.req = &http.Request{}

	// Wrapped
	ExpectCall(t.wrapped, "SendRequest")(t.req).
		WillOnce(oglemock.Return(nil, nil))

	// Call
	t.call()
}

func (t *RetryingConnTest) WrappedReturnsWrongErrorType() {
	// Wrapped
	ExpectCall(t.wrapped, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	t.call()

	ExpectThat(t.err, Error(Equals("taco")))
}

func (t *RetryingConnTest) WrappedReturnsWrongOpErrorType() {
	// Wrapped
	wrappedErr := &net.OpError{
		Op: "taco",
		Err: errors.New("burrito"),
	}

	ExpectCall(t.wrapped, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, wrappedErr))

	// Call
	t.call()

	ExpectThat(t.err, Error(HasSubstr("taco")))
	ExpectThat(t.err, Error(HasSubstr("burrito")))
}

func (t *RetryingConnTest) WrappedReturnsUninterestingErrno() {
	// Wrapped
	wrappedErr := &net.OpError{
		Op: "taco",
		Err: syscall.EMLINK,
	}

	ExpectCall(t.wrapped, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, wrappedErr))

	// Call
	t.call()

	ExpectThat(t.err, Error(HasSubstr("taco")))
	ExpectThat(t.err, Error(HasSubstr("too many links")))
}

func (t *RetryingConnTest) RetriesForBrokenPipe() {
	t.req = &http.Request{}

	// Wrapped
	wrappedErr := &net.OpError{
		Err: syscall.EPIPE,
	}

	ExpectCall(t.wrapped, "SendRequest")(t.req).
		WillOnce(oglemock.Return(nil, wrappedErr)).
		WillOnce(oglemock.Return(nil, wrappedErr)).
		WillOnce(oglemock.Return(nil, nil))

	// Call
	t.call()
}

func (t *RetryingConnTest) WrappedFailsOnThirdCall() {
	// Wrapped
	wrappedErr0 := &net.OpError{
		Op: "taco",
		Err: syscall.EPIPE,
	}

	wrappedErr1 := wrappedErr0

	wrappedErr2 := &net.OpError{
		Op: "burrito",
		Err: syscall.EPIPE,
	}

	ExpectCall(t.wrapped, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, wrappedErr0)).
		WillOnce(oglemock.Return(nil, wrappedErr1)).
		WillOnce(oglemock.Return(nil, wrappedErr2))

	// Call
	t.call()

	ExpectThat(t.err, Error(HasSubstr("burrito")))
	ExpectThat(t.err, Error(HasSubstr("broken pipe")))
}

func (t *RetryingConnTest) WrappedSucceedsOnThirdCall() {
	// Wrapped
	wrappedErr := &net.OpError{
		Err: syscall.EPIPE,
	}

	expected := &http.Response{}

	ExpectCall(t.wrapped, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, wrappedErr)).
		WillOnce(oglemock.Return(nil, wrappedErr)).
		WillOnce(oglemock.Return(expected, nil))

	// Call
	t.call()

	AssertEq(nil, t.err)
	ExpectEq(expected, t.resp)
}
