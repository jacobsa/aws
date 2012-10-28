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

package s3util_test

import (
	"errors"
	"github.com/jacobsa/aws/s3/mock"
	"github.com/jacobsa/aws/s3/s3util"
	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestKeys(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type ListAllKeysTest struct {
	bucket mock_s3.MockBucket

	keys []string
	err  error
}

func init() { RegisterTestSuite(&ListAllKeysTest{}) }

func (t *ListAllKeysTest) SetUp(i *TestInfo) {
	t.bucket = mock_s3.NewMockBucket(i.MockController, "bucket")
}

func (t *ListAllKeysTest) call() {
	t.keys, t.err = s3util.ListAllKeys(t.bucket)
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *ListAllKeysTest) CallsListKeyRepeatedly() {
	// ListKeys (call 0)
	keys0 := []string{"burrito", "enchilada"}

	ExpectCall(t.bucket, "ListKeys")("").
		WillOnce(oglemock.Return(keys0, nil))

	// ListKeys (call 1)
	keys1 := []string{"queso", "taco"}

	ExpectCall(t.bucket, "ListKeys")("enchilada").
		WillOnce(oglemock.Return(keys1, nil))

	// ListKeys (call 2)
	ExpectCall(t.bucket, "ListKeys")("taco").
		WillOnce(oglemock.Return(nil, errors.New("")))

	// Call
	t.call()
}

func (t *ListAllKeysTest) ListKeyReturnsError() {
	// ListKeys
	ExpectCall(t.bucket, "ListKeys")(Any()).
		WillOnce(oglemock.Return([]string{"a"}, nil)).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	t.call()

	ExpectThat(t.err, Error(HasSubstr("ListKeys")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *ListAllKeysTest) ListKeyReturnsNoKeys() {
	// ListKeys
	ExpectCall(t.bucket, "ListKeys")(Any()).
		WillOnce(oglemock.Return([]string{}, nil))

	// Call
	t.call()
	AssertEq(nil, t.err)

	ExpectThat(t.keys, ElementsAre())
}

func (t *ListAllKeysTest) ListKeyReturnsSomeKeys() {
	// ListKeys
	ExpectCall(t.bucket, "ListKeys")(Any()).
		WillOnce(oglemock.Return([]string{"burrito", "enchilada"}, nil)).
		WillOnce(oglemock.Return([]string{"taco"}, nil)).
		WillOnce(oglemock.Return([]string{}, nil))

	// Call
	t.call()
	AssertEq(nil, t.err)

	ExpectThat(
		t.keys,
		ElementsAre(
			"burrito",
			"enchilada",
			"taco",
		),
	)
}
