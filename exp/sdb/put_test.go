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

func (t *PutTest) EmptyItemName() {
	ExpectEq("TODO", "")
}

func (t *PutTest) InvalidItemName() {
	ExpectEq("TODO", "")
}

func (t *PutTest) ZeroUpdates() {
	ExpectEq("TODO", "")
}

func (t *PutTest) TooManyUpdates() {
	ExpectEq("TODO", "")
}

func (t *PutTest) OneAttributeNameEmpty() {
	ExpectEq("TODO", "")
}

func (t *PutTest) OneAttributeNameInvalid() {
	ExpectEq("TODO", "")
}

func (t *PutTest) OneAttributeValueInvalid() {
	ExpectEq("TODO", "")
}

func (t *PutTest) OnePreconditionNameInvalid() {
	ExpectEq("TODO", "")
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
