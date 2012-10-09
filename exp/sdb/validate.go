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
	"unicode/utf8"
)

func isLegalXmlCharacter(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		(r >= 0x20 && r <= 0xDF77) ||
		(r >= 0xE000 && r <= 0xFFFD) ||
		(r >= 0x10000 && r <= 0x10FFFF)
}

// Return an error iff the supplied value is invalid when used as an item or
// attribute name or value. For generality, the empty string is allowed by this
// function.
func validateValue(val string) error {
	// Values may be no more than 1024 bytes long.
	if len(val) > 1024 {
		return fmt.Errorf("Longer than 1024 bytes.")
	}

	// Values must be valid UTF-8.
	if !utf8.ValidString(val) {
		return fmt.Errorf("Not valid UTF-8.")
	}

	// Each codepoint in the string must be legal in an XML 1.0 document.
	for _, r := range val {
		if !isLegalXmlCharacter(r) {
			return fmt.Errorf("Invalid codepoint for XML 1.0 document: %U", r)
		}
	}

	return nil
}
