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

package sdb

import (
	"github.com/jacobsa/aws/exp/sdb/conn/mock"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestDomain(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type DomainTest struct {
	name string
	c mock_conn.MockConn
	domain Domain
}

func init() { RegisterTestSuite(&DomainTest{}) }

func (t *DomainTest) SetUp(i *TestInfo) {
	var err error

	t.name = "some_domain"
	t.c = mock_conn.NewMockConn(i.MockController, "conn")

	t.domain, err = newDomain(t.name, t.c)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *DomainTest) DoesFoo() {
	ExpectEq("TODO", "")
}
