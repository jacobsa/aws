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
// An integration test that uses a real SimpleDB account. Run as follows:
//
//     go run integration_test/*.go \
//         -key_id <key ID> \
//         -domain <domain> \
//         -region sdb.ap-northeast-1.amazonaws.com
//
// Before doing this, create an empty domain (or delete the contents of an
// existing domain).

package main

import (
	"flag"
	"fmt"
	"github.com/jacobsa/aws"
	"github.com/jacobsa/ogletest"
	"os"
	"regexp"
	"testing"
)

////////////////////////////////////////////////////////////////////////
// Globals
////////////////////////////////////////////////////////////////////////

var g_keyId = flag.String("key_id", "", "Access key ID.")
var g_region = flag.String("region", "", "Region endpoint server.")
var g_accessKey aws.AccessKey

////////////////////////////////////////////////////////////////////////
// main
////////////////////////////////////////////////////////////////////////

func main() {
	flag.Parse()

	if *g_keyId == "" {
		fmt.Println("You must set the -key_id flag.")
		fmt.Println("Find a key ID here:")
		fmt.Println("    https://portal.aws.amazon.com/gp/aws/securityCredentials")
		os.Exit(1)
	}

	if *g_region == "" {
		fmt.Println("You must set the -region flag. See region.go.")
		os.Exit(1)
	}

	// Read in the access key.
	g_accessKey.Id = *g_keyId
	g_accessKey.Secret = readPassword("Access key secret: ")

	// Run the tests.
	matchString := func(pat, str string) (bool, error) {
		re, err := regexp.Compile(pat)
		if err != nil {
			return false, err
		}

		return re.MatchString(str), nil
	}

	testing.Main(
		matchString,
		[]testing.InternalTest{
			testing.InternalTest{
				Name: "IntegrationTest",
				F:    func(t *testing.T) { ogletest.RunTests(t) },
			},
		},
		[]testing.InternalBenchmark{},
		[]testing.InternalExample{},
	)
}
