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
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/jacobsa/aws/s3/auth/mock"
	"github.com/jacobsa/aws/s3/http"
	"github.com/jacobsa/aws/s3/http/mock"
	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/oglemock"
	. "github.com/jacobsa/ogletest"
	"strings"
	"testing"
	"time"
)

func TestBucket(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func computeBase64Md5(d []byte) string {
	h := md5.New()
	if _, err := h.Write(d); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	e := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := e.Write(h.Sum(nil)); err != nil {
		panic(err)
	}

	e.Close()
	return buf.String()
}

type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	return c.now
}

type bucketTest struct {
	httpConn mock_http.MockConn
	signer   mock_auth.MockSigner
	bucket   Bucket
	clock    *fakeClock
}

func (t *bucketTest) SetUp(i *TestInfo) {
	var err error

	t.httpConn = mock_http.NewMockConn(i.MockController, "httpConn")
	t.signer = mock_auth.NewMockSigner(i.MockController, "signer")
	t.clock = &fakeClock{}

	t.bucket, err = openBucket("some.bucket", t.httpConn, t.signer, t.clock)
	AssertEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// GetObject
////////////////////////////////////////////////////////////////////////

type GetObjectTest struct {
	bucketTest
}

func init() { RegisterTestSuite(&GetObjectTest{}) }

func (t *GetObjectTest) KeyNotValidUtf8() {
	key := "\x80\x81\x82"

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("valid")))
	ExpectThat(err, Error(HasSubstr("UTF-8")))
}

func (t *GetObjectTest) KeyTooLong() {
	key := strings.Repeat("a", 1025)

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("1024")))
	ExpectThat(err, Error(HasSubstr("bytes")))
}

func (t *GetObjectTest) KeyContainsNullByte() {
	key := "taco\x00burrito"

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("null")))
}

func (t *GetObjectTest) KeyIsEmpty() {
	key := ""

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("empty")))
}

func (t *GetObjectTest) CallsSigner() {
	key := "foo/bar/baz"

	// Clock
	t.clock.now = time.Date(1985, time.March, 18, 15, 33, 17, 123, time.UTC)

	// Signer
	var httpReq *http.Request
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		httpReq = r
		return errors.New("")
	}))

	// Call
	t.bucket.GetObject(key)

	AssertNe(nil, httpReq)
	ExpectEq("GET", httpReq.Verb)
	ExpectEq("/some.bucket/foo/bar/baz", httpReq.Path)
	ExpectEq("Mon, 18 Mar 1985 15:33:17 UTC", httpReq.Headers["Date"])
}

func (t *GetObjectTest) SignerReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(errors.New("taco")))

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("Sign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *GetObjectTest) CallsConn() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		r.Verb = "burrito"
		return nil
	}))

	// Conn
	var httpReq *http.Request
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) (*http.Response, error) {
		httpReq = r
		return nil, errors.New("")
	}))

	// Call
	t.bucket.GetObject(key)

	AssertNe(nil, httpReq)
	ExpectEq("burrito", httpReq.Verb)
}

func (t *GetObjectTest) ConnReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("SendRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *GetObjectTest) ServerReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 500,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	_, err := t.bucket.GetObject(key)

	ExpectThat(err, Error(HasSubstr("server")))
	ExpectThat(err, Error(HasSubstr("500")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *GetObjectTest) ReturnsResponseBody() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	data, err := t.bucket.GetObject(key)
	AssertEq(nil, err)

	ExpectThat(data, DeepEquals([]byte("taco")))
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
	key := strings.Repeat("a", 1025)
	data := []byte{}

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("1024")))
	ExpectThat(err, Error(HasSubstr("bytes")))
}

func (t *StoreObjectTest) KeyContainsNullByte() {
	key := "taco\x00burrito"
	data := []byte{}

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("null")))
}

func (t *StoreObjectTest) KeyIsEmpty() {
	key := ""
	data := []byte{}

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("empty")))
}

func (t *StoreObjectTest) CallsSigner() {
	key := "foo/bar/baz"
	data := []byte{0x00, 0xde, 0xad, 0xbe, 0xef}

	// Clock
	t.clock.now = time.Date(1985, time.March, 18, 15, 33, 17, 123, time.UTC)

	// Signer
	var httpReq *http.Request
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		httpReq = r
		return errors.New("")
	}))

	// Call
	t.bucket.StoreObject(key, data)

	AssertNe(nil, httpReq)
	ExpectEq("PUT", httpReq.Verb)
	ExpectEq("/some.bucket/foo/bar/baz", httpReq.Path)
	ExpectEq("Mon, 18 Mar 1985 15:33:17 UTC", httpReq.Headers["Date"])
	ExpectEq(computeBase64Md5(data), httpReq.Headers["Content-MD5"])
	ExpectThat(httpReq.Body, DeepEquals(data))
}

