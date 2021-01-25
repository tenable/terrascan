/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package utils

import (
	"strings"

	"github.com/hashicorp/terraform/registry/regsrc"
)

var (
	supportedRemoteTypes = []string{"git", "s3", "gcs", "http", "terraform-registry"}
	localSourcePrefixes  = []string{
		"./",
		"../",
		".\\",
		"..\\",
	}
)

// IsValidRemoteType validates the remote type supplied as scan option
func IsValidRemoteType(remoteType string) bool {
	for _, r := range supportedRemoteTypes {
		if strings.EqualFold(r, strings.ToLower(strings.TrimSpace(remoteType))) {
			return true
		}
	}
	return false
}

// IsLocalSourceAddr validates if a source address is a local address or not
func IsLocalSourceAddr(addr string) bool {
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}

// IsRegistrySourceAddr will validate if the source address is a valid registry
// module or not.
// a valid source address is of the form <HOSTNAME>/NAMESPACE>/<NAME>/<PROVIDER>
// regsrc.ParseModuleSource func returns a terraform registry module source.
func IsRegistrySourceAddr(addr string) bool {
	_, err := regsrc.ParseModuleSource(addr)
	return err == nil
}
