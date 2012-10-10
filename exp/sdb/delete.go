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
	"sort"
)

type batchDeletePair struct {
	Item    ItemName
	Updates []DeleteUpdate
}

type batchDeletePairList []batchDeletePair

func (l batchDeletePairList) Len() int           { return len(l) }
func (l batchDeletePairList) Less(i, j int) bool { return l[i].Item < l[j].Item }
func (l batchDeletePairList) Swap(i, j int)      { l[j], l[i] = l[i], l[j] }

// Return the elements of the map sorted by item name.
func getSortedDeletePairs(deleteMap map[ItemName][]DeleteUpdate) batchDeletePairList {
	res := batchDeletePairList{}
	for item, updates := range deleteMap {
		res = append(res, batchDeletePair{item, updates})
	}

	sort.Sort(res)
	return res
}

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
	numUpdates := len(updates)
	if numUpdates > 256 {
		return fmt.Errorf("Illegal number of updates: %d", numUpdates)
	}

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
	precondition *Precondition) (err error) {
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

	// Validate the precondition, if any.
	if precondition != nil {
		if err = validatePrecondition(*precondition); err != nil {
			return fmt.Errorf("Invalid precondition (%v): %v", err, *precondition)
		}
	}

	// Assemble an appropriate request.
	req := conn.Request{
		"Action":     "DeleteAttributes",
		"Version":    apiVersion,
		"DomainName": d.name,
		"ItemName":   string(item),
	}

	for i, u := range deletes {
		keyPrefix := fmt.Sprintf("Attribute.%d.", i+1)
		req[keyPrefix+"Name"] = u.Name

		if u.Value != nil {
			req[keyPrefix+"Value"] = *u.Value
		}
	}

	if precondition != nil {
		keyPrefix := "Expected.1."
		req[keyPrefix+"Name"] = precondition.Name

		if precondition.Value != nil {
			req[keyPrefix+"Value"] = *precondition.Value
		} else if *precondition.Exists {
			req[keyPrefix+"Exists"] = "true"
		} else {
			req[keyPrefix+"Exists"] = "false"
		}
	}

	// Call the connection.
	if _, err = d.c.SendRequest(req); err != nil {
		return fmt.Errorf("SendRequest: %v", err)
	}

	return nil
}

func (d *domain) BatchDeleteAttributes(
	deleteMap map[ItemName][]DeleteUpdate) (err error) {
	// Make sure the size of the request is legal.
	numItems := len(deleteMap)
	if numItems == 0 || numItems > 25 {
		return fmt.Errorf("Illegal number of items: %d", numItems)
	}

	// Make sure each item name and set of updates is legal.
	for item, updates := range deleteMap {
		if item == "" {
			return fmt.Errorf("Invalid item name; names must be non-empty.")
		}

		if err = validateValue(string(item)); err != nil {
			return fmt.Errorf("Invalid item name: %v", err)
		}

		if err = validateDeleteUpdates(updates); err != nil {
			return fmt.Errorf("Updates for item %s: %v", item, err)
		}
	}

	// Build a request.
	req := conn.Request{
		"Action":     "BatchDeleteAttributes",
		"Version":    apiVersion,
		"DomainName": d.name,
	}

	pairs := getSortedDeletePairs(deleteMap)
	for i, pair := range pairs {
		itemPrefix := fmt.Sprintf("Item.%d.", i+1)
		req[itemPrefix+"ItemName"] = string(pair.Item)

		for j, u := range pair.Updates {
			updatePrefix := fmt.Sprintf("%sAttribute.%d.", itemPrefix, j+1)
			req[updatePrefix+"Name"] = u.Name

			if u.Value != nil {
				req[updatePrefix+"Value"] = *u.Value
			}
		}
	}

	// Call the connection.
	if _, err = d.c.SendRequest(req); err != nil {
		return fmt.Errorf("SendRequest: %v", err)
	}

	return nil
}
