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
	"fmt"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestDelete(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func newString(s string) *string {
	return &s
}

////////////////////////////////////////////////////////////////////////
// DeleteAttributes
////////////////////////////////////////////////////////////////////////

type DeleteTest struct {
	domainTest

	item          ItemName
	updates       []DeleteUpdate
	preconditions []Precondition

	err error
}

func init() { RegisterTestSuite(&DeleteTest{}) }

func (t *DeleteTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.item = "foo"
	t.updates = []DeleteUpdate{DeleteUpdate{Name: "bar"}}
}

func (t *DeleteTest) callDomain() {
	t.err = t.domain.DeleteAttributes(t.item, t.updates, t.preconditions)
}

func (t *DeleteTest) EmptyItemName() {
	t.item = ""

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
}

func (t *DeleteTest) InvalidItemName() {
	t.item = "taco\x80\x81\x82"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *DeleteTest) OneAttributeNameEmpty() {
	t.updates = []DeleteUpdate{
		DeleteUpdate{Name: "foo"},
		DeleteUpdate{Name: ""},
		DeleteUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *DeleteTest) OneAttributeNameInvalid() {
	t.updates = []DeleteUpdate{
		DeleteUpdate{Name: "foo"},
		DeleteUpdate{Name: "taco\x80\x81\x82"},
		DeleteUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.updates[1].Name)))
}

func (t *DeleteTest) OneAttributeValueInvalid() {
	t.updates = []DeleteUpdate{
		DeleteUpdate{Name: "foo"},
		DeleteUpdate{Name: "bar", Value: newString("taco\x80\x81\x82")},
		DeleteUpdate{Name: "baz"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
}

func (t *DeleteTest) OnePreconditionNameEmpty() {
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

func (t *DeleteTest) OnePreconditionNameInvalid() {
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

func (t *DeleteTest) OnePreconditionValueInvalid() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Value: new(string)},
		Precondition{Name: "bar", Value: new(string)},
		Precondition{Name: "baz", Value: new(string)},
	}

	*t.preconditions[0].Value = ""
	*t.preconditions[1].Value = "taco\x80\x81\x82"
	*t.preconditions[2].Value = "qux"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
}

func (t *DeleteTest) OnePreconditionMissingOperand() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "bar"},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("precondition")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
}

func (t *DeleteTest) OnePreconditionHasTwoOperands() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "bar", Exists: new(bool), Value: new(string)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("precondition")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
}

func (t *DeleteTest) BasicParameters() {
	t.item = "some_item"
	t.updates = []DeleteUpdate{
		DeleteUpdate{Name: "foo"},
		DeleteUpdate{Name: "bar", Value: newString("")},
		DeleteUpdate{Name: "baz", Value: newString("taco")},
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Attribute.1.Name",
			"Attribute.2.Name",
			"Attribute.2.Value",
			"Attribute.3.Name",
			"Attribute.3.Value",
			"DomainName",
			"ItemName",
		),
	)

	ExpectEq(t.name, t.c.req["DomainName"])
	ExpectEq("some_item", t.c.req["ItemName"])

	ExpectEq("foo", t.c.req["Attribute.1.Name"])
	ExpectEq("bar", t.c.req["Attribute.2.Name"])
	ExpectEq("baz", t.c.req["Attribute.3.Name"])

	ExpectEq("", t.c.req["Attribute.2.Value"])
	ExpectEq("taco", t.c.req["Attribute.3.Value"])
}

func (t *DeleteTest) NoUpdates() {
	t.item = "some_item"
	t.updates = []DeleteUpdate{}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"DomainName",
			"ItemName",
		),
	)

	ExpectEq(t.name, t.c.req["DomainName"])
	ExpectEq("some_item", t.c.req["ItemName"])
}

func (t *DeleteTest) NoPreconditions() {
	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	ExpectThat(getSortedKeys(t.c.req), Not(Contains(HasSubstr("Expected"))))
}

func (t *DeleteTest) SomePreconditions() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "bar", Value: new(string)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	*t.preconditions[0].Exists = false
	*t.preconditions[1].Value = "taco"
	*t.preconditions[2].Exists = true

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		AllOf(
			Contains("Expected.1.Name"),
			Contains("Expected.2.Name"),
			Contains("Expected.3.Name"),
			Contains("Expected.1.Exists"),
			Contains("Expected.2.Value"),
			Contains("Expected.3.Exists"),
		),
	)

	ExpectEq("foo", t.c.req["Expected.1.Name"])
	ExpectEq("bar", t.c.req["Expected.2.Name"])
	ExpectEq("baz", t.c.req["Expected.3.Name"])

	ExpectEq("false", t.c.req["Expected.1.Exists"])
	ExpectEq("taco", t.c.req["Expected.2.Value"])
	ExpectEq("true", t.c.req["Expected.3.Exists"])
}

