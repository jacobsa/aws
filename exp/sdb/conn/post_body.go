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

package conn

import (
	"net/url"
	"sort"
	"strings"
)

type keyValuePair struct {
	Key string
	Val string
}

// Sortable by key.
type keyValueList []keyValuePair

func (l keyValueList) Len() int { return len(l) }
func (l keyValueList) Less(i, j int) bool { return l[i].Key < l[j].Key }
func (l keyValueList) Swap(i, j int) { l[j], l[i] = l[i], l[j] }

// Given a set of request parameters, assemble them into a form usable both as
// a POST body in a request to SimpleDB and as the "canonicalized query string"
// in the SimpleDB signing algorithm.
//
// Reference:
//     http://goo.gl/sRr8w
func assemblePostBody(req Request) string {
	// Make a list of key/value pairs and sort them by key.
	kvPairs := keyValueList{}
	for key, val := range req {
		kvPairs = append(kvPairs, keyValuePair{key, val})
	}

	sort.Sort(kvPairs)

	// Assemble the appropriate parts.
	parts := make([]string, len(kvPairs))
	for i, kvPair := range kvPairs {
		parts[i] = url.QueryEscape(kvPair.Key) + "=" + url.QueryEscape(kvPair.Val)
	}

	return strings.Join(parts, "&")
}
