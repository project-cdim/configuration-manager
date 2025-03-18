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
        
package group_model

import (
	"fmt"
	"unicode/utf8"

	"github.com/project-cdim/configuration-manager/common"
)

// ValidateProperty checks the validity of the provided property map.
// It ensures that the "name" field is a string with a length between 1 and 64 characters,
// and the "description" field is a string with a length of up to 256 characters.
// Returns true if both conditions are met, otherwise returns false.
//
// Parameters:
//   - property: map[string]any - A map containing the property fields to validate.
//
// Returns:
//   - bool: true if the property is valid, false otherwise.
func ValidateProperty(property map[string]any) bool {
	name, ok := property["name"].(string)
	if !ok {
		common.Log.Warn("name is not a string")
		return false
	}

	nameLen := utf8.RuneCountInString(name)
	if nameLen < 1 || nameLen > 64 {
		common.Log.Warn(fmt.Sprintf("name length is invalid. length(%v)", nameLen))
		return false
	}

	description, ok := property["description"].(string)
	if !ok {
		common.Log.Warn("description is not a string")
		return false
	}

	descLen := utf8.RuneCountInString(description)
	if descLen > 256 {
		common.Log.Warn(fmt.Sprintf("description length is invalid. length(%v)", descLen))
		return false
	}
	return true
}