func (t *StoreObjectTest) SignerReturnsError() {
	key := "a"
	data := []byte{}

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(errors.New("taco")))

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("Sign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *StoreObjectTest) CallsConn() {
	key := "a"
	data := []byte{}

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		r.Verb = "burrito"
		return nil
	}))

	// Conn
	var httpReq *http.Request
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) (*http.Response, error) {
		httpReq = r
		return nil, errors.New("")
	}))

	// Call
	t.bucket.StoreObject(key, data)

	AssertNe(nil, httpReq)
	ExpectEq("burrito", httpReq.Verb)
}

func (t *StoreObjectTest) ConnReturnsError() {
	key := "a"
	data := []byte{}

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("SendRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *StoreObjectTest) ServerReturnsError() {
	key := "a"
	data := []byte{}

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 500,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectThat(err, Error(HasSubstr("server")))
	ExpectThat(err, Error(HasSubstr("500")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *StoreObjectTest) ServerSaysOkay() {
	key := "a"
	data := []byte{}

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	err := t.bucket.StoreObject(key, data)

	ExpectEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// DeleteObject
////////////////////////////////////////////////////////////////////////

type DeleteObjectTest struct {
	bucketTest
}

func init() { RegisterTestSuite(&DeleteObjectTest{}) }

func (t *DeleteObjectTest) KeyNotValidUtf8() {
	key := "\x80\x81\x82"

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("valid")))
	ExpectThat(err, Error(HasSubstr("UTF-8")))
}

func (t *DeleteObjectTest) KeyTooLong() {
	key := strings.Repeat("a", 1025)

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("1024")))
	ExpectThat(err, Error(HasSubstr("bytes")))
}

func (t *DeleteObjectTest) KeyContainsNullByte() {
	key := "taco\x00burrito"

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("null")))
}

func (t *DeleteObjectTest) KeyIsEmpty() {
	key := ""

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("empty")))
}

func (t *DeleteObjectTest) CallsSigner() {
	key := "taco burrito"

	// Clock
	t.clock.now = time.Date(1985, time.March, 18, 15, 33, 17, 123, time.UTC)

	// Signer
	var httpReq *http.Request
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		httpReq = r
		return errors.New("")
	}))

	// Call
	t.bucket.DeleteObject(key)

	AssertNe(nil, httpReq)
	ExpectEq("DELETE", httpReq.Verb)
	ExpectEq("/some.bucket/taco burrito", httpReq.Path)
	ExpectEq("Mon, 18 Mar 1985 15:33:17 UTC", httpReq.Headers["Date"])
}

func (t *DeleteObjectTest) SignerReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(errors.New("taco")))

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("Sign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *DeleteObjectTest) CallsConn() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		r.Verb = "burrito"
		return nil
	}))

	// Conn
	var httpReq *http.Request
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) (*http.Response, error) {
		httpReq = r
		return nil, errors.New("")
	}))

	// Call
	t.bucket.DeleteObject(key)

	AssertNe(nil, httpReq)
	ExpectEq("burrito", httpReq.Verb)
}

func (t *DeleteObjectTest) ConnReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("SendRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *DeleteObjectTest) ServerReturnsError() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 500,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	err := t.bucket.DeleteObject(key)

	ExpectThat(err, Error(HasSubstr("server")))
	ExpectThat(err, Error(HasSubstr("500")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *DeleteObjectTest) ServerReturnsNoContent() {
	key := "a"

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 204,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	err := t.bucket.DeleteObject(key)
	ExpectEq(nil, err)
}

////////////////////////////////////////////////////////////////////////
// ListKeys
////////////////////////////////////////////////////////////////////////

type ListKeysTest struct {
	bucketTest
}

func init() { RegisterTestSuite(&ListKeysTest{}) }

func (t *ListKeysTest) CallsSignerWithEmptyMin() {
	prevKey := ""

	// Clock
	t.clock.now = time.Date(1985, time.March, 18, 15, 33, 17, 123, time.UTC)

	// Signer
	var httpReq *http.Request
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		httpReq = r
		return errors.New("")
	}))

	// Call
	t.bucket.ListKeys(prevKey)

	AssertNe(nil, httpReq)
	ExpectEq("GET", httpReq.Verb)
	ExpectEq("/some.bucket", httpReq.Path)
	ExpectEq("Mon, 18 Mar 1985 15:33:17 UTC", httpReq.Headers["Date"])

	marker, containsMarker := httpReq.Parameters["marker"]
	ExpectFalse(containsMarker, "marker: \"%s\"", marker)
}

