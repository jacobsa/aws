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
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/exp/sdb/conn"
	"github.com/jacobsa/aws/exp/sdb/conn/mock"
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

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *ConnTest) CallsSigner() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) SignerReturnsError() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) CallsHttpConn() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) HttpConnReturnsError() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerReturnsError() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerSaysOkay() {
	ExpectEq("TODO", "")
}
