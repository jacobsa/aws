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

	item         ItemName
	updates      []DeleteUpdate
	precondition *Precondition

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
	t.err = t.domain.DeleteAttributes(t.item, t.updates, t.precondition)
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

func (t *DeleteTest) PreconditionNameEmpty() {
	t.precondition = &Precondition{
		Name: "",
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *DeleteTest) PreconditionNameInvalid() {
	t.precondition = &Precondition{
		Name: "taco\x80\x81\x82",
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.precondition.Name)))
}

func (t *DeleteTest) PreconditionValueInvalid() {
	t.precondition = &Precondition{
		Name:  "bar",
		Value: newString("taco\x80\x81\x82"),
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
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
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
			"Attribute.1.Name",
			"Attribute.2.Name",
			"Attribute.2.Value",
			"Attribute.3.Name",
			"Attribute.3.Value",
			"DomainName",
			"ItemName",
			"Version",
		),
	)

	ExpectEq("DeleteAttributes", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])
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
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
			"DomainName",
			"ItemName",
			"Version",
		),
	)

	ExpectEq("DeleteAttributes", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])
	ExpectEq(t.name, t.c.req["DomainName"])
	ExpectEq("some_item", t.c.req["ItemName"])
}

func (t *DeleteTest) NoPrecondition() {
	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	ExpectThat(getSortedKeys(t.c.req), Not(Contains(HasSubstr("Expected"))))
}

func (t *DeleteTest) NonExistencePrecondition() {
	t.precondition = &Precondition{
		Name: "foo",
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		AllOf(
			Contains("Expected.1.Name"),
			Contains("Expected.1.Exists"),
			Not(Contains("Expected.1.Value")),
		),
	)

	ExpectEq("foo", t.c.req["Expected.1.Name"])
	ExpectEq("false", t.c.req["Expected.1.Exists"])
}

func (t *DeleteTest) ValuePrecondition() {
	t.precondition = &Precondition{
		Name:  "foo",
		Value: newString("taco"),
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		AllOf(
			Contains("Expected.1.Name"),
			Contains("Expected.1.Value"),
			Not(Contains("Expected.1.Exists")),
		),
	)

	ExpectEq("foo", t.c.req["Expected.1.Name"])
	ExpectEq("taco", t.c.req["Expected.1.Value"])
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

	updates BatchDeleteMap

	err error
}

func init() { RegisterTestSuite(&BatchDeleteTest{}) }

func (t *BatchDeleteTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.updates = BatchDeleteMap{
		"some_item": []DeleteUpdate{
			DeleteUpdate{Name: "foo"},
		},
	}
}

func (t *BatchDeleteTest) callDomain() {
	t.err = t.domain.BatchDeleteAttributes(t.updates)
}

func (t *BatchDeleteTest) NoItems() {
	t.updates = BatchDeleteMap{}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("items")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *BatchDeleteTest) TooManyItems() {
	t.updates = BatchDeleteMap{}

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
	t.updates = BatchDeleteMap{
		"foo": legalUpdates,
		"":    legalUpdates,
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
	t.updates = BatchDeleteMap{
		"foo":             legalUpdates,
		"bar\x80\x81\x82": legalUpdates,
		"baz":             legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("item")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
}

func (t *BatchDeleteTest) TooManyUpdatesForOneItem() {
	legalUpdates := []DeleteUpdate{DeleteUpdate{Name: "foo"}}
	t.updates = BatchDeleteMap{
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
	t.updates = BatchDeleteMap{
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
	t.updates = BatchDeleteMap{
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
	t.updates = BatchDeleteMap{
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
	t.updates = BatchDeleteMap{
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
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
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
			"Version",
		),
	)

	ExpectEq("BatchDeleteAttributes", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])
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
