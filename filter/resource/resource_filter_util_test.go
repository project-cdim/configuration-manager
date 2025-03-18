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
	"testing"
)

func Test_isEnableStatus(t *testing.T) {
	type args struct {
		status map[string]any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Normal case: If there is no 'state' item in status, return false",
			args{
				map[string]any{"health": "NG"},
			},
			false,
		},
		{
			"Normal case: If there is no 'health' item in status, return false",
			args{
				map[string]any{"state": "Disabled"},
			},
			false,
		},
		{
			"Normal case: If [state: \"Enabled\", health: \"OK\"], return true",
			args{
				map[string]any{"state": "Enabled", "health": "OK"},
			},
			true,
		},
		{
			"Normal case: If [state: not \"Enabled\", health: \"OK\"], return false",
			args{
				map[string]any{"state": "Disabled", "health": "OK"},
			},
			false,
		},
		{
			"Normal case: If [state: \"Enabled\", health: not \"OK\"], return false",
			args{
				map[string]any{"state": "Enabled", "health": "NG"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEnableStatus(tt.args.status); got != tt.want {
				t.Errorf("isEnableStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUnused(t *testing.T) {
	type args struct {
		links []any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Normal case: Returns true when an empty slice is passed",
			args{
				[]any{},
			},
			true,
		},
		{
			"Normal case: Returns false when a non-empty slice is passed",
			args{
				[]any{"aaa"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUnused(tt.args.links); got != tt.want {
				t.Errorf("isUnused() = %v, want %v", got, tt.want)
			}
		})
	}
}
