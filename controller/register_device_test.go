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

func Test_mappingResources(t *testing.T) {
	type args struct {
		dbExistsResources map[string]existingResource // 既存DBリスト(=リソース一覧)
		deviceID          string
	}
	tests := []struct {
		name string
		args args
		want map[string]existingResource
	}{
		{
			"When the target deviceID exists in the resource list",
			args{
				map[string]existingResource{
					"dev11": {isNotDetected: true, resourceType: hwResourceType(Accelerator)},
				}, "dev11",
			},
			map[string]existingResource{
				"dev11": {isNotDetected: false, resourceType: hwResourceType(Accelerator)},
			},
		},
		{
			"When the target deviceID does not exist in the resource list",
			args{
				map[string]existingResource{
					"dev11": {isNotDetected: true, resourceType: hwResourceType(Accelerator)},
				}, "dev99",
			},
			map[string]existingResource{
				"dev11": {isNotDetected: true, resourceType: hwResourceType(Accelerator)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mappingResources(tt.args.dbExistsResources, tt.args.deviceID)
			got := tt.args.dbExistsResources
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mappingResources() = %v, want %v", got, tt.want)
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
MATCH (vrs: CPU)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: Accelerator)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: DSP)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: FPGA)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: GPU)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: UnknownProcessor)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: Memory)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: Storage)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: NetworkInterface)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: GraphicController)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)
UNION ALL
MATCH (vrs: VirtualMedia)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)`
