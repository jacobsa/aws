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

// Loop indefinitely, storing and fetching data, and reporting on the time it
// takes to do so.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3"
	"github.com/jacobsa/util/password"
	"log"
	"math/rand"
	"time"
)

var g_bucket = flag.String("bucket", "", "The bucket to use for benchmarking.")
var g_region = flag.String("region", "", "The region of the bucket.")
var g_keyId = flag.String("key_id", "", "The AWS access key ID.")

var g_valueSize = flag.Uint("value_size", 1 << 10, "The number of bytes to write at a time.")
var g_parallelism = flag.Uint("parallelism", 1, "The parallelism to use.")

func loop(bucket s3.Bucket) {
	for i := 0; true; {
		key := fmt.Sprintf("value_%016x", rand.Int63())
		data := bytes.Repeat([]byte{byte(i)}, int(*g_valueSize))

		timeBefore := time.Now()
		err := bucket.StoreObject(key, data)
		elapsed := time.Since(timeBefore)

		if err != nil {
			log.Println("StoreObject:", err)
			continue
		}

		log.Printf(
			"Stored %d bytes in %v (%6f.0 bytes/s)\n",
			len(data),
			elapsed,
			float64(len(data))/float64((elapsed/time.Second)),
		)
	}
}

func main() {
	flag.Parse()

	// Set up bare logging output.
	log.SetFlags(0)

	// Sanity-check flags.
	if *g_bucket == "" {
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
	bucket, err := s3.OpenBucket(*g_bucket, s3.Region(*g_region), accessKey)
	if err != nil {
		log.Fatalf("OpenBucket:", err)
	}

	// Kick off several workers.
	for i := uint(0); i < *g_parallelism; i++ {
		go loop(bucket)
	}

	// Never return.
	someChan := make(chan bool)
	<-someChan
}