func (t *DeleteTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *DeleteTest) ConnSaysOkay() {
	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDomain()

	ExpectEq(nil, t.err)
}

////////////////////////////////////////////////////////////////////////
// BatchDeleteAttributes
////////////////////////////////////////////////////////////////////////

type BatchDeleteTest struct {
	domainTest

	updates       map[ItemName][]DeleteUpdate

	err error
}

func init() { RegisterTestSuite(&BatchDeleteTest{}) }

func (t *BatchDeleteTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.updates = map[ItemName][]DeleteUpdate{
		"some_item": []DeleteUpdate{
			DeleteUpdate{Name: "foo"},
		},
	}
}

func (t *BatchDeleteTest) callDomain() {
	t.err = t.domain.BatchDeleteAttributes(t.updates)
}

func (t *BatchDeleteTest) NoItems() {
	t.updates = map[ItemName][]DeleteUpdate{
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("items")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *BatchDeleteTest) TooManyItems() {
	t.updates = map[ItemName][]DeleteUpdate{}

	for i := 0; i < 26; i++ {
		t.updates[ItemName(fmt.Sprintf("%d", i))] = []DeleteUpdate{
			DeleteUpdate{Name: "foo"},
		}
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("items")))
	ExpectThat(t.err, Error(HasSubstr("26")))
}

func (t *BatchDeleteTest) OneItemNameEmpty() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"": legalUpdates,
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("item")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("empty")))
}

func (t *BatchDeleteTest) OneItemNameInvalid() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar\x80\x81\x82": legalUpdates,
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("item")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *BatchDeleteTest) ZeroUpdatesForOneItem() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar": []DeleteUpdate{},
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *BatchDeleteTest) TooManyUpdatesForOneItem() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar": make([]DeleteUpdate, 257),
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("257")))
}

func (t *BatchDeleteTest) OneAttributeNameEmpty() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar": []DeleteUpdate{
			DeleteUpdate{Name: "qux"},
			DeleteUpdate{Name: ""},
			DeleteUpdate{Name: "wot"},
		},
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("empty")))
}

func (t *BatchDeleteTest) OneAttributeNameInvalid() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar": []DeleteUpdate{
			DeleteUpdate{Name: "qux"},
			DeleteUpdate{Name: "taco\x80\x81\x82"},
			DeleteUpdate{Name: "wot"},
		},
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *BatchDeleteTest) OneAttributeValueInvalid() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = map[ItemName][]DeleteUpdate{
		"foo": legalUpdates,
		"bar": []DeleteUpdate{
			DeleteUpdate{Name: "a"},
			DeleteUpdate{Name: "b", Value: newString("taco\x80\x81\x82")},
			DeleteUpdate{Name: "c", Value: newString("qux")},
		},
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *BatchDeleteTest) CallsConn() {
	t.updates = map[ItemName][]DeleteUpdate{
		"bar": []DeleteUpdate{
			DeleteUpdate{Name: "a"},
			DeleteUpdate{Name: "b", Value: newString("qux")},
			DeleteUpdate{Name: "c", Value: newString("")},
		},
		"baz": []DeleteUpdate{},
		"foo": []DeleteUpdate{
			DeleteUpdate{Name: "d", Value: newString("wot")},
		},
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"DomainName",
			"Item.1.Attribute.1.Name",
			"Item.1.Attribute.2.Name",
			"Item.1.Attribute.2.Value",
			"Item.1.Attribute.3.Name",
			"Item.1.Attribute.3.Value",
			"Item.1.ItemName",
			"Item.2.ItemName",
			"Item.3.Attribute.1.Name",
			"Item.3.Attribute.1.Value",
			"Item.3.ItemName",
		),
	)

	ExpectEq(t.name, t.c.req["DomainName"])

	ExpectEq("bar", t.c.req["Item.1.ItemName"])
	ExpectEq("baz", t.c.req["Item.2.ItemName"])
	ExpectEq("foo", t.c.req["Item.3.ItemName"])

	ExpectEq("a", t.c.req["Item.1.Attribute.1.Name"])
	ExpectEq("b", t.c.req["Item.1.Attribute.2.Name"])
	ExpectEq("c", t.c.req["Item.1.Attribute.3.Name"])
	ExpectEq("d", t.c.req["Item.3.Attribute.1.Name"])

	ExpectEq("qux", t.c.req["Item.1.Attribute.2.Value"])
	ExpectEq("", t.c.req["Item.1.Attribute.3.Value"])
	ExpectEq("wot", t.c.req["Item.3.Attribute.1.Value"])
}

func (t *BatchDeleteTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *BatchDeleteTest) ConnSaysOkay() {
	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDomain()

	ExpectEq(nil, t.err)
}
