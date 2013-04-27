// Copyright 2013 Aaron Jacobs. All Rights Reserved.
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
	"github.com/jacobsa/util/password"
	"log"
	"sync"
)

var g_bucketName = flag.String("bucket", "", "The bucket to use for benchmarking.")
var g_region = flag.String("region", "", "The region of the bucket.")
var g_keyId = flag.String("key_id", "", "The AWS access key ID.")

var g_bucketOnce sync.Once
var g_bucket s3.Bucket

func initBucket() {
	var err error

	// Sanity-check flags.
	if *g_bucketName == "" {
		log.Fatalln("You must set the -bucket flag.")
	}

	if *g_region == "" {
		log.Fatalln("You must set the -region flag.")
	}

	if *g_keyId == "" {
		log.Fatalln("You must set the -key_id flag.")
	}

	// Set up the access key.
	prompt := fmt.Sprintf(
		"Enter secret for AWS access key %s: ",
		*g_keyId,
	)

	accessKey := aws.AccessKey{
		Id: *g_keyId,
		Secret: password.ReadPassword(prompt),
	}

	// Open the bucket.
	g_bucket, err = s3.OpenBucket(*g_bucketName, s3.Region(*g_region), accessKey)
	if err != nil {
		log.Fatalf("OpenBucket:", err)
	}
}

// Return the globally-configured bucket to use for benchmarking.
func getBucket() s3.Bucket {
	g_bucketOnce.Do(initBucket)
	return g_bucket
}
