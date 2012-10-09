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

func (d *domain) GetAttributes(
	item ItemName,
	constistentRead bool,
	attrNames []string) (attrs []Attribute, err error) {
	// Make sure the item name is legal.
	if item == "" {
		err = fmt.Errorf("Invalid item name; names must be non-empty.")
		return
	}

	if err = validateValue(string(item)); err != nil {
		err = fmt.Errorf("Invalid item name: %v", err)
		return
	}

	// Make sure attribute names are legal.
	for _, name := range attrNames {
		if name == "" {
			err = fmt.Errorf("Invalid attribute name; names must be non-empty.")
			return
		}

		if err = validateValue(name); err != nil {
			err = fmt.Errorf("Invalid attribute name: %v", err)
			return
		}
	}

	// Create an appropriate request.
	//
	// Reference:
	//     http://goo.gl/MmaJA
	req := conn.Request{
		"DomainName": d.name,
		"ItemName": string(item),
	}

	if constistentRead {
		req["ConsistentRead"] = "true"
	}

	for i, name := range attrNames {
		req[fmt.Sprintf("AttributeName.%d", i)] = name
	}

	// Call the connection.
	if _, err = d.c.SendRequest(req); err != nil {
		err = fmt.Errorf("SendRequest: %v", err)
		return
	}

	return nil, fmt.Errorf("TODO")
}
