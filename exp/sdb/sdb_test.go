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

type fakeDomain struct {
	name string
}

func (d *fakeDomain) Name() string {
	return d.name
}

func (d *fakeDomain) PutAttributes(
		item ItemName,
		updates []PutUpdate,
		preconditions []Precondition) error {
	panic("Unsupported")
}

func (d *fakeDomain) BatchPutAttributes(updateMap map[ItemName][]PutUpdate) error {
	panic("Unsupported")
}

func (d *fakeDomain) DeleteAttributes(
		item ItemName,
		deletes []DeleteUpdate,
		preconditions []Precondition) error {
	panic("Unsupported")
}

func (d *fakeDomain) BatchDeleteAttributes(deleteMap map[ItemName][]DeleteUpdate) error {
	panic("Unsupported")
}

func (d *fakeDomain) GetAttributes(
		item ItemName,
		constistentRead bool,
		attrNames []string) (attrs []Attribute, err error) {
	panic("Unsupported")
}

////////////////////////////////////////////////////////////////////////
// OpenDomain
////////////////////////////////////////////////////////////////////////

type OpenDomainTest struct {
	simpleDBTest

	name string

	domain Domain
	err error
}

func init() { RegisterTestSuite(&OpenDomainTest{}) }

func (t *OpenDomainTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.simpleDBTest.SetUp(i)

	// Make the request legal by default.
	t.name = "foo"
}

func (t *OpenDomainTest) callDB() {
	t.domain, t.err = t.db.OpenDomain(t.name)
}

func (t *OpenDomainTest) NameIsEmpty() {
	ExpectFalse(true, "TODO")
}

func (t *OpenDomainTest) NameIsInvalid() {
	ExpectFalse(true, "TODO")
}

func (t *OpenDomainTest) CallsConn() {
	ExpectFalse(true, "TODO")
}

func (t *OpenDomainTest) ConnReturnsError() {
	ExpectFalse(true, "TODO")
}

func (t *OpenDomainTest) CallsFactoryFuncAndReturnsResult() {
	ExpectFalse(true, "TODO")
}

////////////////////////////////////////////////////////////////////////
// DeleteDomain
////////////////////////////////////////////////////////////////////////

type DeleteDomainTest struct {
	simpleDBTest

	domain Domain
	err error
}

func init() { RegisterTestSuite(&DeleteDomainTest{}) }

func (t *DeleteDomainTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.simpleDBTest.SetUp(i)

	// Set up a fake named domain.
	t.domain = &fakeDomain{"some_domain"}
}

func (t *DeleteDomainTest) callDB() {
	t.err = t.db.DeleteDomain(t.domain)
}

func (t *DeleteDomainTest) DoesFoo() {
	ExpectFalse(true, "TODO")
}
