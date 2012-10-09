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

package conn

import (
	"bytes"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"net/url"
	"strings"
	"testing"
)

func TestPostBody(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type PostBodyTest struct {
}

func init() { RegisterTestSuite(&PostBodyTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *PostBodyTest) NoParameters() {
	req := Request{}

	body := assemblePostBody(req)

	ExpectEq("", body)
}

func (t *PostBodyTest) OneParameter() {
	req := Request{
		"taco": "burrito",
	}

	body := assemblePostBody(req)

	ExpectEq("taco=burrito", body)
}

func (t *PostBodyTest) MultipleParameters() {
	req := Request{
		"taco":      "burrito",
		"enchilada": "queso",
		"nachos":    "carnitas",
	}

	body := assemblePostBody(req)
	components := strings.Split(string(body), "&")

	AssertEq(3, len(components), "Components: %v", components)
	ExpectThat(components, Contains("taco=burrito"))
	ExpectThat(components, Contains("enchilada=queso"))
	ExpectThat(components, Contains("nachos=carnitas"))
}

func (t *PostBodyTest) EmptyParameterName() {
	req := Request{
		"":          "burrito",
		"enchilada": "queso",
	}

	body := assemblePostBody(req)
	components := strings.Split(string(body), "&")

	AssertEq(2, len(components), "Components: %v", components)
	ExpectThat(components, Contains("=burrito"))
	ExpectThat(components, Contains("enchilada=queso"))
}

func (t *PostBodyTest) UnreservedCharacters() {
	req := Request{
		"abcdefghijklmnopqrstuvwxyz": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789":                 "-_.~",
	}

	body := assemblePostBody(req)

	ExpectThat(body, HasSubstr("abcdefghijklmnopqrstuvwxyz"))
	ExpectThat(body, HasSubstr("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	ExpectThat(body, HasSubstr("0123456789"))
	ExpectThat(body, HasSubstr("-_.~"))
}

func (t *PostBodyTest) StructuralCharacters() {
	req := Request{
		":/?#[]@": "!$&'()*+,;=",
	}

	body := assemblePostBody(req)

	ExpectThat(body, HasSubstr("%3A%2F%3F%23%5B%5D%40="))
	ExpectThat(body, HasSubstr("=%21%24%26%27%28%29%2A%2B%2C%3B%3D"))
}

func (t *PostBodyTest) PercentCharacter() {
	req := Request{
		"a%b": "c%d",
	}

	body := assemblePostBody(req)

	ExpectThat(body, HasSubstr("a%25b=c%25d"))
}

func (t *PostBodyTest) SpaceAndPlus() {
	req := Request{
		"b+a z": "q+u x",
	}

	body := assemblePostBody(req)
	ExpectEq("b%2Ba%20z=q%2Bu%20x", body)
}

func (t *PostBodyTest) KoreanCharacters() {
	req := Request{
		"음식": "타코",
	}

	body := assemblePostBody(req)
	ExpectEq("%EC%9D%8C%EC%8B%9D=%ED%83%80%EC%BD%94", body)
}

func (t *PostBodyTest) ParameterOrdering() {
	// Sanity check: ordering of Korean.
	AssertEq(-1, bytes.Compare([]byte("음"), []byte("타")))

	// Sanity check: ordering of unescaped strings reversed by escaping. (We
	// should order by parameter name, not be escaped parameter name.)
	AssertLt("f", "|")
	AssertGt(url.QueryEscape("f"), url.QueryEscape("|"))

	// Request
	req := Request{
		// Easy cases
		"bar": "asd",
		"qux": "asd",
		"aaa": "asd",

		// Korean ordering
		"타코": "bar",
		"음식": "foo",

		// Order before escaping
		"foo": "asd",
		"|":   "asd",
	}

	body := assemblePostBody(req)
	components := strings.Split(string(body), "&")

	ExpectThat(
		components,
		ElementsAre(
			HasSubstr("aaa="),
			HasSubstr("bar="),
			HasSubstr("foo="),
			HasSubstr("qux="),
			HasSubstr(url.QueryEscape("|")+"="),
			HasSubstr(url.QueryEscape("음식")+"="),
			HasSubstr(url.QueryEscape("타코")+"=")))
}
