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

// A request to the SimpleDB server, specified as key/value parameters. These
// include names of actions, action-specific parameters, and authentication
// info.
//
// For example:
//
//     "Action": "PutAttributes",
//     "DomainName": "some_domain",
//     "ItemName": "some_item",
//     "Attribute.1.Name": "color",
//     "Attribute.1.Value": "blue",
//     "AWSAccessKeyId": "0123456",
//     "Signature": "some_signature",
//
type Request map[string]string
