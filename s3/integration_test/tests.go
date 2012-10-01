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
	"github.com/jacobsa/aws/s3"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type BucketTest struct {
	bucket s3.Bucket
}

func init() { RegisterTestSuite(&BucketTest{}) }

func (t *BucketTest) ensureDeleted(key string) {
	err := t.bucket.DeleteObject(key)
	AssertEq(nil, err, "Couldn't delete object: %s", key)
}

func (t *BucketTest) SetUp(i *TestInfo) {
	var err error

	// Open a bucket.
	t.bucket, err = s3.OpenBucket(*bucketName, s3.Region(*region), accessKey)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

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
	_, err := t.bucket.GetObject("some_key")

	ExpectThat(err, Error(HasSubstr("404")))
	ExpectThat(err, Error(HasSubstr("some_key")))
	ExpectThat(err, Error(HasSubstr("exist")))
}

func (t *BucketTest) StoreThenGetEmptyObject() {
	key := "some_key"
	defer t.ensureDeleted(key)

	data := []byte{}

	// Store
	err := t.bucket.StoreObject(key, data)
	AssertEq(nil, err)

	// Get
	returnedData, err := t.bucket.GetObject(key)
	AssertEq(nil, err)
	ExpectThat(returnedData, DeepEquals(data))
}

func (t *BucketTest) StoreThenGetNonEmptyObject() {
	key := "some_key"
	defer t.ensureDeleted(key)

	data := []byte{0x17, 0x19, 0x00, 0x02, 0x03}

	// Store
	err := t.bucket.StoreObject(key, data)
	AssertEq(nil, err)

	// Get
	returnedData, err := t.bucket.GetObject(key)
	AssertEq(nil, err)
	ExpectThat(returnedData, DeepEquals(data))
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

func (t *BucketTest) KeysWithSpecialCharacters() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) DeleteNonExistentObject() {
	ExpectFalse(true, "TODO")
}

func (t *BucketTest) DeleteThenListAndGetObject() {
	ExpectFalse(true, "TODO")
}
