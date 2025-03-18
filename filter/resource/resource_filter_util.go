// Copyright (C) 2025 NEC Corporation.
// 
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
        
package resource_filter

import (
	"fmt"

	"github.com/project-cdim/configuration-manager/common"
)

// isEnableStatus checks if the given status map has the "state" and "health" elements
// and verifies if their values are "Enabled" and "OK" respectively.
//
// If the "state" element is missing, it logs a warning and returns false.
// If the "health" element is missing, it logs a warning and returns false.
// If both "state" and "health" elements are present and their values are "Enabled" and "OK",
// it returns true. Otherwise, it returns false.
//
// Parameters:
//
//	status - a map[string]any representing the status to be checked.
//
// Returns:
//
//	A boolean value indicating whether the status is enabled and healthy.
func isEnableStatus(status map[string]any) bool {
	if _, ok := status["state"]; !ok {
		common.Log.Warn(fmt.Sprintf("There is no state element in status. [states : %v]", status))
		return false
	} else if _, ok := status["health"]; !ok {
		common.Log.Warn(fmt.Sprintf("There is no health element in status. [states : %v]", status))
		return false
	} else if status["state"] == "Enabled" && status["health"] == "OK" {
		return true
	}

	return false
}

// isUnused checks if the provided links slice is empty.
//
// This function returns true if the length of the links slice is zero,
// indicating that there are no links. Otherwise, it returns false.
//
// Parameters:
//
//	links - a slice of any type representing the links to be checked.
//
// Returns:
//
//	A boolean value indicating whether the links slice is empty.
func isUnused(links []any) bool {
	return len(links) == 0
}
