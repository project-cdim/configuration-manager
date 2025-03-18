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
        
package resource_repository

import (
	"reflect"
	"testing"

	"github.com/project-cdim/configuration-manager/common"
)

func TestNewUpdateGroupOfResourceRepository(t *testing.T) {
	type args struct {
		deviceID          string
		dbDeviceType      string
		newResourceGroups []string
	}
	tests := []struct {
		name string
		args args
		want AssignResourceToGroupRepository
	}{
		{
			"Normal case: Create an instance of the AssignResourceToGroupRepository struct",
			args{"001", "CPU", []string{common.DefaultGroupId}},
			AssignResourceToGroupRepository{
				"001",
				"CPU",
				[]string{common.DefaultGroupId},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAssignResourceToGroupRepository(tt.args.deviceID, tt.args.dbDeviceType, tt.args.newResourceGroups); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAssignResourceToGroupRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignResourceToGroupRepository_Set(t *testing.T) {
	t.Skip("not test")
}
