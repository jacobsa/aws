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
	t.item = "taco"
	t.constistentRead = false
	t.names = []string{}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"DomainName",
			"ItemName",
		),
	)

	ExpectEq(t.name, t.c.req["DomainName"])
	ExpectEq("taco", t.c.req["ItemName"])
}

func (t *GetTest) ConsistentRead() {
	t.constistentRead = true

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	ExpectEq("true", t.c.req["ConsistentRead"])
}

func (t *GetTest) SomeAttributeNames() {
	t.names = []string{"taco", "burrito"}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req, "Error: %v", t.err)

	AssertThat(
		getSortedKeys(t.c.req),
		AllOf(
			Contains("AttributeName.0"),
			Contains("AttributeName.1"),
		),
	)

	ExpectEq("taco", t.c.req["AttributeName.0"])
	ExpectEq("burrito", t.c.req["AttributeName.1"])
}

func (t *GetTest) ConnReturnsError() {
	// Conn
	t.c.err = errors.New("taco")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("SendRequest")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *GetTest) ConnReturnsJunk() {
	// Conn
	t.c.resp = []byte("asdf")

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("server")))
	ExpectThat(t.err, Error(HasSubstr("asdf")))
}

func (t *GetTest) NoAttributesInResponse() {
	// Conn
	t.c.resp = []byte(`
		<GetAttributesResponse>
		  <GetAttributesResult>
		  </GetAttributesResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</GetAttributesResponse>`)

	// Call
	t.callDomain()

	AssertEq(nil, t.err)
	ExpectThat(t.attributes, ElementsAre())
}

func (t *GetTest) SomeAttributesInResponse() {
	// Conn
	t.c.resp = []byte(`
		<GetAttributesResponse>
		  <GetAttributesResult>
		    <Attribute><Name>taco</Name><Value>burrito</Value></Attribute>
		    <Attribute><Name>enchilada</Name><Value>queso</Value></Attribute>
		  </GetAttributesResult>
		  <ResponseMetadata>
		    <RequestId>b1e8f1f7-42e9-494c-ad09-2674e557526d</RequestId>
		    <BoxUsage>0.0000219907</BoxUsage>
		  </ResponseMetadata>
		</GetAttributesResponse>`)

	// Call
	t.callDomain()

	AssertEq(nil, t.err)
	ExpectThat(
		t.attributes,
		ElementsAre(
			DeepEquals(Attribute{Name: "taco", Value: "burrito"}),
			DeepEquals(Attribute{Name: "enchilada", Value: "queso"}),
		),
	)
}
