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
	"errors"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"strings"
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

func (t *OpenDomainTest) NameIsTooShort() {
	t.name = "aa"

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("domain")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *OpenDomainTest) NameIsTooLong() {
	t.name = strings.Repeat("x", 256)

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("domain")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *OpenDomainTest) NameContainsUnusableCharacter() {
	t.name = "foo%bar"

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("domain")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("foo%bar")))
}

func (t *OpenDomainTest) CallsConn() {
	t.name = "f00_bar.baz-qux"

	// Call
	t.callDB()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
			"DomainName",
			"Version",
		),
	)

	ExpectEq("CreateDomain", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])

	ExpectEq("f00_bar.baz-qux", t.c.req["DomainName"])
}

func (t *OpenDomainTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *OpenDomainTest) ConnSaysOkay() {
	t.name = "taco"

	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDB()
	AssertEq(nil, t.err)

	castedDomain, ok := t.domain.(*domain)
	AssertTrue(ok)
	ExpectEq("taco", castedDomain.name)
	ExpectEq(t.c, castedDomain.c)
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

func (t *DeleteDomainTest) CallsConn() {
	// Call
	t.callDB()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
			"DomainName",
			"Version",
		),
	)

	ExpectEq("DeleteDomain", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])

	ExpectEq(t.domain.Name(), t.c.req["DomainName"])
}

func (t *DeleteDomainTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *DeleteDomainTest) ConnSaysOkay() {
	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDB()

	ExpectEq(nil, t.err)
}
