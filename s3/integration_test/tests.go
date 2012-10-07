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
	"fmt"
	"github.com/jacobsa/aws/s3"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"sync"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

// Run the supplied function for every integer in [0, n) with some degree of
// parallelism, returning an error if any invocation returns an error.
func runForRange(n int, f func (int) error) (err error) {
		// Set up channels. The work channel should be buffered so that we don't
		// have to block writing to it before checking for errors below. The error
		// channel must be buffered so that no worker goroutine gets stuck writing
		// a result to it and never returns. The stop channel must not be buffered
		// so that we can be sure that no more work will be processed when we
		// return below.
		work := make(chan int, n)
		errs := make(chan error, n)
		stop := make(chan bool)

		// Launch worker functions that attempt to do work, returning if a read
		// from the stop channel succeeds.
		processWork := func() {
			for {
				select {
				case i := <-work:
					errs<- f(i)
				case <-stop:
					return
				}
			}
		}

		const numWorkers = 128
		for i := 0; i < numWorkers; i++ {
			go processWork()
		}

		// Feed the workers work.
		for i := 0; i < n; i++ {
			work<- i
		}

		// Read results, stopping immediately if there is an error.
		for i := 0; i < n; i++ {
			err = <-errs
			if err != nil {
				break
			}
		}

		// Stop all of the workers, and wait for them to stop. This ensures that
		// no piece of work is in progress when we return.
		for i := 0; i < numWorkers; i++ {
			stop<- true
		}

		return
	}

type BucketTest struct {
	bucket s3.Bucket

	mutex sync.Mutex
	keysToDelete []string
}

func init() { RegisterTestSuite(&BucketTest{}) }

// Ensure that the given key is deleted before the test finishes.
func (t *BucketTest) ensureDeleted(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.keysToDelete = append(t.keysToDelete, key)
}

func (t *BucketTest) SetUp(i *TestInfo) {
	var err error

	// Open a bucket.
	t.bucket, err = s3.OpenBucket(*g_bucketName, s3.Region(*g_region), g_accessKey)
	AssertEq(nil, err)
}

func (t *BucketTest) TearDown() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	err := runForRange(len(t.keysToDelete), func (i int) error {
		key := t.keysToDelete[i]
		if err := t.bucket.DeleteObject(key); err != nil {
			return fmt.Errorf("Couldn't delete key %s: %v", key, err)
		}

		return nil
	})

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
	t.ensureDeleted(key)

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
	t.ensureDeleted(key)

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

	err = runForRange(len(toCreate), func(i int) error {
		key := toCreate[i]
		t.ensureDeleted(key)
		return t.bucket.StoreObject(key, []byte{})
	})

	AssertEq(nil, err)

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
	var err error

	// Decide on many keys.
	const numKeys = 1001
	allKeys := make([]string, numKeys)

	for i, _ := range allKeys {
		allKeys[i] = fmt.Sprintf("%08x", i)
	}

	// Create them.
	err = runForRange(numKeys, func(i int) error {
		key := allKeys[i]
		t.ensureDeleted(key)
		return t.bucket.StoreObject(key, []byte{})
	})

	AssertEq(nil, err)

	// List them progressively.
	lb := ""
	keysListed := []string{}

	for {
		keys, err := t.bucket.ListKeys(lb)

		AssertEq(nil, err)
		AssertLt(len(keys), numKeys)

		if len(keys) == 0 {
			break
		}

		keysListed = append(keysListed, keys...)
		lb = keys[len(keys) - 1]
	}

	// We should have gotten them all back.
	ExpectThat(keysListed, DeepEquals(allKeys))
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
