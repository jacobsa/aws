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
	"github.com/jacobsa/aws/exp/sdb/conn"
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
	preconditions []Precondition) (err error) {
	// Make sure the item name is legal.
	if item == "" {
		return fmt.Errorf("Invalid item name; names must be non-empty.")
	}

	if err = validateValue(string(item)); err != nil {
		return fmt.Errorf("Invalid item name: %v", err)
	}

	// Validate deletes.
	if err = validateDeleteUpdates(deletes); err != nil {
		return err
	}

	// Validate preconditions.
	for _, p := range preconditions {
		if err = validatePrecondition(p); err != nil {
			return fmt.Errorf("Invalid precondition (%v): %v", err, p)
		}
	}

	// Assemble an appropriate request.
	req := conn.Request{}
	req["DomainName"] = d.name
	req["ItemName"] = string(item)

	for i, u := range deletes {
		keyPrefix := fmt.Sprintf("Attribute.%d.", i+1)
		req[keyPrefix + "Name"] = u.Name

		if u.Value != nil {
			req[keyPrefix + "Value"] = *u.Value
		}
	}

	for i, p := range preconditions {
		keyPrefix := fmt.Sprintf("Expected.%d.", i+1)
		req[keyPrefix + "Name"] = p.Name

		if p.Value != nil {
			req[keyPrefix + "Value"] = *p.Value
		} else if *p.Exists {
			req[keyPrefix + "Exists"] = "true"
		} else {
			req[keyPrefix + "Exists"] = "false"
		}
	}

	// Call the connection.
	if _, err = d.c.SendRequest(req); err != nil {
		return fmt.Errorf("SendRequest: %v", err)
	}

	return nil
}

func (d *domain) BatchDeleteAttributes(deletes map[ItemName][]DeleteUpdate) error {
	return fmt.Errorf("TODO")
}
