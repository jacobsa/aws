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

package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/aws/s3/http"
)

// A Signer knows how to create signatures suitable for inclusion in an HTTP
// request to S3.
//
// Reference:
//     http://goo.gl/bOVNo
type Signer interface {
	// Add an appropriate signature header to the supplied HTTP request.
	Sign(r *http.Request) error
}

// NewSigner creates a Signer using the supplied access key.
func NewSigner(key *aws.AccessKey) (Signer, error) {
	return newSigner(stringToSign, key)
}

// newSigner is a helper used by NewSigner, split out for testability. It
// allows you to inject the function that is used to determine the string to
// sign for any given request.
func newSigner(
	sts func(*http.Request) (string, error),
	key *aws.AccessKey) (Signer, error) {
	return &signer{sts, key}, nil
}

type signer struct {
	sts func(*http.Request) (string, error)
	key *aws.AccessKey
}

func (s *signer) Sign(r *http.Request) error {
	// Canonicalize the request.
	toSign, err := s.sts(r)
	if err != nil {
		return fmt.Errorf("stringToSign: %v", err)
	}

	// Sign the request.
	h := hmac.New(sha1.New, []byte(s.key.Secret))
	if _, err := h.Write([]byte(toSign)); err != nil {
		return fmt.Errorf("hmac.Write: %v", err)
	}

	signature := h.Sum(nil)

	// Base64-encode the result.
	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := encoder.Write(signature); err != nil {
		return fmt.Errorf("encoder.Write: %v", err)
	}

	encoder.Close()

	// Add the appropriate header.
	r.Headers["Authorization"] = fmt.Sprintf("AWS %s:%s", s.key.Id, buf.String())

	return nil
}
