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
	t.bucket, err = s3.OpenBucket(*g_bucketName, s3.Region(*g_region), g_accessKey)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *BucketTest) WrongAccessKeySecret() {
	// Open a bucket with the wrong key.
	wrongKey := g_accessKey
	wrongKey.Secret += "taco"

	bucket, err := s3.OpenBucket(*g_bucketName, s3.Region(*g_region), wrongKey)
	AssertEq(nil, err)

	// Attempt to do something.
	_, err = bucket.ListKeys("")
	ExpectThat(err, Error(HasSubstr("signature")))
}

func (t *BucketTest) InvalidUtf8Keys() {
	ExpectEq("TODO", "")
}

func (t *BucketTest) LongKeys() {
	ExpectEq("TODO", "")
}

func (t *BucketTest) NullBytesInKeys() {
	ExpectEq("TODO", "")
}

func (t *BucketTest) NonGraphicalCharacterInKeys() {
	ExpectEq("TODO", "")
}

func (t *BucketTest) EmptyKeys() {
	ExpectEq("TODO", "")
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
	var keys []string
	var err error

	// Create several keys. S3 returns keys in an XML 1.0 document, and according
	// to Section 2.2 of the spec the smallest legal character is #x9, so a
	// string's successor in that space of strings is computed by appending \x09.
	//
	// S3 will actually allow you to create a smaller key, e.g. "bar\x01", but
	// Go's xml package will then refuse to parse its LIST responses.
	toCreate := []string{
		"foo",
		"bar",
		"bar\x09",
		"bar\x09\x09",
		"baz",
	}

	for _, key := range toCreate {
		defer t.ensureDeleted(key)
		err := t.bucket.StoreObject(key, []byte{})
		AssertEq(nil, err, "Creating object: %s", key)
	}

	// From start.
	keys, err = t.bucket.ListKeys("")
	AssertEq(nil, err)
	ExpectThat(
		keys,
		ElementsAre(
		"bar",
		"bar\x09",
		"bar\x09\x09",
		"baz",
		"foo",
	))

	// Just before bar\x09.
	keys, err = t.bucket.ListKeys("bar")
	AssertEq(nil, err)
	ExpectThat(
		keys,
		ElementsAre(
		"bar\x09",
		"bar\x09\x09",
		"baz",
		"foo",
	))

	// At bar\x09.
	keys, err = t.bucket.ListKeys("bar\x09")
	AssertEq(nil, err)
	ExpectThat(
		keys,
		ElementsAre(
		"bar\x09\x09",
		"baz",
		"foo",
	))

	// Just after bar\x09.
	keys, err = t.bucket.ListKeys("bar\x09\x09")
	AssertEq(nil, err)
	ExpectThat(
		keys,
		ElementsAre(
		"baz",
		"foo",
	))

	// At last key.
	keys, err = t.bucket.ListKeys("foo")
	AssertEq(nil, err)
	ExpectThat(keys, ElementsAre())

	// Just after last key.
	keys, err = t.bucket.ListKeys("foo\x09")
	AssertEq(nil, err)
	ExpectThat(keys, ElementsAre())

	// Well after last key.
	keys, err = t.bucket.ListKeys("qux")
	AssertEq(nil, err)
	ExpectThat(keys, ElementsAre())
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
