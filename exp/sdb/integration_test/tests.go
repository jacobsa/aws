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
	"fmt"
	"github.com/jacobsa/aws/exp/sdb"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"math/rand"
	"sort"
	"sync"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type integrationTest struct {
	db sdb.SimpleDB
}

func (t *integrationTest) SetUp(i *TestInfo) {
	var err error

	// Open a connection.
	t.db, err = sdb.NewSimpleDB(g_region, g_accessKey)
	AssertEq(nil, err)
}

// Generate an item name likely to be unique.
func (t *integrationTest) makeItemName() sdb.ItemName {
	return sdb.ItemName(fmt.Sprintf("item.%16x", uint64(rand.Int63())))
}

type nameSortedAttrList []sdb.Attribute

func (l nameSortedAttrList) Len() int           { return len(l) }
func (l nameSortedAttrList) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l nameSortedAttrList) Swap(i, j int)      { l[j], l[i] = l[i], l[j] }

func sortByName(attrs []sdb.Attribute) []sdb.Attribute {
	res := make(nameSortedAttrList, len(attrs))
	copy(res, attrs)
	sort.Sort(res)
	return res
}

////////////////////////////////////////////////////////////////////////
// Domains
////////////////////////////////////////////////////////////////////////

var g_domainsTestDb sdb.SimpleDB
var g_domainsTestDomain0 sdb.Domain
var g_domainsTestDomain1 sdb.Domain

type DomainsTest struct {
	integrationTest

	mutex           sync.Mutex
	domainsToDelete []sdb.Domain  // Protected by mutex
}

func init() { RegisterTestSuite(&DomainsTest{}) }

func (t *DomainsTest) SetUpTestSuite() {
	var err error

	// Open a connection.
	g_domainsTestDb, err = sdb.NewSimpleDB(g_region, g_accessKey)
	if err != nil {
		panic(err)
	}

	// Create domain 0.
	g_domainsTestDomain0, err = g_domainsTestDb.OpenDomain("DomainsTest.domain0")
	if err != nil {
		panic(err)
	}

	// Create domain 1.
	g_domainsTestDomain1, err = g_domainsTestDb.OpenDomain("DomainsTest.domain1")
	if err != nil {
		panic(err)
	}
}

func (t *DomainsTest) TearDownTestSuite() {
	// Delete both domains.
	if err := g_domainsTestDb.DeleteDomain(g_domainsTestDomain0); err != nil {
		panic(err)
	}

	if err := g_domainsTestDb.DeleteDomain(g_domainsTestDomain1); err != nil {
		panic(err)
	}

	// Clear variables.
	g_domainsTestDb = nil
	g_domainsTestDomain0 = nil
	g_domainsTestDomain1 = nil
}

func (t *DomainsTest) TearDown() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Delete each of the domains created during the test.
	for _, d := range t.domainsToDelete {
		ExpectEq(nil, t.db.DeleteDomain(d), "Domain: %s", d.Name())
	}
}

func (t *DomainsTest) InvalidAccessKey() {
	// Open a connection with an unknown key ID.
	wrongKey := g_accessKey
	wrongKey.Id += "taco"

	db, err := sdb.NewSimpleDB(g_region, wrongKey)
	AssertEq(nil, err)

	// Attempt to create a domain.
	_, err = db.OpenDomain("some_domain")

	ExpectThat(err, Error(HasSubstr("403")))
	ExpectThat(err, Error(HasSubstr("Key Id")))
	ExpectThat(err, Error(HasSubstr("exist")))
}

func (t *DomainsTest) SeparatelyNamedDomainsHaveIndependentItems() {
	var err error

	// Set up an item in the first domain.
	itemName := t.makeItemName()
	err = g_domainsTestDomain0.PutAttributes(
		itemName,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "enchilada", Value: "queso"},
		},
		[]sdb.Precondition{},
	)

	AssertEq(nil, err)

	// Get attributes for the same name in the other domain. There should be
	// none.
	attrs, err := g_domainsTestDomain1.GetAttributes(itemName, true, []string{})
	AssertEq(nil, err)

	ExpectThat(attrs, ElementsAre())
}

func (t *DomainsTest) IdenticallyNamedDomainsHaveIdenticalItems() {
	var err error

	// Set up an item in the first domain.
	itemName := t.makeItemName()
	err = g_domainsTestDomain0.PutAttributes(
		itemName,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "enchilada", Value: "queso"},
		},
		[]sdb.Precondition{},
	)

	AssertEq(nil, err)

	// Get attributes for the same name in another domain object opened with the
	// same name.
	domain1, err := t.db.OpenDomain(g_domainsTestDomain0.Name())
	AssertEq(nil, err)

	attrs, err := domain1.GetAttributes(itemName, true, []string{})
	AssertEq(nil, err)

	ExpectThat(
		attrs,
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "enchilada", Value: "queso"}),
		),
	)
}

func (t *DomainsTest) Delete() {
	var err error
	domainName := "DomainsTest.Delete"

	// Create a domain, then delete it.
	domain, err := t.db.OpenDomain(domainName)
	AssertEq(nil, err)

	err = t.db.DeleteDomain(domain)
	AssertEq(nil, err)

	// Delete again; nothing should go wrong.
	err = t.db.DeleteDomain(domain)
	AssertEq(nil, err)

	// Attempt to write to the domain.
	err = domain.PutAttributes(
		"some_item",
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "bar"},
		},
		[]sdb.Precondition{},
	)

	ExpectThat(err, Error(HasSubstr("NoSuchDomain")))
	ExpectThat(err, Error(HasSubstr("domain")))
	ExpectThat(err, Error(HasSubstr("exist")))
}

////////////////////////////////////////////////////////////////////////
// Items
////////////////////////////////////////////////////////////////////////

var g_itemsTestDb     sdb.SimpleDB
var g_itemsTestDomain sdb.Domain

type ItemsTest struct {
	integrationTest
}

func init() { RegisterTestSuite(&ItemsTest{}) }

func (t *ItemsTest) SetUpTestSuite() {
	var err error

	// Open a connection.
	g_itemsTestDb, err = sdb.NewSimpleDB(g_region, g_accessKey)
	if err != nil {
		panic(err)
	}

	// Create a domain.
	g_itemsTestDomain, err = g_itemsTestDb.OpenDomain("ItemsTest.domain")
	if err != nil {
		panic(err)
	}
}

func (t *ItemsTest) TearDownTestSuite() {
	// Delete the domain.
	if err := g_itemsTestDb.DeleteDomain(g_itemsTestDomain); err != nil {
		panic(err)
	}

	// Clear variables.
	g_itemsTestDb = nil
	g_itemsTestDomain = nil
}

func (t *ItemsTest) PutThenGet() {
	var err error
	item := t.makeItemName()

	// Put
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		[]sdb.Precondition{},
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
			DeepEquals(sdb.Attribute{Name: "baz", Value: "enchilada"}),
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}

func (t *ItemsTest) BatchPutThenGet() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		map[sdb.ItemName][]sdb.PutUpdate{
			item0: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "foo", Value: "taco"},
				sdb.PutUpdate{Name: "bar", Value: "burrito"},
			},
			item1: []sdb.PutUpdate{
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
			},
		},
	)

	AssertEq(nil, err)

	// Get for item 0
	attrs, err := g_itemsTestDomain.GetAttributes(item0, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)

	// Get for item 1
	attrs, err = g_itemsTestDomain.GetAttributes(item1, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "baz", Value: "enchilada"}),
		),
	)
}

func (t *ItemsTest) GetForNonExistentItem() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) GetParticularAttributes() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) GetNonExistentAttributeName() {
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
