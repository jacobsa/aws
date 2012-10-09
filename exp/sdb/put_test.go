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

func TestPut(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// PutAttributes
////////////////////////////////////////////////////////////////////////

type PutTest struct {
	domainTest

	item ItemName
	updates []PutUpdate
	preconditions []Precondition

	err error
}

func init() { RegisterTestSuite(&PutTest{}) }

func (t *PutTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.item = "foo"
	t.updates = []PutUpdate{PutUpdate{"bar", "baz", false}}
}

func (t *PutTest) callDomain() {
	t.err = t.domain.PutAttributes(t.item, t.updates, t.preconditions)
}

func (t *PutTest) EmptyItemName() {
	t.item = ""

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
}

func (t *PutTest) InvalidItemName() {
	t.item = "taco\x80\x81\x82"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
	ExpectThat(t.err, Error(HasSubstr(string(t.item))))
}

func (t *PutTest) ZeroUpdates() {
	t.updates = []PutUpdate{}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *PutTest) TooManyUpdates() {
	t.updates = make([]PutUpdate, 257)

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("256")))
}

func (t *PutTest) OneAttributeNameEmpty() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "", Value: "taco"},
		PutUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *PutTest) OneAttributeNameInvalid() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "taco\x80\x81\x82"},
		PutUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.updates[1].Name)))
}

func (t *PutTest) OneAttributeValueInvalid() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "bar", Value: "taco\x80\x81\x82"},
		PutUpdate{Name: "baz"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
	ExpectThat(t.err, Error(HasSubstr(t.updates[1].Value)))
}

func (t *PutTest) OnePreconditionNameEmpty() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "", Exists: new(bool)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *PutTest) OnePreconditionNameInvalid() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "taco\x80\x81\x82", Exists: new(bool)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.preconditions[1].Name)))
}

func (t *PutTest) OnePreconditionValueInvalid() {
	ExpectEq("TODO", "")
}

func (t *PutTest) NoPreconditions() {
	ExpectEq("TODO", "")
}

func (t *PutTest) SomePreconditions() {
	ExpectEq("TODO", "")
}

func (t *PutTest) ConnReturnsError() {
	ExpectEq("TODO", "")
}

func (t *PutTest) ConnSaysOkay() {
	ExpectEq("TODO", "")
}

////////////////////////////////////////////////////////////////////////
// BatchPutAttributes
////////////////////////////////////////////////////////////////////////

type BatchPutTest struct {
	domainTest
}

func init() { RegisterTestSuite(&BatchPutTest{}) }

func (t *BatchPutTest) DoesFoo() {
	ExpectEq("TODO", "")
}
