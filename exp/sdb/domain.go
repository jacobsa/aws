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
	"github.com/jacobsa/aws/exp/sdb/conn"
)

// A precondition for a conditional Put or Delete operation. Preconditions may
// specify a value that an attribute must have or that the attribute must not
// exist.
type Precondition struct {
	// The name of the attribute to be inspected. Attributes with multiple values
	// are not supported.
	Name string

	// If non-nil, the value that the attribute must possess at the time of the
	// update. If nil, the attribute must not exist at the time of the update.
	Value *string
}

// An update to make to a particular attribute as part of a Put request.
type PutUpdate struct {
	// The name of the attribute.
	Name string

	// The value to set for the attribute.
	Value string

	// If true, add this (name, value) attribute rather than replacing any
	// existing attributes with the given name (the default).
	Add bool
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
	// Return the name of this domain.
	Name() string

	// Atomically apply the supplied updates to the attributes of the named item,
	// but only if the supplied precondition holds. If no precondition is
	// desired, pass nil.
	//
	// The length of updates must be in [1, 256].
	PutAttributes(
		item ItemName,
		updates []PutUpdate,
		precondition *Precondition) error

	// Atomically apply updates to multiple items simultaneously.
	//
	// The length of the map must be in [1, 25]. The length of each of its values
	// must be in [1, 256]. An error may be returned if the underlying request to
	// SimpleDB is too large.
	BatchPutAttributes(updateMap map[ItemName][]PutUpdate) error

	// Retrieve a set of attributes for the named item, or all attributes if the
	// attribute name slice is empty.
	//
	// If the named item doesn't exist, the empty set is returned.
	//
	// constistentRead specifies whether completely fresh data is needed or not.
	GetAttributes(
		item ItemName,
		constistentRead bool,
		attrNames []string) (attrs []Attribute, err error)

	// Atomically delete attributes from the named item, but only if the supplied
	// precondition holds. If no precondition is desired, pass nil.
	//
	// If deletes is empty, delete all attributes from the item. Otherwise
	// perform only the deletes is specifies. Deleting a non-existent attribute
	// does not result in an error.
	DeleteAttributes(
		item ItemName,
		deletes []DeleteUpdate,
		precondition *Precondition) error

	// Atomically delete attributes from multiple items simultaneously.
	//
	// If no updates are supplied for a particular item, delete all of its
	// attributes.
	BatchDeleteAttributes(deleteMap map[ItemName][]DeleteUpdate) error
}

func newDomain(name string, c conn.Conn) (Domain, error) {
	return &domain{name, c}, nil
}

type domain struct {
	name string
	c    conn.Conn
}

func (d *domain) Name() string {
	return d.name
}
