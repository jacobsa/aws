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
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/exp/sdb/conn"
	"github.com/jacobsa/aws/time"
	"net/url"
	"regexp"
)

// A regexp for valid domain names.
//
// Reference:
//     http://goo.gl/Kkbnf
var domainNameRe = regexp.MustCompile(`^[-a-zA-Z0-9_.]{3,255}$`)

// The API version for which this package was written.
const apiVersion = "2009-04-15"

// The name of an item within a SimpleDB domain. Item names must be UTF-8
// strings no longer than 1024 bytes. They must contain only characters that
// are valid in XML 1.0 documents, as defined by Section 2.2 of the XML 1.0
// spec. (Note that this is a more restrictive condition than imposed by
// SimpleDB itself, and is done for the sake of Go's XML 1.0 parser.)
//
// For more info:
//
//     http://goo.gl/Fkjnz
//     http://goo.gl/csem8
//
type ItemName string

// An attribute is a (name, value) pair possessed by an item. Items contain
// sets of attributes; they may contain multiple attributes with the same name,
// but not with the same (name, value) pair.
//
// Attribute names and values share the same restrictions as those on item
// names.
type Attribute struct {
	Name  string
	Value string
}

// A single result from a Select operation.
type SelectedItem struct {
	Name ItemName
	Attributes []Attribute
}

// SimpleDB represents an authenticated connection to a particular endpoint the
// SimpleDB service.
type SimpleDB interface {
	// Ensure that the domain with the specified name exists, then return it.
	OpenDomain(name string) (Domain, error)

	// Remove the supplied domain from the DB. Any domains previously opened with
	// this domain's name will no longer be valid. It is not an error to call
	// this method for a domain that has already been deleted.
	DeleteDomain(d Domain) error

	// Retrieve a set of items and their attributes based on a query string.
	//
	// constistentRead specifies whether completely fresh data is needed or not.
	//
	// If the selected result set is too large, a "next token" will be returned.
	// It can be passed to the Select method to resume where the previous result
	// set left off. For the initial query, use nil.
	//
	// For more info:
	//
	//     http://goo.gl/GTsSZ
	//
	Select(
		query string,
		constistentRead bool,
		nextToken []byte) (results []SelectedItem, tok []byte, err error)
}

// Return a SimpleDB connection tied to the given region, using the sipplied
// access key to authenticate requests.
func NewSimpleDB(region Region, key aws.AccessKey) (db SimpleDB, err error) {
	// Open an appropriate HTTP connection.
	endpoint := &url.URL{
		Scheme: "https",
		Host:   string(region),
	}

	httpConn, err := conn.NewHttpConn(endpoint)
	if err != nil {
		err = fmt.Errorf("Opening HTTP connection: %v", err)
		return
	}

	// Create a request signer.
	signer, err := conn.NewSigner(key, endpoint.Host)
	if err != nil {
		err = fmt.Errorf("Creating signer: %v", err)
		return
	}

	// Create a connection to the server.
	c, err := conn.NewConn(key, httpConn, signer, time.RealClock())
	if err != nil {
		err = fmt.Errorf("Creating connection: %v", err)
		return
	}

	return newSimpleDB(c)
}

// Create a SimpleDB with the supplied underlying connection.
func newSimpleDB(c conn.Conn) (SimpleDB, error) {
	return &simpleDB{c}, nil
}

type simpleDB struct {
	c conn.Conn
}

func (db *simpleDB) OpenDomain(name string) (d Domain, err error) {
	// Make sure the requested name is legal.
	if !domainNameRe.MatchString(name) {
		err = fmt.Errorf("Illegal domain name: %s", name)
		return
	}

	// Build a request.
	req := conn.Request{
		"Action":     "CreateDomain",
		"DomainName": name,
		"Version":    apiVersion,
	}

	// Call the connection.
	if _, err = db.c.SendRequest(req); err != nil {
		err = fmt.Errorf("SendRequest: %v", err)
		return
	}

	// Create the object.
	return newDomain(name, db.c)
}

func (db *simpleDB) DeleteDomain(d Domain) (err error) {
	// Build a request.
	req := conn.Request{
		"Action":     "DeleteDomain",
		"DomainName": d.Name(),
		"Version":    apiVersion,
	}

	// Call the connection.
	if _, err = db.c.SendRequest(req); err != nil {
		err = fmt.Errorf("SendRequest: %v", err)
		return
	}

	return
}
