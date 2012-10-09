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
	"github.com/jacobsa/aws/exp/sdb/conn"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"sort"
	"testing"
)

func TestPut(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func getSortedKeys(r conn.Request) []string {
	result := sort.StringSlice{}
	for key, _ := range r {
		result = append(result, key)
	}

	sort.Sort(result)
	return result
}

////////////////////////////////////////////////////////////////////////
// PutAttributes
////////////////////////////////////////////////////////////////////////

type PutTest struct {
	domainTest

	item          ItemName
	updates       []PutUpdate
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
	ExpectThat(t.err, Error(HasSubstr("UTF-8")))
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
	ExpectThat(t.err, Error(HasSubstr("257")))
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

func (t *PutTest) OnePreconditionMissingOperand() {
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

func (t *PutTest) OnePreconditionHasTwoOperands() {
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

func (t *PutTest) BasicParameters() {
	t.item = "some_item"
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "bar", Value: "taco", Replace: true},
		PutUpdate{Name: "baz", Value: "burrito"},
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Attribute.1.Name",
			"Attribute.1.Value",
			"Attribute.2.Name",
			"Attribute.2.Replace",
			"Attribute.2.Value",
			"Attribute.3.Name",
			"Attribute.3.Value",
			"DomainName",
			"ItemName",
		),
	)

	ExpectEq("foo", t.c.req["Attribute.1.Name"])
	ExpectEq("bar", t.c.req["Attribute.2.Name"])
	ExpectEq("baz", t.c.req["Attribute.3.Name"])

	ExpectEq("", t.c.req["Attribute.1.Value"])
	ExpectEq("taco", t.c.req["Attribute.2.Value"])
	ExpectEq("burrito", t.c.req["Attribute.3.Value"])

	ExpectEq("true", t.c.req["Attribute.2.Replace"])

	ExpectEq("some_item", t.c.req["ItemName"])
	ExpectEq(t.name, t.c.req["DomainName"])
}

func (t *PutTest) NoPreconditions() {
	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	ExpectThat(getSortedKeys(t.c.req), Not(Contains(HasSubstr("Expected"))))
}

func (t *PutTest) SomePreconditions() {
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

func (t *PutTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *PutTest) ConnSaysOkay() {
	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDomain()

	ExpectEq(nil, t.err)
}

////////////////////////////////////////////////////////////////////////
// BatchPutAttributes
////////////////////////////////////////////////////////////////////////

type BatchPutTest struct {
	domainTest

	updates map[ItemName][]PutUpdate

	err error
}

func init() { RegisterTestSuite(&BatchPutTest{}) }

func (t *BatchPutTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.updates = map[ItemName][]PutUpdate{
		"some_item": []PutUpdate{
			PutUpdate{Name: "foo"},
		},
	}
}

func (t *BatchPutTest) callDomain() {
	t.err = t.domain.BatchPutAttributes(t.updates)
}

func (t *BatchPutTest) NoItems() {
	t.updates = map[ItemName][]PutUpdate{}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("items")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *BatchPutTest) TooManyItems() {
	t.updates = map[ItemName][]PutUpdate{}

	for i := 0; i < 26; i++ {
		t.updates[ItemName(fmt.Sprintf("%d", i))] = []PutUpdate{
			PutUpdate{Name: "foo"},
		}
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("items")))
	ExpectThat(t.err, Error(HasSubstr("26")))
}

func (t *BatchPutTest) OneItemNameEmpty() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
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

func (t *BatchPutTest) OneItemNameInvalid() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
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

func (t *BatchPutTest) ZeroUpdatesForOneItem() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
		"foo": legalUpdates,
		"bar": []PutUpdate{},
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *BatchPutTest) TooManyUpdatesForOneItem() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
		"foo": legalUpdates,
		"bar": make([]PutUpdate, 257),
		"baz": legalUpdates,
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
	ExpectThat(t.err, Error(HasSubstr("257")))
}

func (t *BatchPutTest) OneAttributeNameEmpty() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
		"foo": legalUpdates,
		"bar": []PutUpdate{
			PutUpdate{Name: "qux"},
			PutUpdate{Name: ""},
			PutUpdate{Name: "wot"},
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

func (t *BatchPutTest) OneAttributeNameInvalid() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
		"foo": legalUpdates,
		"bar": []PutUpdate{
			PutUpdate{Name: "qux"},
			PutUpdate{Name: "taco\x80\x81\x82"},
			PutUpdate{Name: "wot"},
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

func (t *BatchPutTest) OneAttributeValueInvalid() {
	legalUpdates := []PutUpdate{PutUpdate{Name: "foo"}}
	t.updates = map[ItemName][]PutUpdate{
		"foo": legalUpdates,
		"bar": []PutUpdate{
			PutUpdate{Name: "a", Value: "qux"},
			PutUpdate{Name: "b", Value: "taco\x80\x81\x82"},
			PutUpdate{Name: "c", Value: "qux"},
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

func (t *BatchPutTest) CallsConn() {
	t.updates = map[ItemName][]PutUpdate{
		"bar": []PutUpdate{
			PutUpdate{Name: "a", Value: ""},
			PutUpdate{Name: "b", Value: "qux", Replace: true},
		},
		"foo": []PutUpdate{
			PutUpdate{Name: "c", Value: "wot"},
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
			"Item.1.Attribute.1.Value",
			"Item.1.Attribute.2.Name",
			"Item.1.Attribute.2.Replace",
			"Item.1.Attribute.2.Value",
			"Item.1.ItemName",
			"Item.2.Attribute.1.Name",
			"Item.2.Attribute.1.Value",
			"Item.2.ItemName",
		),
	)

	ExpectEq(t.name, t.c.req["DomainName"])

	ExpectEq("bar", t.c.req["Item.1.ItemName"])
	ExpectEq("foo", t.c.req["Item.2.ItemName"])

	ExpectEq("a", t.c.req["Item.1.Attribute.1.Name"])
	ExpectEq("b", t.c.req["Item.1.Attribute.2.Name"])
	ExpectEq("c", t.c.req["Item.2.Attribute.1.Name"])

	ExpectEq("", t.c.req["Item.1.Attribute.1.Value"])
	ExpectEq("qux", t.c.req["Item.1.Attribute.2.Value"])
	ExpectEq("wot", t.c.req["Item.2.Attribute.1.Value"])

	ExpectEq("true", t.c.req["Item.1.Attribute.2.Replace"])
}

func (t *BatchPutTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *BatchPutTest) ConnSaysOkay() {
	// Conn
	t.c.resp = []byte{}

	// Call
	t.callDomain()

	ExpectEq(nil, t.err)
}