func (t *ListKeysTest) CallsSignerWithNonEmptyMin() {
	prevKey := "taco burrito"

	// Clock
	t.clock.now = time.Date(1985, time.March, 18, 15, 33, 17, 123, time.UTC)

	// Signer
	var httpReq *http.Request
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		httpReq = r
		return errors.New("")
	}))

	// Call
	t.bucket.ListKeys(prevKey)

	AssertNe(nil, httpReq)
	ExpectEq("GET", httpReq.Verb)
	ExpectEq("/some.bucket", httpReq.Path)
	ExpectEq("taco burrito", httpReq.Parameters["marker"])
	ExpectEq("Mon, 18 Mar 1985 15:33:17 UTC", httpReq.Headers["Date"])
}

func (t *ListKeysTest) SignerReturnsError() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(errors.New("taco")))

	// Call
	_, err := t.bucket.ListKeys(prevKey)

	ExpectThat(err, Error(HasSubstr("Sign")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ListKeysTest) CallsConn() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) error {
		r.Verb = "burrito"
		return nil
	}))

	// Conn
	var httpReq *http.Request
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Invoke(func(r *http.Request) (*http.Response, error) {
		httpReq = r
		return nil, errors.New("")
	}))

	// Call
	t.bucket.ListKeys(prevKey)

	AssertNe(nil, httpReq)
	ExpectEq("burrito", httpReq.Verb)
}

func (t *ListKeysTest) ConnReturnsError() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(nil, errors.New("taco")))

	// Call
	_, err := t.bucket.ListKeys(prevKey)

	ExpectThat(err, Error(HasSubstr("SendRequest")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ListKeysTest) ServerReturnsError() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 500,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	_, err := t.bucket.ListKeys(prevKey)

	ExpectThat(err, Error(HasSubstr("server")))
	ExpectThat(err, Error(HasSubstr("500")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ListKeysTest) ResponseBodyIsJunk() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body:       []byte("taco"),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	_, err := t.bucket.ListKeys(prevKey)

	ExpectThat(err, Error(HasSubstr("Invalid")))
	ExpectThat(err, Error(HasSubstr("taco")))
}

func (t *ListKeysTest) WrongRootTag() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body: []byte(`
			<?xml version="1.0" encoding="UTF-8"?>
			<FooBar xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
				<Contents>
					<Key>some_key</Key>
				</Contents>
			</FooBar>`),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	_, err := t.bucket.ListKeys(prevKey)

	ExpectThat(err, Error(HasSubstr("Invalid")))
	ExpectThat(err, Error(HasSubstr("FooBar")))
}

func (t *ListKeysTest) ResponseContainsNoKeys() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body: []byte(`
			<?xml version="1.0" encoding="UTF-8"?>
			<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
			</ListBucketResult>`),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	keys, err := t.bucket.ListKeys(prevKey)
	AssertEq(nil, err)

	ExpectThat(keys, ElementsAre())
}

func (t *ListKeysTest) ResponseContainsSomeKeys() {
	prevKey := ""

	// Signer
	ExpectCall(t.signer, "Sign")(Any()).
		WillOnce(oglemock.Return(nil))

	// Conn
	resp := &http.Response{
		StatusCode: 200,
		Body: []byte(`
			<?xml version="1.0" encoding="UTF-8"?>
			<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
				<Contents>
					<Key>bar</Key>
				</Contents>
				<Contents>
					<Key>baz</Key>
				</Contents>
				<Contents>
					<Key>foo</Key>
				</Contents>
			</ListBucketResult>`),
	}

	ExpectCall(t.httpConn, "SendRequest")(Any()).
		WillOnce(oglemock.Return(resp, nil))

	// Call
	keys, err := t.bucket.ListKeys(prevKey)
	AssertEq(nil, err)

	ExpectThat(keys, ElementsAre("bar", "baz", "foo"))
}
