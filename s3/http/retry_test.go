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
	"github.com/jacobsa/aws/s3/http/mock"
	"github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
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

	// Wrapped (first call)
	ExpectCall(t.wrapped, "SendRequest")(t.req).
		WillOnce(oglemock.Return(nil, nil))

	// Call
	t.call()
}

func (t *RetryingConnTest) WrappedReturnsWrongErrorType() {
	ExpectEq("TODO", "")
}

func (t *RetryingConnTest) WrappedReturnsWrongOpErrorType() {
	ExpectEq("TODO", "")
}

func (t *RetryingConnTest) WrappedReturnsUnknownErrno() {
	ExpectEq("TODO", "")
}

func (t *RetryingConnTest) RetriesForBrokenPipe() {
	ExpectEq("TODO", "")
}

func (t *RetryingConnTest) WrappedFailsOnThirdCall() {
	ExpectEq("TODO", "")
}

func (t *RetryingConnTest) WrappedSucceedsOnThirdCall() {
	ExpectEq("TODO", "")
}
