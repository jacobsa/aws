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

package conn

import (
	"github.com/jacobsa/aws"
)

type Signer interface {
	// Add appropriate parameters to the supplied request in order to sign it.
	SignRequest(req Request) error
}

// Create a signer that uses the supplied key for requests to the given host.
func NewSigner(key aws.AccessKey, host string) (Signer, error) {
	return newSigner(key, host, computeStringToSign), nil
}

// The underlying constructor, which accepts a "string to sign" function for
// testability.
func newSigner(
	key aws.AccessKey,
	host string,
	sts func (Request, string) (string, error)) Signer {
	return &signer{}
}

type signer struct {
}

func (s *signer) SignRequest(req Request) error {
	return nil
}
