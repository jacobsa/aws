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
	"github.com/jacobsa/aws/s3"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

func storeData(
	bucket s3.Bucket,
	dataSize uint) (key string, err error) {
	// Choose a random key and create some data.
	key = fmt.Sprintf("value_%016x", rand.Int63())
	data := bytes.Repeat([]byte{byte(rand.Uint32())}, int(dataSize))

	// Store the object.
	err = bucket.StoreObject(key, data)
	return
}

////////////////////////////////////////////////////////////////////////
// Latency
////////////////////////////////////////////////////////////////////////

func measureLatency(bucket s3.Bucket) (avg time.Duration, err error) {
	// Average over several runs until we've taken at least this much time.
	const minDuration time.Duration = 4 * time.Second

	var total time.Duration
	var numRuns int

	for ; total < minDuration; numRuns++ {
		timeBefore := time.Now()
		if _, err = storeData(bucket, 1); err != nil {
			return
		}

		total += time.Since(timeBefore)
	}

	avg = time.Duration(float64(total) / float64(numRuns))
	return
}

////////////////////////////////////////////////////////////////////////
// Downstream
////////////////////////////////////////////////////////////////////////

func measureDownstreamBandwidth_SingleRun(
	bucket s3.Bucket,
	keys []string,
	parallelism uint) (bytesPerSecond float64, err error) {
	// Time the whole process.
	timeBefore := time.Now()

	// Start several workers.
	errs := make([]error, int(parallelism))
	done := make(chan bool)

	var totalLoaded uint64
	for i := 0; i < int(parallelism); i++ {
		go func(i int) {
			// Load data from one of the keys.
			data, err := bucket.GetObject(keys[i%len(keys)])
			if err != nil {
				errs[i] = err
			}

			atomic.AddUint64(&totalLoaded, uint64(len(data)))
			done <- true
		}(i)
	}

	// Wait for all of the workers.
	for i := 0; i < int(parallelism); i++ {
		<-done
	}

	elapsed := time.Since(timeBefore)

	// Did any of the workers fail?
	for i := 0; i < int(parallelism); i++ {
		if errs[i] != nil {
			err = errs[i]
			return
		}
	}

	// Estimate the bandwidth.
	bytesTransferred := float64(totalLoaded)
	secondsElapsed := float64(elapsed) / float64(time.Second)

	bytesPerSecond = bytesTransferred / secondsElapsed
	return
}

func measureDownstreamBandwidth(
	bucket s3.Bucket,
	dataSize uint,
	parallelism uint) (bytesPerSecond float64, err error) {
	// Store several objects with the given data size.
	const minUploadDuration time.Duration = 4 * time.Second
	keys := []string{}
	for timeBefore := time.Now(); time.Since(timeBefore) < minUploadDuration; {
		var key string
		key, err = storeData(bucket, dataSize)
		if err != nil {
			return
		}

		keys = append(keys, key)
	}

	// Average over several runs until we've taken at least this much time.
	const minDuration time.Duration = 5 * time.Second

	var bandwidthTotal float64
	var numRuns int

	for timeBefore := time.Now(); time.Since(timeBefore) < minDuration; numRuns++ {
		var singleResult float64
		singleResult, err = measureDownstreamBandwidth_SingleRun(
			bucket,
			keys,
			parallelism,
		)

		if err != nil {
			return
		}

		bandwidthTotal += singleResult
	}

	bytesPerSecond = bandwidthTotal / float64(numRuns)
	return
}

////////////////////////////////////////////////////////////////////////
// Upstream
////////////////////////////////////////////////////////////////////////

func measureUpstreamBandwidth_SingleRun(
	bucket s3.Bucket,
	dataSize uint,
	parallelism uint) (bytesPerSecond float64, err error) {
	// Time the whole process.
	timeBefore := time.Now()

	// Start several workers.
	errs := make([]error, int(parallelism))
	done := make(chan bool)

	for i := 0; i < int(parallelism); i++ {
		go func(i int) {
			_, errs[i] = storeData(bucket, dataSize)
			done <- true
		}(i)
	}

	// Wait for all of the workers.
	for i := 0; i < int(parallelism); i++ {
		<-done
	}

	elapsed := time.Since(timeBefore)

	// Did any of the workers fail?
	for i := 0; i < int(parallelism); i++ {
		if errs[i] != nil {
			err = errs[i]
			return
		}
	}

	// Estimate the bandwidth.
	bytesTransferred := float64(dataSize) * float64(parallelism)
	secondsElapsed := float64(elapsed) / float64(time.Second)

	bytesPerSecond = bytesTransferred / secondsElapsed
	return
}

