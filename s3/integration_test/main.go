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
//
// An integration test that uses a real S3 account. Run as follows:
//
//     go run integration_test/*.go \
//         -key_id <key ID> \
//         -bucket <bucket> \
//         -region s3-ap-northeast-1.amazonaws.com
//
// Before doing this, create an empty bucket (or delete the contents of an
// existing bucket) using the S3 management console.

package main

import (
	"flag"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"os"
	"regexp"
	"testing"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

var keyId = flag.String("key_id", "", "Access key ID.")
var bucketName = flag.String("bucket", "", "Bucket name.")
var region = flag.String("region", "", "Region endpoint server.")
var accessKey aws.AccessKey

type integrationTest struct {
}

////////////////////////////////////////////////////////////////////////
// Bucket
////////////////////////////////////////////////////////////////////////

type BucketTest struct {
	integrationTest
	bucket s3.Bucket
}

func init() { RegisterTestSuite(&BucketTest{}) }

func (t *BucketTest) SetUp(i *TestInfo) {
	var err error

	// Open a bucket.
	t.bucket, err = s3.OpenBucket(*bucketName, s3.Region(*region), accessKey)
	AssertEq(nil, err)
}

func (t *BucketTest) WrongAccessKeySecret() {
	// Open a bucket with the wrong key.
	wrongKey := accessKey
	wrongKey.Secret += "taco"

	bucket, err := s3.OpenBucket(*bucketName, s3.Region(*region), wrongKey)
	AssertEq(nil, err)

	// Attempt to do something.
	_, err = bucket.ListKeys("")
	ExpectThat(err, Error(HasSubstr("signature")))
}

func (t *BucketTest) GetNonExistentObject() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) StoreThenGetObject() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) StoreThenDeleteObject() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListEmptyBucket() {
	var keys []string
	var err error

	// From start.
	keys, err = t.bucket.ListKeys("")
	AssertEq(nil, err)
	ExpectThat(keys, ElementsAre())

	// From middle.
	keys, err = t.bucket.ListKeys("foo")
	AssertEq(nil, err)
	ExpectThat(keys, ElementsAre())
}

func (t *BucketTest) ListWithEmptyMinimum() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListWithInvalidUtf8Minimum() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListWithLongMinimum() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListWithNullByteInMinimum() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListFewKeys() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) ListManyKeys() {
	ExpectFalse(true, "TODO")
}

////////////////////////////////////////////////////////////////////////
// main
////////////////////////////////////////////////////////////////////////

func main() {
	flag.Parse()

	if *keyId == "" {
		fmt.Println("You must set the -key_id flag.")
		fmt.Println("Find a key ID here:")
		fmt.Println("    https://portal.aws.amazon.com/gp/aws/securityCredentials")
		os.Exit(1)
	}

	if *bucketName == "" {
		fmt.Println("You must set the -bucket flag.")
		fmt.Println("Manage your buckets here:")
		fmt.Println("    https://console.aws.amazon.com/s3/")
		os.Exit(1)
	}

	if *region == "" {
		fmt.Println("You must set the -region flag. See region.go.")
		os.Exit(1)
	}

	// Read in the access key.
	accessKey.Id = *keyId
	accessKey.Secret = readPassword("Access key secret: ")

	// Run the tests.
	matchString := func(pat, str string) (bool, error) {
		re, err := regexp.Compile(pat)
		if err != nil {
			return false, err
		}

		return re.MatchString(str), nil
	}

	testing.Main(
		matchString,
		[]testing.InternalTest{
			testing.InternalTest{
				Name: "IntegrationTest",
				F:    func(t *testing.T) { RunTests(t) },
			},
		},
		[]testing.InternalBenchmark{},
		[]testing.InternalExample{},
	)
}
