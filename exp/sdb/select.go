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
	"encoding/xml"
	"fmt"
	"github.com/jacobsa/aws/exp/sdb/conn"
)

type selectedItem struct {
	Name       ItemName
	Attributes []Attribute `xml:"Attribute"`
}

type selectResult struct {
	Items     []selectedItem `xml:"Item"`
	NextToken []byte
}

type selectResponse struct {
	SelectResult selectResult
}

func parseSelectResponse(resp []byte) (
	results []SelectedItem,
	tok []byte,
	err error) {
	responseStruct := &selectResponse{}
	if err = xml.Unmarshal(resp, responseStruct); err != nil {
		err = fmt.Errorf("Invalid response from server (%v): %s", err, resp)
		return
	}

	for _, item := range responseStruct.SelectResult.Items {
		selectedItem := SelectedItem{item.Name, item.Attributes}
		results = append(results, selectedItem)
	}

	tok = responseStruct.SelectResult.NextToken
	return
}

func (db *simpleDB) Select(
	query string,
	constistentRead bool,
	nextToken []byte) (results []SelectedItem, tok []byte, err error) {
	// Create an appropriate request.
	//
	// Reference:
	//     http://goo.gl/GTsSZ
	req := conn.Request{
		"Action":           "Select",
		"Version":          apiVersion,
		"SelectExpression": query,
	}

	if constistentRead {
		req["ConsistentRead"] = "true"
	}

	if nextToken != nil {
		req["NextToken"] = string(nextToken)
	}

	// Call the connection.
	resp, err := db.c.SendRequest(req)
	if err != nil {
		err = fmt.Errorf("SendRequest: %v", err)
		return
	}

	return parseSelectResponse(resp)
}
