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
)

func validateDeleteUpdate(u DeleteUpdate) (err error) {
	// Make sure the attribute name is legal.
	if u.Name == "" {
		return fmt.Errorf("Invalid attribute name; names must be non-empty.")
	}

	if err = validateValue(string(u.Name)); err != nil {
		return fmt.Errorf("Invalid attribute name: %v", err)
	}

	// Make sure the attribute value is legal, if it is specified.
	if u.Value != nil {
		if err = validateValue(string(*u.Value)); err != nil {
			return fmt.Errorf("Invalid attribute value: %v", err)
		}
	}

	return nil
}

func validateDeleteUpdates(updates []DeleteUpdate) (err error) {
	for _, u := range updates {
		if err = validateDeleteUpdate(u); err != nil {
			return fmt.Errorf("Invalid update (%v): %v", err, u)
		}
	}

	return nil
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
