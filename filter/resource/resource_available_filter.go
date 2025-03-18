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

	"golang.org/x/exp/slices"
)

// ResourceAvailableFilter is a struct that holds the filter criteria for resource availability.
// It contains a slice of resource group IDs that are targeted for search.
type ResourceAvailableFilter struct {
	TargetResourceGroupIDs []string // TargetResourceGroupIDs is a slice of resource group IDs to be targeted for search.
}

// NewResourceAvailableFilter creates a new instance of resourceAvailableFilter.
// It takes a slice of resource group IDs as input and returns a resourceAvailableFilter
// with the provided resource group IDs set as the TargetResourceGroupIDs.
//
// Parameters:
//
//	resourceGroupIDs - a slice of strings representing the resource group IDs to be targeted.
//
// Returns:
//
//	A new instance of resourceAvailableFilter with the TargetResourceGroupIDs set to the provided resource group IDs.
func NewResourceAvailableFilter(resourceGroupIDs []string) ResourceAvailableFilter {
	return ResourceAvailableFilter{
		TargetResourceGroupIDs: resourceGroupIDs,
	}
}

// FilterByCondition evaluates if a given record matches the conditions set in the resourceAvailableFilter.
// It returns true if the record matches the conditions, false otherwise.
//
// This function checks for the following conditions:
// - If NotDetected is true, it filters out records where 'detected' is false.
// - If Available is true, it ensures that the 'device' and 'annotation' maps exist and that the 'available' status within 'annotation' matches the 'status' in 'device'.
// - It also checks if the record belongs to any of the TargetResourceGroupIDs, if specified.
//
// Arguments:
// record: The record to evaluate, expected to be a map with keys like 'detected', 'device', and 'annotation'.
// recordOption: Optional parameters for future use.
//
// Returns:
// A boolean indicating if the record matches the filter conditions.
func (raf ResourceAvailableFilter) FilterByCondition(record map[string]any, recordOption ...any) bool {
	detected := record["detected"].(bool)
	if !detected {
		return false
	}

	device := record["device"].(map[string]any)
	status, ok := device["status"].(map[string]any)
	if !ok {
		common.Log.Warn(fmt.Sprintf("There is no status element in rec[\"device\"], or status element is not a map. [device : %v]", device))
		return false
	}
	if !isEnableStatus(status) {
		return false
	}

	annotation := record["annotation"].(map[string]any)
	available, ok := annotation["available"].(bool)
	if !ok {
		common.Log.Warn(fmt.Sprintf("There is no available element in rec[\"annotation\"], or available element is not a map. [annotation : %v]", annotation))
		return false
	}
	if !available {
		return false
	}

	if len(raf.TargetResourceGroupIDs) > 0 {
		resourceGroupIDs := record["resourceGroupIDs"].([]string)
		for _, resourceGroupID := range resourceGroupIDs {
			if slices.Contains(raf.TargetResourceGroupIDs, resourceGroupID) {
				return true
			}
		}

		return false
	}

	return true
}
