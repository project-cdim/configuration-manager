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

	"github.com/project-cdim/configuration-manager/common"
)

// GroupList represents a list of groups.
// It contains a slice of Group objects.
type GroupList struct {
	Groups []Group
}

// NewGroupList creates and returns a new GroupList instance with an empty list of Groups.
func NewGroupList() GroupList {
	return GroupList{
		Groups: []Group{},
	}
}

// ToObject converts the GroupList to a slice of maps, where each map represents a group.
// Only valid groups (those that pass validation) are included in the result.
//
// Returns:
//
//	[]map[string]any: A slice of maps, each containing the data of a valid group.
func (gl *GroupList) ToObject() []map[string]any {
	res := []map[string]any{}
	for _, group := range gl.Groups {
		if group.Validate() {
			groupObj := group.ToObject()
			res = append(res, groupObj)
		}
	}

	return res
}

// ToObjectWithResources converts the GroupList to a slice of maps,
// where each map represents a group with its associated resources.
// Only valid groups (those that pass validation) are included in the result.
//
// Returns:
//
//	[]map[string]any: A slice of maps, each containing the data of a valid group.
func (gl *GroupList) ToObjectWithResources() []map[string]any {
	res := []map[string]any{}
	for _, group := range gl.Groups {
		if group.Validate() {
			res = append(res, group.ToObjectWithResources())
		} else {
			common.Log.Warn(fmt.Sprintf("Not added to list. group(%v)", group))
		}
	}

	return res
}
