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

type batchPutPair struct {
	Item    ItemName
	Updates []PutUpdate
}

type batchPutPairList []batchPutPair

func (l batchPutPairList) Len() int           { return len(l) }
func (l batchPutPairList) Less(i, j int) bool { return l[i].Item < l[j].Item }
func (l batchPutPairList) Swap(i, j int)      { l[j], l[i] = l[i], l[j] }

// Return the elements of the map sorted by item name.
func getSortedPutPairs(updateMap map[ItemName][]PutUpdate) batchPutPairList {
	res := batchPutPairList{}
	for item, updates := range updateMap {
		res = append(res, batchPutPair{item, updates})
	}

	sort.Sort(res)
	return res
}

func validatePutUpdate(u PutUpdate) (err error) {
	// Make sure the attribute name is legal.
	if u.Name == "" {
		return fmt.Errorf("Invalid attribute name; names must be non-empty.")
	}

	if err = validateValue(string(u.Name)); err != nil {
		return fmt.Errorf("Invalid attribute name: %v", err)
	}

	// Make sure the attribute value is legal.
	if err = validateValue(string(u.Value)); err != nil {
		return fmt.Errorf("Invalid attribute value: %v", err)
	}

	return nil
}

func validatePutUpdates(updates []PutUpdate) (err error) {
	numUpdates := len(updates)
	if numUpdates == 0 || numUpdates > 256 {
		return fmt.Errorf("Illegal number of updates: %d", numUpdates)
	}

	for _, u := range updates {
		if err = validatePutUpdate(u); err != nil {
			return fmt.Errorf("Invalid update (%v): %v", err, u)
		}
	}

	return nil
}

func (d *domain) PutAttributes(
	item ItemName,
	updates []PutUpdate,
	preconditions []Precondition) (err error) {
	// Make sure the item name is legal.
	if item == "" {
		return fmt.Errorf("Invalid item name; names must be non-empty.")
	}

	if err = validateValue(string(item)); err != nil {
		return fmt.Errorf("Invalid item name: %v", err)
	}

	// Validate updates.
	if err = validatePutUpdates(updates); err != nil {
		return err
	}

	// Validate preconditions.
	for _, p := range preconditions {
		if err = validatePrecondition(p); err != nil {
			return fmt.Errorf("Invalid precondition (%v): %v", err, p)
		}
	}

	// Assemble an appropriate request.
	req := conn.Request{
		"Action": "PutAttributes",
		"Version": apiVersion,
		"DomainName": d.name,
		"ItemName": string(item),
	}

	for i, u := range updates {
		keyPrefix := fmt.Sprintf("Attribute.%d.", i+1)
		req[keyPrefix+"Name"] = u.Name
		req[keyPrefix+"Value"] = u.Value

		if !u.Add {
			req[keyPrefix+"Replace"] = "true"
		}
	}

	for i, p := range preconditions {
		keyPrefix := fmt.Sprintf("Expected.%d.", i+1)
		req[keyPrefix+"Name"] = p.Name

		if p.Value != nil {
			req[keyPrefix+"Value"] = *p.Value
		} else if *p.Exists {
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

func (d *domain) BatchPutAttributes(updateMap map[ItemName][]PutUpdate) (err error) {
	// Make sure the size of the request is legal.
	numItems := len(updateMap)
	if numItems == 0 || numItems > 25 {
		return fmt.Errorf("Illegal number of items: %d", numItems)
	}

	// Make sure each item name and set of updates is legal.
	for item, updates := range updateMap {
		if item == "" {
			return fmt.Errorf("Invalid item name; names must be non-empty.")
		}

		if err = validateValue(string(item)); err != nil {
			return fmt.Errorf("Invalid item name: %v", err)
		}

		if err = validatePutUpdates(updates); err != nil {
			return fmt.Errorf("Updates for item %s: %v", item, err)
		}
	}

	// Build a request.
	req := conn.Request{
		"Action": "BatchPutAttributes",
		"Version": apiVersion,
		"DomainName": d.name,
	}

	pairs := getSortedPutPairs(updateMap)
	for i, pair := range pairs {
		itemPrefix := fmt.Sprintf("Item.%d.", i+1)
		req[itemPrefix+"ItemName"] = string(pair.Item)

		for j, u := range pair.Updates {
			updatePrefix := fmt.Sprintf("%sAttribute.%d.", itemPrefix, j+1)
			req[updatePrefix+"Name"] = u.Name
			req[updatePrefix+"Value"] = u.Value

			if !u.Add {
				req[updatePrefix+"Replace"] = "true"
			}
		}
	}

	// Call the connection.
	if _, err = d.c.SendRequest(req); err != nil {
		return fmt.Errorf("SendRequest: %v", err)
	}

	return nil
}
