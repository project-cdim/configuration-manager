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
        
package filter

// A struct that is defined for situations where no condition is required.
type noFilter struct{}

// NewNoFilter is a constructor for filters that don't require conditions.
// It creates and returns an instance of a filter that does not perform any filtering based on conditions.
// This type of filter can be used in scenarios where all records are to be included without any filtering.
//
// Returns:
// An instance of noFilter, which does not apply any conditions to the records.
func NewNoFilter() noFilter {
	return noFilter{}
}

// FilterByCondition always returns true, indicating that any given record matches the conditions.
// This implementation of the filter does not hold any conditions, thus it does not perform any actual filtering.
// It is used in contexts where filtering based on conditions is not required.
//
// Parameters:
// - record: A map representing the record to be checked. This parameter is not used in this implementation.
// - recordOption: Optional parameters that can be used to extend the logic. These are not used in this implementation.
//
// Returns:
// - Always returns true, as this filter does not perform any condition checks.
func (nsc noFilter) FilterByCondition(record map[string]any, recordOption ...any) bool {
	return true
}
