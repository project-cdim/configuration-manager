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
	"golang.org/x/exp/slices"
)

// ResourceUnusedFilter is a struct that holds the filter criteria for resources that are unused.
// It contains a slice of resource group IDs that are targeted for the search.
type ResourceUnusedFilter struct {
	TargetResourceGroupIDs []string // TargetResourceGroupIDs is a slice of resource group IDs to be targeted for search.
}

// NewResourceUnusedFilter creates a new instance of resourceUnusedFilter.
// It takes a slice of resource group IDs as input and returns a resourceUnusedFilter
// with the provided resource group IDs set as the TargetResourceGroupIDs.
//
// Parameters:
//
//	resourceGroupIDs - a slice of strings representing the resource group IDs to be targeted.
//
// Returns:
//
//	A new instance of resourceUnusedFilter with the TargetResourceGroupIDs set to the provided resource group IDs.
func NewResourceUnusedFilter(resourceGroupIDs []string) ResourceUnusedFilter {
	return ResourceUnusedFilter{
		TargetResourceGroupIDs: resourceGroupIDs,
	}
}

// FilterByCondition filters a record based on specific conditions defined in the resourceUnusedFilter.
// It checks if the record meets the following criteria:
// 1. The "detected" field must be true.
// 2. The "device" field must contain a "status" map with "state" as "Enabled" and "health" as "OK".
// 3. The "annotation" field must contain an "available" boolean set to true.
// 4. The links array within the device field should be empty.
// 5. If TargetResourceGroupIDs is not empty, the record's "resourceGroupIDs" must contain at least one of the target IDs.
//
// Parameters:
//
//	record - a map[string]any representing the record to be filtered.
//	recordOption - optional additional parameters (not used in this function).
//
// Returns:
//
//	A boolean value indicating whether the record meets the filter conditions.
func (ruf ResourceUnusedFilter) FilterByCondition(record map[string]any, recordOption ...any) bool {
	detected := record["detected"].(bool)
	if !detected {
		return false
	}

	device := record["device"].(map[string]any)
	status, _ := device["status"].(map[string]any)
	if !isEnableStatus(status) {
		return false
	}

	annotation := record["annotation"].(map[string]any)
	available, _ := annotation["available"].(bool)
	if !available {
		return false
	}

	links, _ := device["links"].([]any)
	if !isUnused(links) {
		return false
	}

	if len(ruf.TargetResourceGroupIDs) > 0 {
		resourceGroupIDs := record["resourceGroupIDs"].([]string)
		for _, resourceGroupID := range resourceGroupIDs {
			if slices.Contains(ruf.TargetResourceGroupIDs, resourceGroupID) {
				return true
			}
		}

		return false
	}

	return true
}
