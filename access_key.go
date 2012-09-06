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

package aws

// AccessKey represents an AWS access key, including its ID and associated
// secret.
//
// To manage your access keys, go here:
//
//     https://portal.aws.amazon.com/gp/aws/securityCredentials
//
type AccessKey struct {
	// The key's ID.
	//
	// For example:
	//
	//     AKIAIOSFODNN7EXAMPLE
	//
	Id string

	// The key's associated secret.
	//
	// For example:
	//
	//     wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
	//
	Secret string
}