func measureUpstreamBandwidth(
	bucket s3.Bucket,
	dataSize uint,
	parallelism uint) (bytesPerSecond float64, err error) {
	// Average over several runs until we've taken at least this much time.
	const minDuration time.Duration = 5 * time.Second

	var bandwidthTotal float64
	var numRuns int

	for timeBefore := time.Now(); time.Since(timeBefore) < minDuration; numRuns++ {
		var singleResult float64
		singleResult, err = measureUpstreamBandwidth_SingleRun(
			bucket,
			dataSize,
			parallelism)

		if err != nil {
			return
		}

		bandwidthTotal += singleResult
	}

	bytesPerSecond = bandwidthTotal / float64(numRuns)
	return
}

////////////////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////////////////

func formattedFloatOrError(
	v float64,
	err error) string {
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%6f.0", v)
}

func printHeading(testName string) {
	fmt.Printf("------------------------------------\n")
	fmt.Println(testName)
	fmt.Println()
}

func formatBytes(n uint64) string {
	type exponentAndSuffix struct {
		exponent uint
		suffix   string
	}

	suffixes := []exponentAndSuffix{
		{0, "bytes"},
		{10, "KiB"},
		{20, "MiB"},
		{30, "GiB"},
	}

	for i, element := range suffixes {
		exponent := element.exponent
		nextExponent := exponent + 10
		if n < (1<<nextExponent) || i == len(suffixes)-1 {
			scaled := float64(n) / float64(uint(1)<<exponent)
			return fmt.Sprintf("%.2f %s", scaled, element.suffix)
		}
	}

	panic("Shouldn't reach here.")
}

func main() {
	flag.Parse()

	// Set up bare logging output.
	log.SetFlags(0)

	// Grab the bucket.
	bucket := getBucket()

	/////////////////////////////////////////////
	// Latency
	/////////////////////////////////////////////

	printHeading("Latency")

	avgLatency, err := measureLatency(bucket)
	if err != nil {
		log.Fatalf("measureLatency: %v", err)
	}

	log.Printf(
		"Average latency: %d ms",
		uint64(float64(avgLatency)/float64(time.Millisecond)),
	)

	/////////////////////////////////////////////
	// Downstream bandwidth
	/////////////////////////////////////////////

	printHeading("Downstream bandwidth")

	downstreamDataSizes := []uint{1 << 18, 1 << 20}
	downstreamParallelisms := []uint{1, 2}

	for _, dataSize := range downstreamDataSizes {
		for _, parallelism := range downstreamParallelisms {
			bytesPerSecond, err := measureDownstreamBandwidth(
				bucket,
				dataSize,
				parallelism)

			if err != nil {
				log.Fatalf("measureDownstreamBandwidth: %v", err)
			}

			log.Printf(
				"%s, parallelism %d: %s/s\n",
				formatBytes(uint64(dataSize)),
				parallelism,
				formatBytes(uint64(bytesPerSecond)),
			)
		}
	}

	/////////////////////////////////////////////
	// Upstream bandwidth
	/////////////////////////////////////////////

	printHeading("Upstream bandwidth")

	upstreamDataSizes := []uint{1 << 14, 1 << 18, 1 << 20}
	upstreamParallelisms := []uint{1, 2}

	for _, dataSize := range upstreamDataSizes {
		for _, parallelism := range upstreamParallelisms {
			bytesPerSecond, err := measureUpstreamBandwidth(
				bucket,
				dataSize,
				parallelism)

			if err != nil {
				log.Fatalf("measureUpstreamBandwidth: %v", err)
			}

			log.Printf(
				"%s, parallelism %d: %s/s\n",
				formatBytes(uint64(dataSize)),
				parallelism,
				formatBytes(uint64(bytesPerSecond)),
			)
		}
	}
}
