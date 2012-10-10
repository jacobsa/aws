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

	mutex sync.Mutex
	deleteRequest sdb.BatchDeleteMap  // Protected by mutex
}

func (t *integrationTest) SetUp(i *TestInfo) {
	var err error

	// Set up the record of what item names to delete.
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.deleteRequest = sdb.BatchDeleteMap{}

	// Open a connection.
	t.db, err = sdb.NewSimpleDB(g_region, g_accessKey)
	AssertEq(nil, err)
}

func (t *integrationTest) TearDown() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Delete items, if appropriate.
	if len(t.deleteRequest) > 0 {
		AssertEq(nil, g_itemsTestDomain.BatchDeleteAttributes(t.deleteRequest))
	}
}

// Generate an item name likely to be unique, and make sure it is later
// deleted.
func (t *integrationTest) makeItemName() sdb.ItemName {
	name := sdb.ItemName(fmt.Sprintf("item.%16x", uint64(rand.Int63())))

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.deleteRequest[name] = nil

	return name
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

func makeStrPtr(s string) *string { return &s }

////////////////////////////////////////////////////////////////////////
// Domains
////////////////////////////////////////////////////////////////////////

var g_domainsTestDb sdb.SimpleDB
var g_domainsTestDomain0 sdb.Domain
var g_domainsTestDomain1 sdb.Domain

type DomainsTest struct {
	integrationTest

	mutex           sync.Mutex
	domainsToDelete []sdb.Domain // Protected by mutex
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
		nil,
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
		nil,
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
		nil,
	)

	ExpectThat(err, Error(HasSubstr("NoSuchDomain")))
	ExpectThat(err, Error(HasSubstr("domain")))
	ExpectThat(err, Error(HasSubstr("exist")))
}

////////////////////////////////////////////////////////////////////////
// Items
////////////////////////////////////////////////////////////////////////

var g_itemsTestDb sdb.SimpleDB
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
		nil,
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

func (t *ItemsTest) PutThenAddAndReplace() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "queso", Add: false},
			sdb.PutUpdate{Name: "bar", Value: "burrito", Add: true}, // Same as first
			sdb.PutUpdate{Name: "bar", Value: "carnitas", Add: true},
			sdb.PutUpdate{Name: "baz", Value: "enchilada", Add: false}, // Same as first
		},
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
			DeepEquals(sdb.Attribute{Name: "bar", Value: "carnitas"}),
			DeepEquals(sdb.Attribute{Name: "baz", Value: "enchilada"}),
			DeepEquals(sdb.Attribute{Name: "foo", Value: "queso"}),
		),
	)
}

func (t *ItemsTest) PutThenAddThenReplace() {
	var err error
	item := t.makeItemName()

	// Create the first value for an attribute.
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Add two more values.
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "burrito", Add: true},
			sdb.PutUpdate{Name: "foo", Value: "enchilada", Add: true},
		},
		nil,
	)

	AssertEq(nil, err)

	// Replace all three.
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "queso"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "queso"}),
		),
	)
}

func (t *ItemsTest) BatchPutThenGet() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
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
	var err error
	itemName := t.makeItemName()

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(itemName, true, nil)

	AssertEq(nil, err)
	ExpectThat(attrs, ElementsAre())
}

func (t *ItemsTest) GetOneParticularAttribute() {
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
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, []string{"bar"})

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
		),
	)
}

func (t *ItemsTest) GetTwoParticularAttributes() {
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
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, []string{"foo", "baz"})

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "baz", Value: "enchilada"}),
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}

func (t *ItemsTest) GetNonExistentAttributeName() {
	var err error
	item := t.makeItemName()

	// Put
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, []string{"foo", "baz"})

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}

func (t *ItemsTest) DeleteParticularAttributes() {
	var err error
	item := t.makeItemName()

	// Put
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito", Add: true},
			sdb.PutUpdate{Name: "bar", Value: "enchilada", Add: true},
			sdb.PutUpdate{Name: "baz", Value: "queso", Add: true},
			sdb.PutUpdate{Name: "baz", Value: "carnitas", Add: true},
		},
		nil,
	)

	AssertEq(nil, err)

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{
			sdb.DeleteUpdate{Name: "foo"},
			sdb.DeleteUpdate{Name: "bar", Value: makeStrPtr("enchilada")},
			sdb.DeleteUpdate{Name: "baz"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
		),
	)
}

func (t *ItemsTest) DeleteAllAttributes() {
	var err error
	item := t.makeItemName()

	// Put
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito", Add: true},
			sdb.PutUpdate{Name: "bar", Value: "enchilada", Add: true},
		},
		nil,
	)

	AssertEq(nil, err)

	// Delete
	err = g_itemsTestDomain.DeleteAttributes( item, nil, nil)
	AssertEq(nil, err)

	// Get
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(attrs, ElementsAre())
}

