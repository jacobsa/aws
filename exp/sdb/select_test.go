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
	domainTest

	query string
	constistentRead bool
	nextToken []byte

	attrMap map[ItemName][]Attribute
	tok []byte
	err error
}

func (t *SelectTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)
}

func (t *SelectTest) callDomain() {
	t.attrMap, t.tok, t.err = t.domain.Select(t.query, t.constistentRead, t.nextToken)
}

func init() { RegisterTestSuite(&SelectTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *SelectTest) NoExtraOptions() {
	t.query = "taco"

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"SelectExpression",
		),
	)

	ExpectEq("taco", t.c.req["SelectExpression"])
}

func (t *SelectTest) ConistentRead() {
	t.constistentRead = true

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	ExpectEq("true", t.c.req["ConsistentRead"])
}

func (t *SelectTest) TokenPresent() {
	t.nextToken = []byte("taco")

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	ExpectEq("taco", t.c.req["NextToken"])
}

func (t *SelectTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *SelectTest) ConnReturnsJunk() {
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
	t.callDomain()

	AssertEq(nil, t.err)
	ExpectEq(0, len(t.attrMap))
	ExpectEq(nil, t.tok)
}

func (t *SelectTest) NoItemsInResponse() {
	ExpectEq("TODO", "")
}

func (t *SelectTest) SomeItemsInResponse() {
	ExpectEq("TODO", "")
}

func (t *SelectTest) NoNextTokenInResponse() {
	ExpectEq("TODO", "")
}

func (t *SelectTest) NextTokenInResponse() {
	ExpectEq("TODO", "")
}
