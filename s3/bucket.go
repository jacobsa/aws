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

import (
	"errors"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3/auth"
	"github.com/jacobsa/aws/s3/http"
	"github.com/jacobsa/aws/s3/time"
	"net/url"
)

// NonExistentBucketError represents an error due to an attempt to work with a
// bucket that doesn't exist according to S3.
type NonExistentBucketError struct {
	s string
}

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

// OpenBucket returns a Bucket tied to a given name in whe given region. You
// must have previously created the bucket in the region, and the supplied
// access key must have access to it.
//
// If the supplied bucket doesn't exist, a *NonExistentBucketError is returned.
//
// To easily create a bucket, use the AWS Console:
//
//     http://aws.amazon.com/console/
//
func OpenBucket(name string, region Region, key aws.AccessKey) (Bucket, error) {
	// Create a connection to the given region's endpoint.
	endpoint := &url.URL{Scheme: "https", Host: string(region)}
	httpConn, err := http.NewConn(endpoint)
	if err != nil {
		return nil, fmt.Errorf("http.NewConn: %v", err)
	}

	// Create an appropriate request signer.
	signer, err := auth.NewSigner(&key)
	if err != nil {
		return nil, fmt.Errorf("auth.NewSigner: %v", err)
	}

	return openBucket(name, httpConn, signer, time.RealClock())
}

func (e *NonExistentBucketError) Error() string {
	return e.s
}

// A version of OpenBucket with the ability to inject dependencies, for
// testability.
func openBucket(
	name string,
	httpConn http.Conn,
	signer auth.Signer,
	clock time.Clock) (Bucket, error) {
	return nil, errors.New("TODO: Implement openBucket.")
}
