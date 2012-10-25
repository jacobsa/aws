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

package s3util

import (
	"fmt"
	"github.com/jacobsa/aws/s3"
)

// List all keys currently contained by the bucket.
func ListAllKeys(bucket s3.Bucket) (keys []string, err error) {
	for {
		var prevKey string
		if len(keys) > 0 {
			prevKey = keys[len(keys)-1]
		}

		var partialKeys []string
		partialKeys, err = bucket.ListKeys(prevKey)
		if err != nil {
			err = fmt.Errorf("ListKeys: %v", err)
			return
		}

		if len(partialKeys) == 0 {
			break
		}

		keys = append(keys, partialKeys...)
	}

	return
}
