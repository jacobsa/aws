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

package main

import (
	"flag"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3"
	"os"
)

var bucketName = flag.String("bucket", "", "Bucket name.")
var region = flag.String("region", "", "Region endpoint server.")

func main() {
	flag.Parse()

	if *bucketName == "" {
		fmt.Println("You must set the -bucket flag.")
		os.Exit(1)
	}

	if *region == "" {
		fmt.Println("You must set the -region flag.")
		os.Exit(1)
	}

	// Read in the access key.
	accessKey := aws.AccessKey{}
	accessKey.Id = readPassword("Access key ID: ")
	accessKey.Secret = readPassword("Access key secret: ")

	// Create a bucket object.
	_, err := s3.OpenBucket(*bucketName, s3.Region(*region), accessKey)
	if err != nil {
		fmt.Println("Opening bucket:", err)
		os.Exit(1)
	}
}

