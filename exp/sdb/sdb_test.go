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
	. "github.com/jacobsa/ogletest"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

// A common helper class.
type simpleDBTest struct {
	c *fakeConn
	db SimpleDB
}

func (t *simpleDBTest) SetUp(i *TestInfo) {
	var err error

	t.c = &fakeConn{}

	t.db, err = newSimpleDB(t.c)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// CreateDomain
////////////////////////////////////////////////////////////////////////

type CreateDomainTest struct {
	simpleDBTest

	name string
	err error
}

func init() { RegisterTestSuite(&CreateDomainTest{}) }

func (t *CreateDomainTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.simpleDBTest.SetUp(i)

	// Make the request legal by default.
	t.name = "foo"
}

func (t *CreateDomainTest) callDB() {
	t.err = t.db.CreateDomain(t.name)
}

func (t *CreateDomainTest) DoesFoo() {
	ExpectFalse(true, "TODO")
}
