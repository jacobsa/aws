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
)

type ItemName string

// TODO(jacobsa): Comments.
type Attribute struct {
	Name string
	Value string
}

// TODO(jacobsa): Comments.
type Expectation struct {
	Name string
	Value *string
	Exists *bool
}

// TODO(jacobsa): Comments.
type AttributeUpdate struct {
	Name string
	Value *string
	Replace bool
}

// TODO(jacobsa): Comments.
type Domain interface {
	PutAttributes(
		item ItemName,
		updates []AttributeUpdate,
		expectations []Expectation) error

	BatchPutAttributes(updates map[ItemName][]AttributeUpdate) error

	DeleteAttributes(
		item ItemName,
		deletes []AttributeUpdate,
		expectations []Expectation) error

	BatchDeleteAttributes(deletes map[ItemName][]AttributeUpdate) error

	GetAttributes(
		item ItemName,
		constistentRead bool,
		attributes []string) ([]Attribute, error)

	Select(
		query string,
		constistentRead bool,
		nextToken []byte) (res map[ItemName][]Attribute, tok []byte, err error)
}
