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
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_newUnitResources(t *testing.T) {
	type args struct {
		requestResource       map[string]any
		nonRemovableDeviceIDs []string
	}
	tests := []struct {
		name string
		args args
		want unitResources
	}{
		{
			name: "empty nonRemovableDeviceIDs - should include only unit device ID",
			args: args{
				requestResource: map[string]any{
					"deviceID": "unit123",
					"type":     "CPU",
				},
				nonRemovableDeviceIDs: []string{},
			},
			want: unitResources{
				unitDeviceID:      "unit123",
				resourceDeviceIDs: []string{"unit123"},
			},
		},
		{
			name: "nil nonRemovableDeviceIDs - should include only unit device ID",
			args: args{
				requestResource: map[string]any{
					"deviceID": "unit456",
					"type":     "GPU",
				},
				nonRemovableDeviceIDs: nil,
			},
			want: unitResources{
				unitDeviceID:      "unit456",
				resourceDeviceIDs: []string{"unit456"},
			},
		},
		{
			name: "processor type (CPU) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "cpu001",
					"type":     "CPU",
				},
				nonRemovableDeviceIDs: []string{"dev001", "dev002"},
			},
			want: unitResources{
				unitDeviceID:      "cpu001",
				resourceDeviceIDs: []string{"cpu001", "dev001", "dev002"},
			},
		},
		{
			name: "processor type (Accelerator) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "acc001",
					"type":     "Accelerator",
				},
				nonRemovableDeviceIDs: []string{"mem001"},
			},
			want: unitResources{
				unitDeviceID:      "acc001",
				resourceDeviceIDs: []string{"acc001", "mem001"},
			},
		},
		{
			name: "processor type (DSP) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "dsp001",
					"type":     "DSP",
				},
				nonRemovableDeviceIDs: []string{"dev003", "dev004", "dev005"},
			},
			want: unitResources{
				unitDeviceID:      "dsp001",
				resourceDeviceIDs: []string{"dsp001", "dev003", "dev004", "dev005"},
			},
		},
		{
			name: "processor type (FPGA) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "fpga001",
					"type":     "FPGA",
				},
				nonRemovableDeviceIDs: []string{"ctrl001"},
			},
			want: unitResources{
				unitDeviceID:      "fpga001",
				resourceDeviceIDs: []string{"fpga001", "ctrl001"},
			},
		},
		{
			name: "processor type (GPU) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "gpu001",
					"type":     "GPU",
				},
				nonRemovableDeviceIDs: []string{"vram001", "bus001"},
			},
			want: unitResources{
				unitDeviceID:      "gpu001",
				resourceDeviceIDs: []string{"gpu001", "vram001", "bus001"},
			},
		},
		{
			name: "processor type (UnknownProcessor) with nonRemovableDeviceIDs - should include unit and non-removable devices",
			args: args{
				requestResource: map[string]any{
					"deviceID": "proc001",
					"type":     "UnknownProcessor",
				},
				nonRemovableDeviceIDs: []string{"cache001"},
			},
			want: unitResources{
				unitDeviceID:      "proc001",
				resourceDeviceIDs: []string{"proc001", "cache001"},
			},
		},
		{
			name: "non-processor type (Memory) with nonRemovableDeviceIDs - should include only unit device ID",
			args: args{
				requestResource: map[string]any{
					"deviceID": "mem001",
					"type":     "Memory",
				},
				nonRemovableDeviceIDs: []string{"dev001", "dev002"},
			},
			want: unitResources{
				unitDeviceID:      "mem001",
				resourceDeviceIDs: []string{},
			},
		},
		{
			name: "non-processor type (Storage) with nonRemovableDeviceIDs - should include only unit device ID",
			args: args{
				requestResource: map[string]any{
					"deviceID": "storage001",
					"type":     "Storage",
				},
				nonRemovableDeviceIDs: []string{"controller001"},
			},
			want: unitResources{
				unitDeviceID:      "storage001",
				resourceDeviceIDs: []string{},
			},
		},
		{
			name: "non-processor type (NetworkInterface) with nonRemovableDeviceIDs - should include only unit device ID",
			args: args{
				requestResource: map[string]any{
					"deviceID": "nic001",
					"type":     "NetworkInterface",
				},
				nonRemovableDeviceIDs: []string{"port001", "port002"},
			},
			want: unitResources{
				unitDeviceID:      "nic001",
				resourceDeviceIDs: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUnitResources(tt.args.requestResource, tt.args.nonRemovableDeviceIDs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUnitResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unitResourceRelation_isRegisterable(t *testing.T) {
	tests := []struct {
		name string
		urr  unitResources
		want bool
	}{
		{
			name: "empty relatedDeviceIDs - should return false",
			urr: unitResources{
				unitDeviceID:      "unit001",
				resourceDeviceIDs: []string{},
			},
			want: false,
		},
		{
			name: "nil relatedDeviceIDs - should return false",
			urr: unitResources{
				unitDeviceID:      "unit002",
				resourceDeviceIDs: nil,
			},
			want: false,
		},
		{
			name: "single relatedDeviceID - should return true",
			urr: unitResources{
				unitDeviceID:      "unit003",
				resourceDeviceIDs: []string{"dev001"},
			},
			want: true,
		},
		{
			name: "multiple relatedDeviceIDs - should return true",
			urr: unitResources{
				unitDeviceID:      "unit004",
				resourceDeviceIDs: []string{"dev001", "dev002", "dev003"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urr.isRegisterable(); got != tt.want {
				t.Errorf("unitResourceRelation.isRegisterable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getQueryResourceList(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Normal case",
			queryResourceList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getQueryResourceList(); got != tt.want {
				t.Errorf("getQueryResourceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterDevice(t *testing.T) {
	t.Skip("not test")
}

func Test_getDeviceIDList(t *testing.T) {
	t.Skip("not test")
}

func Test_getNodeList(t *testing.T) {
	t.Skip("not test")
}

func Test_getCxlSwitchList(t *testing.T) {
	t.Skip("not test")
}

func Test_validateRegisterData(t *testing.T) {
	type args struct {
		body []map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    *resourceRegister
		wantErr bool
	}{
		{
			"Error case: [deviceID] element does not exist",
			args{
				[]map[string]any{
					{"aa": "id12", "type": Accelerator},
				},
			},
			nil,
			true,
		},
		{
			"Error case: [deviceID] element is not a string",
			args{
				[]map[string]any{
					{"deviceID": 123, "type": Accelerator},
				},
			},
			nil,
			true,
		},
		{
			"Error case: [type] element does not exist",
			args{
				[]map[string]any{
					{"deviceID": "id12", "aa": Accelerator},
				},
			},
			nil,
			true,
		},
		{
			"Error case: [type] element is not a string",
			args{
				[]map[string]any{
					{"deviceID": "id12", "type": 123},
				},
			},
			nil,
			true,
		},
		{
			"Normal case",
			args{
				[]map[string]any{
					{"deviceID": "id12", "type": Accelerator},
					{"deviceID": "id13", "type": DSP},
					{"deviceID": "id14", "type": FPGA},
					{"deviceID": "id11", "type": CPU},
					{"deviceID": "id17", "type": Memory},
					{"deviceID": "id99", "type": CPU},
					{"deviceID": "id15", "type": GPU},
					{"deviceID": "id16", "type": UnknownProcessor},
					{"deviceID": "id88", "type": Memory},
					{"deviceID": "id18", "type": Storage},
					{"deviceID": "id19", "type": NetworkInterface},
					{"deviceID": "id20", "type": GraphicController},
					{"deviceID": "id21", "type": VirtualMedia},
				},
			},
			createTestValue_resourceRegister(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateRegisterData(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRegisterData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateRegisterData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerResources(t *testing.T) {
	t.Skip("not test")
}

func Test_updateResourcesAsDetected(t *testing.T) {
	type args struct {
		dbExistsResources map[string]existingResource
		deviceID          string
		resourceType      hwResourceType
	}
	tests := []struct {
		name string
		args args
		want map[string]existingResource
	}{
		{
			name: "deviceID exists, isNotDetected should become false",
			args: args{
				dbExistsResources: map[string]existingResource{
					"dev1": {isNotDetected: true, resourceType: hwResourceType(CPU)},
				},
				deviceID:     "dev1",
				resourceType: hwResourceType(CPU),
			},
			want: map[string]existingResource{
				"dev1": {isNotDetected: false, resourceType: hwResourceType(CPU)},
			},
		},
		{
			name: "deviceID does not exist, should add new entry",
			args: args{
				dbExistsResources: map[string]existingResource{
					"dev1": {isNotDetected: true, resourceType: hwResourceType(CPU)},
				},
				deviceID:     "dev2",
				resourceType: hwResourceType(GPU),
			},
			want: map[string]existingResource{
				"dev1": {isNotDetected: true, resourceType: hwResourceType(CPU)},
				"dev2": {isNotDetected: false, resourceType: hwResourceType(GPU)},
			},
		},
		{
			name: "empty dbExistsResources, should add new entry",
			args: args{
				dbExistsResources: map[string]existingResource{},
				deviceID:          "dev3",
				resourceType:      hwResourceType(Memory),
			},
			want: map[string]existingResource{
				"dev3": {isNotDetected: false, resourceType: hwResourceType(Memory)},
			},
		},
		{
			name: "deviceID exists with isNotDetected already false, should remain false",
			args: args{
				dbExistsResources: map[string]existingResource{
					"dev4": {isNotDetected: false, resourceType: hwResourceType(Storage)},
				},
				deviceID:     "dev4",
				resourceType: hwResourceType(Storage),
			},
			want: map[string]existingResource{
				"dev4": {isNotDetected: false, resourceType: hwResourceType(Storage)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateResourcesAsDetected(tt.args.dbExistsResources, tt.args.deviceID, tt.args.resourceType)
			if !reflect.DeepEqual(tt.args.dbExistsResources, tt.want) {
				t.Errorf("updateResourcesAsDetected() = %v, want %v", tt.args.dbExistsResources, tt.want)
			}
		})
	}
}

func Test_mappingNodes(t *testing.T) {
	type args struct {
		requestResource map[string]any                // Target registration resource
		dbExistsNodes   map[string]existingNodeSwitch // Existing DB list
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]existingNodeSwitch
		want2 string
	}{
		{
			"If there is no 'links' element in the target registration resource and the deviceID of the target registration resource is not in the list, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the target registration resource has a 'links' element but it is not a slice, and the deviceID of the target registration resource is not in the list, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
					"links":    "aaaaa",
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the target registration resource has a 'links' element but it is a slice with 0 elements, and the deviceID of the target registration resource is not in the list, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
					"links":    []any{},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the type of the target registration resource is CPU, and the deviceID is a nodeID not present in the list, add the node to the list and return the nodeID",
			args{
				map[string]any{
					"deviceID": "dev11",
					"type":     "CPU",
					"links": []any{
						map[string]any{"deviceID": "devmem11", "type": "memory"},
					},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"dev11": hwResourceType(CPU)},
				},
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"dev11",
		},
		{
			"If the type of the target registration resource is CPU, and the deviceID is a nodeID that exists in the list, do not update the node in the list but add the CPU under the node and return the nodeID",
			args{
				map[string]any{
					"deviceID": "dev11",
					"type":     "CPU",
					"links": []any{
						map[string]any{"deviceID": "devmem11", "type": "memory"},
					},
				},
				map[string]existingNodeSwitch{
					"dev11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"dev11": hwResourceType(CPU)},
				},
			},
			"dev11",
		},
		{
			"If the type of the target registration resource is not CPU, and the first element of 'links' is not a map type, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
					"links":    []any{"test"},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the type of the target registration resource is not CPU, and the first element of 'links' does not have a 'deviceID' element, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
					"links": []any{
						map[string]any{"test": "dev11", "type": "CPU"},
					},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the type of the target registration resource is not CPU, and the first element of 'links' has a 'deviceID' element but it is not a string type, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
					"links": []any{
						map[string]any{"deviceID": 111, "type": "CPU"},
					},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the type of the target registration resource is not CPU, and the 'deviceID' in the 'links' is a nodeID not present in the list, add the node to the list and return the nodeID",
			args{
				map[string]any{
					"deviceID": "devmem11",
					"type":     "memory",
					"links": []any{
						map[string]any{"deviceID": "dev11", "type": "CPU"},
					},
				},
				map[string]existingNodeSwitch{
					"dev12": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devcpu11": hwResourceType(CPU)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
				"dev12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devcpu11": hwResourceType(CPU)},
				},
			},
			"dev11",
		},
		{
			"If the type of the target registration resource is not CPU, and the 'deviceID' in the 'links' is a nodeID that exists in the list, do not update the node in the list but add the Memory under the node and return the nodeID",
			args{
				map[string]any{
					"deviceID": "devmem11",
					"type":     "memory",
					"links": []any{
						map[string]any{"deviceID": "dev11", "type": "CPU"},
					},
				},
				map[string]existingNodeSwitch{
					"dev11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{},
					},
				},
			},
			map[string]existingNodeSwitch{
				"dev11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"dev11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got2 := mappingNodes(tt.args.requestResource, tt.args.dbExistsNodes)
			got := tt.args.dbExistsNodes
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mappingNodes() dbExistsNodes = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("mappingNodes() nodeID = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_mappingSwitches(t *testing.T) {
	type args struct {
		requestResource  map[string]any                // Target registration resource
		dbExistsSwitches map[string]existingNodeSwitch // Existing DB list
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]existingNodeSwitch
		want2 string
	}{
		{
			"If there is no deviceSwitchInfo element in the target registration resource, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID": "devmem12",
					"type":     "memory",
				},
				map[string]existingNodeSwitch{
					"switch11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"switch11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the value of the deviceSwitchInfo element in the target registration resource is not a string type, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID":         "devmem12",
					"type":             "memory",
					"deviceSwitchInfo": 111,
				},
				map[string]existingNodeSwitch{
					"switch11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"switch11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the value of the deviceSwitchInfo element in the target registration resource is an empty string, return an empty string without changing the list",
			args{
				map[string]any{
					"deviceID":         "devmem12",
					"type":             "memory",
					"deviceSwitchInfo": "",
				},
				map[string]existingNodeSwitch{
					"switch11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"switch11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
			},
			"",
		},
		{
			"If the 'switchID' obtained from deviceSwitchInfo does not exist in the list, add the switch to the list and return the switchID",
			args{
				map[string]any{
					"deviceID":         "devmem12",
					"type":             "memory",
					"deviceSwitchInfo": "switch12",
				},
				map[string]existingNodeSwitch{
					"switch11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
					},
				},
			},
			map[string]existingNodeSwitch{
				"switch11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem11": hwResourceType(Memory)},
				},
				"switch12": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem12": hwResourceType(Memory)},
				},
			},
			"switch12",
		},
		{
			"If the 'switchID' obtained from deviceSwitchInfo exists in the list, add the deviceID under the switch in the list and return the switchID",
			args{
				map[string]any{
					"deviceID":         "devmem12",
					"type":             "memory",
					"deviceSwitchInfo": "switch11",
				},
				map[string]existingNodeSwitch{
					"switch11": {
						isNotDetected:    false,
						deviceDictionary: map[string]hwResourceType{},
					},
				},
			},
			map[string]existingNodeSwitch{
				"switch11": {
					isNotDetected:    false,
					deviceDictionary: map[string]hwResourceType{"devmem12": hwResourceType(Memory)},
				},
			},
			"switch11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got2 := mappingSwitches(tt.args.requestResource, tt.args.dbExistsSwitches)
			got := tt.args.dbExistsSwitches
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mappingSwitches() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("mappingNodes() nodeID = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_deleteDeviceIDFromOtherNodeSwitches(t *testing.T) {
	type args struct {
		deviceID             string                        // Target deviceID to delete
		excludedNodeSwitchID string                        // ID of the node/switch to be excluded from deletion
		dbExists             map[string]existingNodeSwitch // Existing DB list of nodes/switches
	}
	tests := []struct {
		name string
		args args
		want map[string]existingNodeSwitch
	}{
		{
			"Normal case: The \"node ID to be excluded\" exists in the list, and the deviceID in the list is deleted",
			args{
				"cpu01",
				"node01",
				map[string]existingNodeSwitch{
					"node01": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
					"node02": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
				},
			},
			map[string]existingNodeSwitch{
				"node01": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"cpu01": hwResourceType(CPU),
						"mem01": hwResourceType(Memory),
					},
				},
				"node02": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"mem01": hwResourceType(Memory),
					},
				},
			},
		},
		{
			"Normal case: The \"node ID to be excluded\" exists in the list, and the deviceID in the list is not deleted",
			args{
				"cpu02",
				"node01",
				map[string]existingNodeSwitch{
					"node01": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
					"node02": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
				},
			},
			map[string]existingNodeSwitch{
				"node01": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"cpu01": hwResourceType(CPU),
						"mem01": hwResourceType(Memory),
					},
				},
				"node02": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"cpu01": hwResourceType(CPU),
						"mem01": hwResourceType(Memory),
					},
				},
			},
		},
		{
			"Normal case: The \"node ID to be excluded\" does not exist in the list, and all deviceIDs in the list are deleted",
			args{
				"cpu01",
				"node03",
				map[string]existingNodeSwitch{
					"node01": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
					"node02": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
				},
			},
			map[string]existingNodeSwitch{
				"node01": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"mem01": hwResourceType(Memory),
					},
				},
				"node02": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"mem01": hwResourceType(Memory),
					},
				},
			},
		},
		{
			"Normal case: The \"node ID to be excluded\" does not exist in the list, and the deviceIDs in the list are not deleted",
			args{
				"cpu02",
				"node03",
				map[string]existingNodeSwitch{
					"node01": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
					"node02": {
						isNotDetected: false,
						deviceDictionary: map[string]hwResourceType{
							"cpu01": hwResourceType(CPU),
							"mem01": hwResourceType(Memory),
						},
					},
				},
			},
			map[string]existingNodeSwitch{
				"node01": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"cpu01": hwResourceType(CPU),
						"mem01": hwResourceType(Memory),
					},
				},
				"node02": {
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						"cpu01": hwResourceType(CPU),
						"mem01": hwResourceType(Memory),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteDeviceIDFromOtherNodeSwitches(tt.args.deviceID, tt.args.excludedNodeSwitchID, tt.args.dbExists)
			got := tt.args.dbExists
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteDeviceIDFromOtherNodeSwitches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncNotDetectedResource(t *testing.T) {
	t.Skip("not test")
}

func Test_syncNode(t *testing.T) {
	t.Skip("not test")
}

func Test_syncSwitch(t *testing.T) {
	t.Skip("not test")
}

func Test_mergeResource(t *testing.T) {
	t.Skip("not test")
}

func Test_mergeUnit(t *testing.T) {
	t.Skip("not test")
}

func Test_getNonRemovableDeviceIds(t *testing.T) {
	type args struct {
		requestResource map[string]any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "constraints field missing",
			args: args{
				requestResource: map[string]any{
					"deviceID": "dev1",
				},
			},
			want: []string{},
		},
		{
			name: "constraints field not a map",
			args: args{
				requestResource: map[string]any{
					"constraints": "not_a_map",
				},
			},
			want: []string{},
		},
		{
			name: "nonRemovableDevices field missing",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{},
				},
			},
			want: []string{},
		},
		{
			name: "nonRemovableDevices not a list",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": "not_a_list",
					},
				},
			},
			want: []string{},
		},
		{
			name: "nonRemovableDevices empty list",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{},
					},
				},
			},
			want: []string{},
		},
		{
			name: "nonRemovableDevices contains non-map element",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{
							"not_a_map",
							map[string]any{"deviceID": "dev2"},
						},
					},
				},
			},
			want: []string{"dev2"},
		},
		{
			name: "nonRemovableDevices contains map without deviceID",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{
							map[string]any{"notDeviceID": "dev3"},
							map[string]any{"deviceID": "dev4"},
						},
					},
				},
			},
			want: []string{"dev4"},
		},
		{
			name: "nonRemovableDevices contains map with non-string deviceID",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{
							map[string]any{"deviceID": 123},
							map[string]any{"deviceID": "dev5"},
						},
					},
				},
			},
			want: []string{"dev5"},
		},
		{
			name: "multiple valid deviceIDs",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{
							map[string]any{"deviceID": "devA"},
							map[string]any{"deviceID": "devB"},
						},
					},
				},
			},
			want: []string{"devA", "devB"},
		},
		{
			name: "all valid",
			args: args{
				requestResource: map[string]any{
					"constraints": map[string]any{
						"nonRemovableDevices": []any{
							map[string]any{"deviceID": "devX"},
						},
					},
				},
			},
			want: []string{"devX"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNonRemovableDeviceIds(tt.args.requestResource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNonRemovableDeviceIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerUnitGraphics(t *testing.T) {
	t.Skip("not test")
}

func Test_createContainQuery(t *testing.T) {
	type args struct {
		unitResourceRelation unitResources
		dbExistsResources    map[string]existingResource
	}
	tests := []struct {
		name        string
		args        args
		wantMatches string
		wantCreates string
		wantErr     bool
	}{
		{
			name: "empty relatedDeviceIDs - should return empty strings",
			args: args{
				unitResourceRelation: unitResources{
					unitDeviceID:      "unit001",
					resourceDeviceIDs: []string{},
				},
				dbExistsResources: map[string]existingResource{},
			},
			wantMatches: "",
			wantCreates: "",
			wantErr:     false,
		},
		{
			name: "single related device - should create match and create parts",
			args: args{
				unitResourceRelation: unitResources{
					unitDeviceID:      "unit001",
					resourceDeviceIDs: []string{"dev001"},
				},
				dbExistsResources: map[string]existingResource{
					"dev001": {isNotDetected: false, resourceType: hwResourceType(CPU)},
				},
			},
			wantMatches: fmt.Sprintf(cypherCreateContainMatchParts, 0, "CPU", "dev001"),
			wantCreates: fmt.Sprintf(cypherCreateContainCreateParts, 0),
			wantErr:     false,
		},
		{
			name: "multiple related devices - should create multiple match and create parts",
			args: args{
				unitResourceRelation: unitResources{
					unitDeviceID:      "unit001",
					resourceDeviceIDs: []string{"dev001", "dev002", "dev003"},
				},
				dbExistsResources: map[string]existingResource{
					"dev001": {isNotDetected: false, resourceType: hwResourceType(CPU)},
					"dev002": {isNotDetected: false, resourceType: hwResourceType(Memory)},
					"dev003": {isNotDetected: false, resourceType: hwResourceType(GPU)},
				},
			},
			wantMatches: strings.Join([]string{
				fmt.Sprintf(cypherCreateContainMatchParts, 0, "CPU", "dev001"),
				fmt.Sprintf(cypherCreateContainMatchParts, 1, "Memory", "dev002"),
				fmt.Sprintf(cypherCreateContainMatchParts, 2, "GPU", "dev003"),
			}, ", "),
			wantCreates: strings.Join([]string{
				fmt.Sprintf(cypherCreateContainCreateParts, 0),
				fmt.Sprintf(cypherCreateContainCreateParts, 1),
				fmt.Sprintf(cypherCreateContainCreateParts, 2),
			}, ", "),
			wantErr: false,
		},
		{
			name: "related device not in dbExistsResources - should skip missing device",
			args: args{
				unitResourceRelation: unitResources{
					unitDeviceID:      "unit001",
					resourceDeviceIDs: []string{"dev001", "dev999", "dev002"},
				},
				dbExistsResources: map[string]existingResource{
					"dev001": {isNotDetected: false, resourceType: hwResourceType(CPU)},
					"dev002": {isNotDetected: false, resourceType: hwResourceType(Memory)},
				},
			},
			wantMatches: strings.Join([]string{
				fmt.Sprintf(cypherCreateContainMatchParts, 0, "CPU", "dev001"),
				fmt.Sprintf(cypherCreateContainMatchParts, 2, "Memory", "dev002"),
			}, ", "),
			wantCreates: strings.Join([]string{
				fmt.Sprintf(cypherCreateContainCreateParts, 0),
				fmt.Sprintf(cypherCreateContainCreateParts, 2),
			}, ", "),
			wantErr: false,
		},
		{
			name: "all related devices missing from dbExistsResources - should return empty strings",
			args: args{
				unitResourceRelation: unitResources{
					unitDeviceID:      "unit001",
					resourceDeviceIDs: []string{"dev999", "dev998"},
				},
				dbExistsResources: map[string]existingResource{
					"dev001": {isNotDetected: false, resourceType: hwResourceType(CPU)},
				},
			},
			wantMatches: "",
			wantCreates: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, gotCreates, err := createContainQuery(tt.args.unitResourceRelation, tt.args.dbExistsResources)
			if (err != nil) != tt.wantErr {
				t.Errorf("createContainQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatches != tt.wantMatches {
				t.Errorf("createContainQuery() gotMatches = %v, want %v", gotMatches, tt.wantMatches)
			}
			if gotCreates != tt.wantCreates {
				t.Errorf("createContainQuery() gotCreates = %v, want %v", gotCreates, tt.wantCreates)
			}
		})
	}
}

func createTestValue_resourceRegister() *resourceRegister {
	res := resourceRegister{resource: []map[string]any{
		{"deviceID": "id12", "type": Accelerator},
		{"deviceID": "id13", "type": DSP},
		{"deviceID": "id14", "type": FPGA},
		{"deviceID": "id11", "type": CPU},
		{"deviceID": "id17", "type": Memory},
		{"deviceID": "id99", "type": CPU},
		{"deviceID": "id15", "type": GPU},
		{"deviceID": "id16", "type": UnknownProcessor},
		{"deviceID": "id88", "type": Memory},
		{"deviceID": "id18", "type": Storage},
		{"deviceID": "id19", "type": NetworkInterface},
		{"deviceID": "id20", "type": GraphicController},
		{"deviceID": "id21", "type": VirtualMedia},
	}}
	return &res
}

const queryResourceList string = `
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)`
