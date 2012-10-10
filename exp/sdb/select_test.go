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
	"testing"
)

func TestSelect(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type SelectTest struct {
	simpleDBTest

	query           string
	constistentRead bool
	nextToken       []byte

	results []SelectedItem
	tok     []byte
	err     error
}

func (t *SelectTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.simpleDBTest.SetUp(i)
}

func (t *SelectTest) callDB() {
	t.results, t.tok, t.err = t.db.Select(t.query, t.constistentRead, t.nextToken)
}

func init() { RegisterTestSuite(&SelectTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *SelectTest) NoExtraOptions() {
	t.query = "taco"

	// Call
	t.callDB()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Action",
			"SelectExpression",
			"Version",
		),
	)

	ExpectEq("Select", t.c.req["Action"])
	ExpectEq("2009-04-15", t.c.req["Version"])

	ExpectEq("taco", t.c.req["SelectExpression"])
}

func (t *SelectTest) ConistentRead() {
	t.constistentRead = true

	// Call
	t.callDB()
	AssertNe(nil, t.c.req)

	ExpectEq("true", t.c.req["ConsistentRead"])
}

func (t *SelectTest) TokenPresent() {
	t.nextToken = []byte("taco")

	// Call
	t.callDB()
	AssertNe(nil, t.c.req)

	ExpectEq("taco", t.c.req["NextToken"])
}

func (t *SelectTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *SelectTest) ConnReturnsJunk() {
	// Conn
	t.c.resp = []byte("asdf")

	// Call
	t.callDB()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("asdf")))
}

func (t *SelectTest) NoItemsInResponse() {
	// Conn
	t.c.resp = []byte(`
		<SelectResponse>
		  <SelectResult>
		  </SelectResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</SelectResponse>`)

	// Call
	t.callDB()

	AssertEq(nil, t.err)
	ExpectThat(t.results, ElementsAre())
}

func (t *SelectTest) SomeItemsInResponse() {
	// Conn
	t.c.resp = []byte(`
		<SelectResponse>
		  <SelectResult>
				<Item>
					<Name>item_0</Name>
					<Attribute><Name>taco</Name><Value>burrito</Value></Attribute>
					<Attribute><Name>enchilada</Name><Value>queso</Value></Attribute>
				</Item>
				<Item>
					<Name>item_1</Name>
					<Attribute><Name>nachos</Name><Value>carnitas</Value></Attribute>
				</Item>
		  </SelectResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</SelectResponse>`)

	// Call
	t.callDB()

	AssertEq(nil, t.err)

	ExpectThat(
		t.results,
		ElementsAre(
			DeepEquals(
				SelectedItem{
					Name: "item_0",
					Attributes: []Attribute{
						Attribute{Name: "taco", Value: "burrito"},
						Attribute{Name: "enchilada", Value: "queso"},
					},
				},
			),
			DeepEquals(
				SelectedItem{
					Name: "item_1",
					Attributes: []Attribute{
						Attribute{Name: "nachos", Value: "carnitas"},
					},
				},
			),
		),
	)
}

func (t *SelectTest) NoNextTokenInResponse() {
	// Conn
	t.c.resp = []byte(`
		<SelectResponse>
		  <SelectResult>
				<Item><Name>item_0</Name></Item>
		  </SelectResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</SelectResponse>`)

	// Call
	t.callDB()

	AssertEq(nil, t.err)
	ExpectEq(nil, t.tok)
}

func (t *SelectTest) NextTokenInResponse() {
	// Conn
	t.c.resp = []byte(`
		<SelectResponse>
		  <SelectResult>
				<Item><Name>item_0</Name></Item>
				<NextToken>taco</NextToken>
		  </SelectResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</SelectResponse>`)

	// Call
	t.callDB()

	AssertEq(nil, t.err)
	ExpectThat(t.tok, DeepEquals([]byte("taco")))
}
