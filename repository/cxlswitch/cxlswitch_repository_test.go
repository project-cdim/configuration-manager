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
        
package cxlswitch_repository

import (
	"reflect"
	"testing"
)

func TestNewCXLSwitchRepository(t *testing.T) {
	type args struct {
		cxlSwitchID string
	}
	tests := []struct {
		name string
		args args
		want CXLSwitchRepository
	}{
		{
			"Normal case: Create an instance of the CXLSwitchRepository struct ('sw11')",
			args{"sw11"},
			CXLSwitchRepository{
				"sw11",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCXLSwitchRepository(tt.args.cxlSwitchID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCXLSwitchRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCXLSwitchRepository_Find(t *testing.T) {
	t.Skip("not test")
}
