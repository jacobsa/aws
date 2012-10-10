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

package main

import (
	"github.com/jacobsa/aws/exp/sdb"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"sync"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type integrationTest struct {
	db sdb.SimpleDB

	mutex           sync.Mutex
	domainsToDelete []sdb.Domain  // Protected by mutex
}

func (t *integrationTest) SetUp(i *TestInfo) {
	var err error

	// Open a connection.
	t.db, err = sdb.NewSimpleDB(g_region, g_accessKey)
	AssertEq(nil, err)
}

func (t *integrationTest) ensureDeleted(d sdb.Domain) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.domainsToDelete = append(t.domainsToDelete, d)
}

func (t *integrationTest) TearDown() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Delete each of the domains created during the test.
	for _, d := range t.domainsToDelete {
		ExpectEq(nil, t.db.DeleteDomain(d), "Domain: %s", d.Name())
	}
}

////////////////////////////////////////////////////////////////////////
// Domains
////////////////////////////////////////////////////////////////////////

type DomainsTest struct {
	integrationTest
}

func init() { RegisterTestSuite(&DomainsTest{}) }

func (t *DomainsTest) InvalidAccessKey() {
	// Open a connection with an unknown key ID.
	wrongKey := g_accessKey
	wrongKey.Id += "taco"

	db, err := sdb.NewSimpleDB(g_region, wrongKey)
	AssertEq(nil, err)

	// Attempt to create a domain.
	_, err = db.OpenDomain("some_domain")

	ExpectThat(err, Error(HasSubstr("TODO")))
}

func (t *DomainsTest) DomainsHaveIndependentItems() {
	ExpectEq("TODO", "")
}

func (t *DomainsTest) Delete() {
	ExpectEq("TODO", "")
}

func (t *DomainsTest) OpeningTwiceDoesntDeleteExistingItems() {
	ExpectEq("TODO", "")
}

func (t *DomainsTest) DeleteTwice() {
	ExpectEq("TODO", "")
}

////////////////////////////////////////////////////////////////////////
// Items
////////////////////////////////////////////////////////////////////////

type ItemsTest struct {
	integrationTest
}

func init() { RegisterTestSuite(&ItemsTest{}) }

func (t *ItemsTest) WrongAccessKeySecret() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) InvalidUtf8ItemName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) InvalidUtf8AttributeName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) InvalidUtf8AttributeValue() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) LongItemName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) LongAttributeName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) LongAttributeValue() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) PutThenGet() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchPutThenGet() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchPutThenBatchGet() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) GetForNonExistentItem() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) GetParticularAttributes() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchGetParticularAttributes() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchGetForNonExistentItems() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) GetNonExistentAttributeName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchGetNonExistentAttributeName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) FailedValuePrecondition() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) FailedExistencePrecondition() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) FailedNonExistencePrecondition() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SuccessfulPreconditions() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) DeleteParticularAttributes() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) DeleteAllAttributes() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) BatchDelete() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) InvalidSelectQuery() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectAll() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectItemName() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectCount() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectWithPredicates() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectWithSortOrder() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectWithLimit() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectEmptyResultSet() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) SelectLargeResultSet() {
	ExpectEq("TODO", "")
}
