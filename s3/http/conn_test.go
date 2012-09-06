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

package http_test

import (
	"github.com/jacobsa/aws/s3/http"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestConn(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type ConnTest struct {
}

func init() { RegisterTestSuite(&ConnTest{}) }

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *ConnTest) InvalidScheme() {
	_, err := http.NewConn("localhost", "taco")

	ExpectThat(err, Error(HasSubstr("scheme")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ConnTest) UnknownHost() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerReturns200() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerReturns404() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerReturns500() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) ServerReturnsEmptyBody() {
	ExpectEq("TODO", "")
}

func (t *ConnTest) HttpsWorksProperly() {
	ExpectEq("TODO", "")
}
