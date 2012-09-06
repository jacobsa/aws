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

package s3

// Bucket represents an S3 bucket, which is a collection of objects keyed on
// Unicode strings. The UTF-8 encoding of a key must be no more than 1024 bytes
// long.
//
// See here for more info:
//
//     http://goo.gl/Nd63t
//
type Bucket interface {
	// Retrieve data for the object with the given key.
	GetObject(key string) (data []byte, err error)

	// Store the supplied data with the given key, overwriting any previous
	// version.
	StoreObject(key string, data []byte) error
}
