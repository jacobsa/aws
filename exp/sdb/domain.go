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
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/exp/sdb/conn"
)

// The name of an item within a SimpleDB domain. Item names must be UTF-8
// strings no longer than 1024 bytes. They must contain only characters that
// are valid in XML 1.0 documents, as defined by Section 2.2 of the XML 1.0
// spec. (Note that this is a more restrictive condition than imposed by
// SimpleDB itself, and is done for the sake of Go's XML 1.0 parser.)
//
// For more info:
//
//     http://goo.gl/Fkjnz
//     http://goo.gl/csem8
//
type ItemName string

// An attribute is a (name, value) pair possessed by an item. Items contain
// sets of attributes; they may contain multiple attributes with the same name,
// but not with the same (name, value) pair.
//
// Attribute names and values share the same restrictions as those on item
// names.
type Attribute struct {
	Name  string
	Value string
}

// A precondition for a conditional Put or Delete operation. Preconditions may
// specify a value that an attribute must have or whether the attribute must
// exist or not.
type Precondition struct {
	// The name of the attribute to be inspected. Attributes with multiple values
	// are not supported.
	Name string

	// If present, the value that the attribute must possess at the time of the
	// update. Must be present iff Exists is not present.
	Value *string

	// If present, whether the attribute must exist at the time of the update.
	// Must be present iff Value is not present.
	Exists *bool
}

// An update to make to a particular attribute as part of a Put request.
type PutUpdate struct {
	// The name of the attribute.
	Name string

	// The value to set for the attribute.
	Value string

	// Whether to replace existing values for the attribute or to simply add a
	// new one.
	Replace bool
}

// An update to make to a particular attribute as part of a Delete request.
type DeleteUpdate struct {
	// The name of the attribute.
	Name string

	// Te requests, the particular value of the attribute to delete if present.
	// Otherwise, all values will be deleted.
	Value *string
}

// A domain represents a named domain within the SimpleDB service. It is a
// collection of named items, each of which possesses a set of attributes.
type Domain interface {
	// Atomically apply the supplied updates to the attributes of the named item,
	// but only if the supplied preconditions hold.
	//
	// The length of updates must be in [1, 256].
	PutAttributes(
		item ItemName,
		updates []PutUpdate,
		preconditions []Precondition) error

	// Atomically apply updates to multiple items simultaneously.
	//
	// The length of updates must be in [1, 25]. The length of each of its values
	// must be in [1, 256]. An error may be returned if the underlying request to
	// SimpleDB is too large.
	BatchPutAttributes(updates map[ItemName][]PutUpdate) error

	// Atomically delete attributes from the named item, but only if the supplied
	// preconditions hold.
	//
	// If deletes is empty, delete all attributes from the item. Otherwise
	// perform only the deletes is specifies. Deleting a non-existent attribute
	// does not result in an error.
	DeleteAttributes(
		item ItemName,
		deletes []DeleteUpdate,
		preconditions []Precondition) error

	// Atomically delete attributes from multiple items simultaneously.
	//
	// If no updates are supplied for a particular item, delete all of its
	// attributes.
	BatchDeleteAttributes(deletes map[ItemName][]DeleteUpdate) error

	// Retrieve a set of attributes for the named item, or all attributes if the
	// attributes slice is empty.
	//
	// If the named item doesn't exist, the empty set is returned.
	//
	// constistentRead specifies whether completely fresh data is needed or not.
	GetAttributes(
		item ItemName,
		constistentRead bool,
		attributes []string) ([]Attribute, error)

	// Retrieve a set of items and their attributes based on a query string.
	//
	// constistentRead specifies whether completely fresh data is needed or not.
	//
	// If the selected result set is too large, a "next token" will be returned.
	// It can be passed to the Select method to resume where the previous result
	// set left off. For the initial query, use nil.
	//
	// For more info:
	//
	//     http://goo.gl/GTsSZ
	//
	Select(
		query string,
		constistentRead bool,
		nextToken []byte) (res map[ItemName][]Attribute, tok []byte, err error)
}

// Return a Domain tied to a given name in a given region. You must have
// previously created the domain in the region, and the supplied access key
// must have access to it.
//
// To easily create a domain, use the Scratchpad app:
//
//     http://goo.gl/C9BMz
//
func OpenDomain(name string, region Region, key aws.AccessKey) (Domain, error)

// As above, but allows injecting a Conn directly.
func newDomain(name string, c conn.Conn) (Domain, error) {
	return &domain{name, c}, nil
}

type domain struct {
	name string
	c conn.Conn
}

func (d *domain) DeleteAttributes(
	item ItemName,
	deletes []DeleteUpdate,
	preconditions []Precondition) error {
	return fmt.Errorf("TODO")
}

func (d *domain) BatchDeleteAttributes(deletes map[ItemName][]DeleteUpdate) error {
	return fmt.Errorf("TODO")
}

func (d *domain) GetAttributes(
	item ItemName,
	constistentRead bool,
	attributes []string) ([]Attribute, error) {
	return nil, fmt.Errorf("TODO")
}

func (d *domain) Select(
		query string,
		constistentRead bool,
		nextToken []byte) (res map[ItemName][]Attribute, tok []byte, err error) {
	return nil, nil, fmt.Errorf("TODO")
}
