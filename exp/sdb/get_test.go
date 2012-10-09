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
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestGet(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type GetTest struct {
	domainTest

	item            string
	constistentRead bool
	names           []string

	attributes []Attribute
	err        error
}

func (t *GetTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.item = "foo"
}

func (t *GetTest) callDomain() {
	t.attributes, t.err = t.domain.GetAttributes(
		ItemName(t.item),
		t.constistentRead,
		t.names)
}

func init() { RegisterTestSuite(&GetTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *GetTest) ItemNameEmpty() {
	t.item = ""

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("item")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("empty")))
}

func (t *GetTest) ItemNameInvalid() {
	t.item = "taco\x80\x81\x82"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("item")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *GetTest) OneAttributeNameEmpty() {
	t.names = []string{"taco", "", "burrito"}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("empty")))
}

func (t *GetTest) OneAttributeNameInvalid() {
	t.names = []string{"taco", "\x80\x81\x82", "burrito"}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *GetTest) InconsistentReadWithNoAttributeNames() {
	ExpectEq("TODO", "")
}

func (t *GetTest) ConsistentRead() {
	ExpectEq("TODO", "")
}

func (t *GetTest) SomeAttributeNames() {
	ExpectEq("TODO", "")
}

func (t *GetTest) ConnReturnsError() {
	ExpectEq("TODO", "")
}

func (t *GetTest) ConnReturnsJunk() {
	ExpectEq("TODO", "")
}

func (t *GetTest) NoAttributesInResponse() {
	ExpectEq("TODO", "")
}

func (t *GetTest) SomeAttributesInResponse() {
	ExpectEq("TODO", "")
}
