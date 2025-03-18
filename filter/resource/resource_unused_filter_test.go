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
	"reflect"
	"testing"
)

func TestNewResourceUnusedFilter(t *testing.T) {
	type args struct {
		resourceGroupIDs []string
	}
	tests := []struct {
		name string
		args args
		want ResourceUnusedFilter
	}{
		{
			"Normal case: Create an instance of the ResourceUnusedFilter structure (arguments: empty array)",
			args{[]string{}},
			ResourceUnusedFilter{[]string{}},
		},
		{
			"Normal case: Create an instance of the ResourceUnusedFilter structure (arguments: non-empty array)",
			args{[]string{"aa", "bb"}},
			ResourceUnusedFilter{[]string{"aa", "bb"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceUnusedFilter(tt.args.resourceGroupIDs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceUnusedFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceUnusedFilter_FilterByCondition(t *testing.T) {
	type fields struct {
		TargetResourceGroupIDs []string
	}
	type args struct {
		record map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"Normal case: if detected is false, the result is false",
			fields{[]string{}},
			args{
				map[string]any{"detected": false},
			},
			false,
		},
		{
			"Normal case: if the value of the state element under the device.status element is not Enabled, the result is false",
			fields{[]string{}},
			args{
				map[string]any{
					"device":   map[string]any{"status": map[string]any{"state": "bbb", "health": "ccc"}},
					"detected": true,
				},
			},
			false,
		},
		{
			"Normal case: if the value of the health element under the device.status element is not OK, the result is false",
			fields{[]string{}},
			args{
				map[string]any{
					"device":   map[string]any{"status": map[string]any{"state": "Enabled", "health": "ccc"}},
					"detected": true,
				},
			},
			false,
		},
		{
			"Normal case: if the value of the available element under the annotation element is not true, the result is false",
			fields{[]string{}},
			args{
				map[string]any{
					"device":     map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}},
					"annotation": map[string]any{"available": false},
					"detected":   true,
				},
			},
			false,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field is empty and the links element is not empty, the result is false",
			fields{[]string{}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"aaa"}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"aaa"},
				},
			},
			false,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field is empty and the links element is empty, the result is true",
			fields{[]string{}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"aaa"},
				},
			},
			true,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field has a value that matches the input data and the links element is not empty, the result is false",
			fields{[]string{"aaa"}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{"aaa"}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"aaa"},
				},
			},
			false,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field has a value that matches the input data and the links element is empty, the result is true",
			fields{[]string{"aaa"}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"aaa"},
				},
			},
			true,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field has a value that not matches the input data and the links element is empty, the result is false",
			fields{[]string{"aaa"}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}, "links": []any{}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"bbb"},
				},
			},
			false,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field has a value that does not match the input data and links elements nothing, the result is false",
			fields{[]string{"bbb"}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"aaa"},
				},
			},
			false,
		},
		{
			"Normal case: if detected, status, and available elements are normal, and the TargetResourceGroupIDs field has a value that matches the input data and links elements nothing, the result is true",
			fields{[]string{"bbb"}},
			args{
				map[string]any{
					"device":           map[string]any{"status": map[string]any{"state": "Enabled", "health": "OK"}},
					"annotation":       map[string]any{"available": true},
					"detected":         true,
					"resourceGroupIDs": []string{"bbb"},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruf := ResourceUnusedFilter{
				TargetResourceGroupIDs: tt.fields.TargetResourceGroupIDs,
			}
			if got := ruf.FilterByCondition(tt.args.record); got != tt.want {
				t.Errorf("resourceUnusedFilter.FilterByCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
