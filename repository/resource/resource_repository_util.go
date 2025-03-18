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
	"sort"
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	annotation_model "github.com/project-cdim/configuration-manager/model/annotation"
	resource_model "github.com/project-cdim/configuration-manager/model/resource"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"

	"github.com/apache/age/drivers/golang/age"
)

// ComposeResource assembles and returns Resource information from a single record of search results.
// If the detail argument is true, all fields are populated; if false, only key fields are populated.
func ComposeResource(resVertex *age.Vertex, annotationVertex *age.Vertex, resourceGroupIDs *age.SimpleEntity, nodeIDs *age.SimpleEntity, detected bool, detail bool) resource_model.Resource {
	resource := resource_model.NewResource()

	// Retrieve the 'available' information from the resource annotation, defaulting to true (available for use in the design proposal) if not present
	annotationProp := annotationVertex.Props()
	available := true
	if _, ok := annotationProp["available"]; ok {
		switch annotationProp["available"].(type) {
		case bool:
			available = annotationProp["available"].(bool)
		}
	}
	annotation := annotation_model.Annotation{
		Properties: map[string]any{"available": available},
	}

	// Retrieve the Property information from the resource Vertex data (Properties are obtained in map format)
	device := resVertex.Props()
	if len(device) <= 0 {
		// Return an empty model if there are no Properties
		return resource
	}
	if !detail {
		device = extractPrimaryDeviceProp(device)
	}
	// If any item value is an empty array, store an empty array
	common.Nil2EmptyFromMap(device)

	// Add device information to the resource Property
	resource.Device = device
	// Add resourceGroupIDs information to the resource Property
	resource.ResourceGroupIDs = cmapi_repository.ExtractEntitySlice(resourceGroupIDs)
	// Add available information to the resource Property
	resource.Annotation = annotation
	// Add detected information to the resource Property
	resource.Detected = detected
	// Add nodeIDs information to the resource Property
	resource.NodeIDs = cmapi_repository.ExtractEntitySlice(nodeIDs)

	return resource
}

// extractPrimaryDeviceProp returns only the primary properties of a device.
func extractPrimaryDeviceProp(prop map[string]any) map[string]any {
	res := map[string]any{}

	res["deviceID"] = prop["deviceID"]
	res["type"] = prop["type"]
	if status, ok := prop["status"]; ok {
		res["status"] = status
	}

	return res
}

// Sort by resourceType and deviceID
func sortResourceList(resources []resource_model.Resource) {
	// Sort from the less priority sort key
	sort.Slice(resources, func(i, j int) bool {
		return strings.Compare(resources[j].Device["deviceID"].(string), resources[i].Device["deviceID"].(string)) > 0
	})
	sort.SliceStable(resources, func(i, j int) bool {
		return strings.Compare(resources[j].Device["type"].(string), resources[i].Device["type"].(string)) > 0
	})
}
