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

package sdb

import (
	"github.com/jacobsa/aws/sdb/conn"
	. "github.com/jacobsa/ogletest"
)

////////////////////////////////////////////////////////////////////////
// Fake Conn
////////////////////////////////////////////////////////////////////////

type fakeConn struct {
	// Argument received
	req conn.Request

	// Response to return
	resp []byte
	err  error
}

func (c *fakeConn) SendRequest(r conn.Request) ([]byte, error) {
	if c.req != nil {
		panic("Already called!")
	}

	c.req = r
	return c.resp, c.err
}

////////////////////////////////////////////////////////////////////////
// Common test class
////////////////////////////////////////////////////////////////////////

// A common helper class.
type domainTest struct {
	name   string
	c      *fakeConn
	domain Domain
}

func (t *domainTest) SetUp(i *TestInfo) {
	var err error

	t.name = "some_domain"
	t.c = &fakeConn{}

	t.domain, err = newDomain(t.name, t.c)
	AssertEq(nil, err)
}
