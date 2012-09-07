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
	"bytes"
	"flag"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3"
	"os"
	"strings"
)

var keyId = flag.String("key_id", "", "Access key ID.")
var bucketName = flag.String("bucket", "", "Bucket name.")
var region = flag.String("region", "", "Region endpoint server.")

func main() {
	flag.Parse()

	if *keyId == "" {
		fmt.Println("You must set the -key_id flag.")
		os.Exit(1)
	}

	if *bucketName == "" {
		fmt.Println("You must set the -bucket flag.")
		os.Exit(1)
	}

	if *region == "" {
		fmt.Println("You must set the -region flag.")
		os.Exit(1)
	}

	// Read in the access key.
	accessKey := aws.AccessKey{Id: *keyId}
	accessKey.Secret = readPassword("Access key secret: ")

	// Open a bucket.
	bucket, err := s3.OpenBucket(*bucketName, s3.Region(*region), accessKey)
	if err != nil {
		fmt.Println("Opening bucket:", err)
		os.Exit(1)
	}

	// Attempt to create an object.
	data := []byte("taco")
	data = append(data, 0x00)
	data = append(data, []byte("burrito")...)

	if err := bucket.StoreObject("some_taco", data); err != nil {
		fmt.Println("StoreObject:", err)
		os.Exit(1)
	}

	// Read the object back.
	dataRead, err := bucket.GetObject("some_taco")
	if err != nil {
		fmt.Println("GetObject:", err)
		os.Exit(1)
	}

	// Make sure the result is identical.
	if !bytes.Equal(data, dataRead) {
		fmt.Printf("Mismatch; %x vs. %x\n", data, dataRead)
		os.Exit(1)
	}

	// Attempt to load a non-existent object. We should get a 404 back.
	_, err = bucket.GetObject("other_name")
	if err == nil || strings.Count(err.Error(), "404") != 1 {
		fmt.Println("Unexpected 404 error:", err)
		os.Exit(1)
	}
}
