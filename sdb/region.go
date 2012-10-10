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

// Region represents a regional endpoint to SimpleDB. Domains created within
// one region are entirely independent of those created in others. You should
// use one of the region constants defined by this package when referring to
// regions.
//
// See here for more info:
//
//     http://goo.gl/BkF9n
//
type Region string

const (
	RegionUsEastNorthernVirginia Region = "sdb.amazonaws.com"
	RegionUsWestOregon           Region = "sdb.us-west-2.amazonaws.com"
	RegionUsWestNorCal           Region = "sdb.us-west-1.amazonaws.com"
	RegionEuIreland              Region = "sdb.eu-west-1.amazonaws.com"
	RegionApacSingapore          Region = "sdb.ap-southeast-1.amazonaws.com"
	RegionApacTokyo              Region = "sdb.ap-northeast-1.amazonaws.com"
	RegionSouthAmericaSaoPaulo   Region = "sdb.sa-east-1.amazonaws.com"
)