func (t *ItemsTest) BatchDelete() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
			item0: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "foo", Value: "taco"},
				sdb.PutUpdate{Name: "bar", Value: "burrito", Add: true},
				sdb.PutUpdate{Name: "bar", Value: "carnitas", Add: true},
			},
			item1: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "baz", Value: "enchilada"},
			},
		},
	)

	// Batch delete
	err = g_itemsTestDomain.BatchDeleteAttributes(
		sdb.BatchDeleteMap{
			item0: []sdb.DeleteUpdate{
				sdb.DeleteUpdate{Name: "foo"},
				sdb.DeleteUpdate{Name: "bar", Value: makeStrPtr("carnitas")},
			},
			item1: nil,
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
		),
	)

	// Get for item 1
	attrs, err = g_itemsTestDomain.GetAttributes(item1, true, nil)

	AssertEq(nil, err)
	ExpectThat(attrs, ElementsAre())
}

func (t *ItemsTest) InvalidSelectQuery() {
	var err error

	// Select
	_, _, err = g_itemsTestDb.Select(
		"select foo bar baz",
		true,
		nil,
	)

	ExpectThat(err, Error(HasSubstr("400")))
	ExpectThat(err, Error(HasSubstr("InvalidQueryExpression")))
	ExpectThat(err, Error(HasSubstr("syntax")))
}

func (t *ItemsTest) SelectAll() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
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

	// Select
	query := fmt.Sprintf(
		"select * from `%s`",
		g_itemsTestDomain.Name())

	results, tok, err := g_itemsTestDb.Select( query, true, nil)

	AssertEq(nil, err)
	ExpectEq(nil, tok)

	AssertEq(2, len(results), "Results: %v", results)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item0,
					Attributes: []sdb.Attribute{
						sdb.Attribute{Name: "bar", Value: "burrito"},
						sdb.Attribute{Name: "foo", Value: "taco"},
					},
				},
			),
		),
	)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item1,
					Attributes: []sdb.Attribute{
						sdb.Attribute{Name: "baz", Value: "enchilada"},
					},
				},
			),
		),
	)
}

func (t *ItemsTest) SelectItemName() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
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

	// Select
	query := fmt.Sprintf(
		"select itemName() from `%s`",
		g_itemsTestDomain.Name())

	results, tok, err := g_itemsTestDb.Select( query, true, nil)

	AssertEq(nil, err)
	ExpectEq(nil, tok)

	AssertEq(2, len(results), "Results: %v", results)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item0,
				},
			),
		),
	)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item1,
				},
			),
		),
	)
}

func (t *ItemsTest) SelectCount() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
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

	// Select
	query := fmt.Sprintf(
		"select count(*) from `%s`",
		g_itemsTestDomain.Name())

	results, tok, err := g_itemsTestDb.Select( query, true, nil)

	AssertEq(nil, err)
	ExpectEq(nil, tok)

	AssertEq(1, len(results), "Results: %v", results)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: "Domain",
					Attributes: []sdb.Attribute{
						sdb.Attribute{Name: "Count", Value: "2"},
					},
				},
			),
		),
	)
}

func (t *ItemsTest) SelectWithPredicatesAndParticularAttributes() {
	var err error
	item0 := t.makeItemName()
	item1 := t.makeItemName()
	item2 := t.makeItemName()

	// Batch put
	err = g_itemsTestDomain.BatchPutAttributes(
		sdb.BatchPutMap{
			item0: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "foo", Value: "017"},
				sdb.PutUpdate{Name: "bar", Value: "taco"},
			},
			item1: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "foo", Value: "013"},
				sdb.PutUpdate{Name: "bar", Value: "burrito"},
			},
			item2: []sdb.PutUpdate{
				sdb.PutUpdate{Name: "foo", Value: "031"},
				sdb.PutUpdate{Name: "bar", Value: "enchilada"},
			},
		},
	)

	AssertEq(nil, err)

	// Select
	query := fmt.Sprintf(
		"select bar from `%s` where foo > '013'",
		g_itemsTestDomain.Name())

	results, tok, err := g_itemsTestDb.Select( query, true, nil)

	AssertEq(nil, err)
	ExpectEq(nil, tok)

	AssertEq(2, len(results), "Results: %v", results)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item0,
					Attributes: []sdb.Attribute{
						sdb.Attribute{Name: "bar", Value: "taco"},
					},
				},
			),
		),
	)

	ExpectThat(
		results,
		Contains(
			DeepEquals(
				sdb.SelectedItem{
					Name: item2,
					Attributes: []sdb.Attribute{
						sdb.Attribute{Name: "bar", Value: "enchilada"},
					},
				},
			),
		),
	)
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

func (t *ItemsTest) ItemNamesAreCaseSensitive() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) AttributeNamesAreCaseSensitive() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) AttributeValuesAreCaseSensitive() {
	ExpectEq("TODO", "")
}

