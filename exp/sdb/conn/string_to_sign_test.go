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
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestStringToSign(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type StringToSignTest struct {
}

func init() { RegisterTestSuite(&StringToSignTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *StringToSignTest) NoParameters() {
	req := Request{
	}

	str, err := computeStringToSign(req, "some_host.com")
	AssertEq(nil, err)

	ExpectEq(
		"POST\n" +
		"some_host.com\n" +
		"/\n" +
		"",
		str)
}

func (t *StringToSignTest) OneParameter() {
	ExpectEq("TODO", "")
}

func (t *StringToSignTest) MultipleParameters() {
	ExpectEq("TODO", "")
}

func (t *StringToSignTest) MixedCaseHost() {
	ExpectEq("TODO", "")
}

func (t *StringToSignTest) GoldenTest() {
	ExpectEq("TODO", "")
}
