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

// Region represents a regional endpoint to S3. Resources created within one
// region are entirely independent of those created in others. You should use
// one of the region constants defined by this package when referring to
// regions.
//
// See here for more info:
//
//     http://goo.gl/FhIRw
//
type Region string

const (
	RegionUsStandard           Region = "s3.amazonaws.com"
	RegionUsWestOregon         Region = "s3-us-west-2.amazonaws.com"
	RegionUsWestNorCal         Region = "s3-us-west-1.amazonaws.com"
	RegionEuIreland            Region = "s3-eu-west-1.amazonaws.com"
	RegionApacSingapore        Region = "s3-ap-southeast-1.amazonaws.com"
	RegionApacSydney           Region = "s3-ap-southeast-2.amazonaws.com"
	RegionApacTokyo            Region = "s3-ap-northeast-1.amazonaws.com"
	RegionSouthAmericaSaoPaulo Region = "s3-sa-east-1.amazonaws.com"
)