func (t *ItemsTest) FailedValuePrecondition() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "blahblah"},
			sdb.PutUpdate{Name: "qux", Value: "queso"},
		},
		&sdb.Precondition{Name: "bar", Value: makeStrPtr("asdf")},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("ConditionalCheckFailed")))
	ExpectThat(err, Error(HasSubstr("bar")))
	ExpectThat(err, Error(HasSubstr("burrito")))
	ExpectThat(err, Error(HasSubstr("expected")))
	ExpectThat(err, Error(HasSubstr("asdf")))

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{},
		&sdb.Precondition{Name: "bar", Value: makeStrPtr("asdf")},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("ConditionalCheckFailed")))
	ExpectThat(err, Error(HasSubstr("bar")))
	ExpectThat(err, Error(HasSubstr("burrito")))
	ExpectThat(err, Error(HasSubstr("expected")))
	ExpectThat(err, Error(HasSubstr("asdf")))

	// Get -- neither the second put nor the delete should have taken effect.
	attrs, err := g_itemsTestDomain.GetAttributes(
		item,
		true,
		[]string{"foo", "qux"},
	)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}

func (t *ItemsTest) FailedNonExistencePrecondition() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "blahblah"},
			sdb.PutUpdate{Name: "qux", Value: "queso"},
		},
		&sdb.Precondition{Name: "bar"},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("ConditionalCheckFailed")))
	ExpectThat(err, Error(HasSubstr("bar")))
	ExpectThat(err, Error(HasSubstr("exists")))

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{},
		&sdb.Precondition{Name: "bar"},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("ConditionalCheckFailed")))
	ExpectThat(err, Error(HasSubstr("bar")))
	ExpectThat(err, Error(HasSubstr("exists")))

	// Get -- neither the second put nor the delete should have taken effect.
	attrs, err := g_itemsTestDomain.GetAttributes(
		item,
		true,
		[]string{"foo", "qux"},
	)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}

func (t *ItemsTest) SuccessfulValuePrecondition() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito"},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "queso"},
		},
		&sdb.Precondition{Name: "bar", Value: makeStrPtr("burrito")},
	)

	AssertEq(nil, err)

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{
			sdb.DeleteUpdate{Name: "baz"},
		},
		&sdb.Precondition{Name: "foo", Value: makeStrPtr("queso")},
	)

	AssertEq(nil, err)

	// Get -- both the second put and the delete should have taken effect.
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "bar", Value: "burrito"}),
			DeepEquals(sdb.Attribute{Name: "foo", Value: "queso"}),
		),
	)
}

func (t *ItemsTest) SuccessfulNonExistencePrecondition() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "queso"},
		},
		&sdb.Precondition{Name: "bar", Value: nil},
	)

	AssertEq(nil, err)

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{
			sdb.DeleteUpdate{Name: "baz"},
		},
		&sdb.Precondition{Name: "bar", Value: nil},
	)

	AssertEq(nil, err)

	// Get -- both the second put and the delete should have taken effect.
	attrs, err := g_itemsTestDomain.GetAttributes(item, true, nil)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "queso"}),
		),
	)
}

func (t *ItemsTest) PreconditionWithMultiValuedAttribute() {
	var err error
	item := t.makeItemName()

	// Put (first call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "taco"},
			sdb.PutUpdate{Name: "bar", Value: "burrito", Add: true},
			sdb.PutUpdate{Name: "bar", Value: "carnitas", Add: true},
			sdb.PutUpdate{Name: "baz", Value: "enchilada"},
		},
		nil,
	)

	AssertEq(nil, err)

	// Put (second call)
	err = g_itemsTestDomain.PutAttributes(
		item,
		[]sdb.PutUpdate{
			sdb.PutUpdate{Name: "foo", Value: "queso"},
		},
		&sdb.Precondition{Name: "bar", Value: makeStrPtr("burrito")},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("MultiValuedAttribute")))
	ExpectThat(err, Error(HasSubstr("bar")))

	// Delete
	err = g_itemsTestDomain.DeleteAttributes(
		item,
		[]sdb.DeleteUpdate{},
		&sdb.Precondition{Name: "bar", Value: makeStrPtr("burrito")},
	)

	ExpectThat(err, Error(HasSubstr("409")))
	ExpectThat(err, Error(HasSubstr("MultiValuedAttribute")))
	ExpectThat(err, Error(HasSubstr("bar")))

	// Get -- neither the second put nor the delete should have taken effect.
	attrs, err := g_itemsTestDomain.GetAttributes(
		item,
		true,
		[]string{"foo", "qux"},
	)

	AssertEq(nil, err)
	ExpectThat(
		sortByName(attrs),
		ElementsAre(
			DeepEquals(sdb.Attribute{Name: "foo", Value: "taco"}),
		),
	)
}
