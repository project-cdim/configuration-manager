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

package controller

import (
	"reflect"
	"testing"
)

func TestUpdateAnnotation(t *testing.T) {
	t.Skip("not test")
}

func TestGetNonRemovableDeviceIDs(t *testing.T) {
	testCases := []struct {
		name     string
		device   map[string]any
		expected []string
	}{
		{
			name:     "Nil Device",
			device:   nil,
			expected: []string{},
		},
		{
			name: "Missing Constraints",
			device: map[string]any{
				"some_other_field": "some_value",
			},
			expected: []string{},
		},
		{
			name: "Constraints Not a Map",
			device: map[string]any{
				"constraints": "not a map",
			},
			expected: []string{},
		},
		{
			name: "Missing NonRemovableDevices",
			device: map[string]any{
				"constraints": map[string]any{},
			},
			expected: []string{},
		},
		{
			name: "NonRemovableDevices Not a List",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": "not a list",
				},
			},
			expected: []string{},
		},
		{
			name: "Empty NonRemovableDevices",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{},
				},
			},
			expected: []string{},
		},
		{
			name: "NonRemovableDevice Not a Map",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{"not a map"},
				},
			},
			expected: []string{},
		},
		{
			name: "Missing DeviceID",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{
						map[string]any{},
					},
				},
			},
			expected: []string{},
		},
		{
			name: "DeviceID Not a String",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{
						map[string]any{
							"deviceID": 123,
						},
					},
				},
			},
			expected: []string{},
		},
		{
			name: "Valid NonRemovableDevices",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{
						map[string]any{
							"deviceID": "device1",
						},
						map[string]any{
							"deviceID": "device2",
						},
					},
				},
			},
			expected: []string{"device1", "device2"},
		},
		{
			name: "Mixed Valid and Invalid NonRemovableDevices",
			device: map[string]any{
				"constraints": map[string]any{
					"nonRemovableDevices": []any{
						map[string]any{
							"deviceID": "device1",
						},
						"not a map",
						map[string]any{
							"deviceID": 123,
						},
						map[string]any{
							"deviceID": "device2",
						},
					},
				},
			},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getNonRemovableDeviceIDs(tc.device)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("getNonRemovableDeviceIDs(%v) = %v, expected %v", tc.device, result, tc.expected)
			}
		})
	}
}
