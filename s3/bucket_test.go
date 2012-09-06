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
	"github.com/jacobsa/aws/s3/auth/mock"
	"github.com/jacobsa/aws/s3/http/mock"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

func TestBucket(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type bucketTest struct {
	httpConn mock_http.MockConn
	signer mock_auth.MockSigner
	bucket Bucket
}

func (t *bucketTest) SetUp(i *TestInfo) {
	var err error

	t.httpConn = mock_http.NewMockConn(i.MockController, "httpConn")
	t.signer = mock_auth.NewMockSigner(i.MockController, "signer")
	t.bucket, err = openBucket("some.bucket", t.httpConn, t.signer)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// GetObject
////////////////////////////////////////////////////////////////////////

type GetObjectTest struct {
	bucketTest
}

func init() { RegisterTestSuite(&GetObjectTest{}) }

func (t *GetObjectTest) DoesFoo() {
	ExpectEq("TODO", "")
}

////////////////////////////////////////////////////////////////////////
// StoreObject
////////////////////////////////////////////////////////////////////////

type StoreObjectTest struct {
	bucketTest
}

func init() { RegisterTestSuite(&StoreObjectTest{}) }

func (t *StoreObjectTest) KeyNotValidUtf8() {
	key := "\x80\x81\x82"
	data := []byte{}

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("valid")))
	ExpectThat(err, Error(HasSubstr("UTF-8")))
}

func (t *StoreObjectTest) KeyTooLong() {
	ExpectEq("TODO", "")
}

func (t *StoreObjectTest) CallsSignerAndConn() {
	ExpectEq("TODO", "")
}

func (t *StoreObjectTest) ConnReturnsError() {
	ExpectEq("TODO", "")
}

func (t *StoreObjectTest) ServerReturnsError() {
	ExpectEq("TODO", "")
}

func (t *StoreObjectTest) ServerSaysOkay() {
	ExpectEq("TODO", "")
}
