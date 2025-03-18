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
        
package rack_repository

import (
	"reflect"
	"testing"
)

func TestNewRackRepository(t *testing.T) {
	type args struct {
		rackID string
		detail bool
	}
	tests := []struct {
		name string
		args args
		want RackRepository
	}{
		{
			"Normal case: Create an instance of the RackRepository struct ('rack11', true)",
			args{"rack11", true},
			RackRepository{
				"rack11",
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRackRepository(tt.args.rackID, tt.args.detail); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRackRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRackRepository_Find(t *testing.T) {
	t.Skip("not test")
}
