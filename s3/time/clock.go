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

// Package time is an implementation detail. You should not use it directly.
package time

import (
	"time"
)

type Clock interface {
	// Return the current time.
	Now() time.Time
}

// Return a clock that uses the real time, with locations set to UTC.
func RealClock() Clock {
	return &realClock{}
}

type realClock struct {}

func (c *realClock) Now() time.Time {
	return time.Now().UTC()
}
