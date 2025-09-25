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
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/database"
	cmapi_repository "github.com/project-cdim/configuration-manager/repository"

	"github.com/apache/age/drivers/golang/age"
	"github.com/gin-gonic/gin"

	dapr "github.com/dapr/go-sdk/client"
)

// Structure for registration information
type resourceRegister struct {
	resource []map[string]any
}

// Structure for storing resource information when fetching the list of existing resources
type existingResource struct {
	isNotDetected    bool
	resourceType     hwResourceType
	resourceGroupIDs []string
}

// Structure for storing node or switch information when fetching the list of existing nodes or switches
type existingNodeSwitch struct {
	isNotDetected    bool
	deviceDictionary map[string]hwResourceType
}

// unitResources represents the relationship between a unit device and its associated resources.
type unitResources struct {
	unitDeviceID      string
	resourceDeviceIDs []string
}

// newUnitResources creates a unitResources instance based on the provided
// request resource and non-removable device IDs. It extracts the device ID from the
// request resource and determines the related device IDs based on the resource type.
//
// For processor-type resources (Accelerator, CPU, DSP, FPGA, GPU, UnknownProcessor),
// it includes both the unit device ID and all non-removable device IDs in the
// relatedDeviceIDs slice. If no non-removable device IDs are provided, only the
// unit device ID is included regardless of resource type.
//
// Parameters:
//   - requestResource: A map containing resource information including "deviceID" and "type"
//   - nonRemovableDeviceIDs: A slice of device IDs that cannot be removed
//
// Returns:
//   - unitResources: A configured relation object with unit device ID and related device IDs
func newUnitResources(requestResource map[string]any, nonRemovableDeviceIDs []string) unitResources {
	res := unitResources{}
	res.unitDeviceID = requestResource["deviceID"].(string)
	res.resourceDeviceIDs = []string{}

	if len(nonRemovableDeviceIDs) == 0 {
		res.resourceDeviceIDs = append(res.resourceDeviceIDs, res.unitDeviceID)
	} else {
		resourceType := hwResourceType(requestResource["type"].(string))
		switch resourceType {
		case Accelerator, CPU, DSP, FPGA, GPU, UnknownProcessor:
			res.resourceDeviceIDs = append(res.resourceDeviceIDs, res.unitDeviceID)
			res.resourceDeviceIDs = append(res.resourceDeviceIDs, nonRemovableDeviceIDs...)
		}
	}

	return res
}

// isRegisterable checks whether the unit resource relation can be registered.
// It returns true if there are related device IDs associated with this unit resource relation,
// false otherwise. A unit resource relation is considered registerable only when it has
// at least one related device ID.
func (urr *unitResources) isRegisterable() bool {
	return len(urr.resourceDeviceIDs) != 0
}

// resourceTypeList is a list of resource types.
// Note: The type has been changed from [...]string to []any to allow expanding
// this list as variadic arguments when needed.
var resourceTypeList = []any{
	DB_CPU,
	DB_Accelerator,
	DB_DSP,
	DB_FPGA,
	DB_GPU,
	DB_UnknownProcessor,
	DB_Memory,
	DB_Storage,
	DB_NetworkInterface,
	DB_GraphicController,
	DB_VirtualMedia,
}

// Parts of the Cypher query to fetch specific resources
const queryResourceList_match_return string = `
MATCH (vrs:%s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)`

const queryResourceList_unionall string = `
UNION ALL`

// getQueryResourceList generates a Cypher query to retrieve a list of resources.
// It iterates over a list of resource types, formatting each into a part of the query using a predefined template.
// These parts are then joined together using a UNION ALL clause to combine the results from different resource types into a single list.
func getQueryResourceList() string {
	items := []string{}
	for range resourceTypeList {
		items = append(items, queryResourceList_match_return)
	}
	return strings.Join(items, queryResourceList_unionall)
}

const selectDeviceListColumnCount = 3
const (
	selectDeviceListIndexDeviceID = iota
	selectDeviceListIndexType
	selectDeviceListIndexResourceGroupIDs
)

// cypher query to search node
const cypherSelectNodeList string = `
	MATCH (vnd:Node)-[:Compose]->(vrs)
	WITH vnd, vrs
	ORDER BY vnd.id
	RETURN vnd.id, vrs.deviceID, vrs.type
`
const selectNodeListColumnCount = 3
const (
	selectNodeListIndexNodeID = iota
	selectNodeListIndexDeviceID
	selectNodeListIndexType
)

// cypher query to search switch
const cyperSelectSwitchList string = `
	MATCH (vcx:CXLswitch)
	OPTIONAL MATCH (vcx)-[:Connect]->(vrs)
	WITH vcx, vrs
	ORDER BY vcx.id
	RETURN vcx.id, CASE WHEN vrs.deviceID IS NULL THEN "" ELSE vrs.deviceID END, CASE WHEN vrs.type IS NULL THEN "" ELSE vrs.type END
`
const selectSwitchListColumnCount = 3
const (
	selectSwitchListIndexSwitchID = iota
	selectSwitchListIndexDeviceID
	selectSwitchListIndexType
)

// cypher query to merge resource
const cyperMergeResource = `
	MERGE (vrs:%s {deviceID: '%s'})
	SET vrs = %s
`

const (
	mergeColumnCount  = 0
	deleteColumnCount = 0
)

// cypher query to create annotation vertex and have edge
const cypherCreateAnnotation = `
	MATCH (vrs:%s {deviceID: '%s'})
	CREATE (van:Annotation {available: true})
	CREATE (vrs)-[:Have]->(van)
`

// cypher query to delete notDetected edge from resource vertex
const cypherDeleteResourceNotdetectedEdge = `
	MATCH (:%s {deviceID: '%s'})-[endt:NotDetected]->(:NotDetectedDevice)
	DELETE endt
`

// cypher query to create notDetected edge from resource vertex
const cypherCreateResourceNotdetectedEdge = `
	MATCH (vrs:%s {deviceID: '%s'}), (vndd:NotDetectedDevice)
	CREATE (vrs)-[:NotDetected]->(vndd)
`

// cypher query to delete notDetected edge from node vertex
const cypherDeleteNodeNotdetectedEdge = `
	MATCH (:Node {id: '%s'})-[endt:NotDetected]->(:NotDetectedDevice)
	DELETE endt
`

// cypher query to delete notDetected edge from switch vertex
const cypherDeleteSwitchNotdetectedEdge = `
	MATCH (:CXLswitch {id: '%s'})-[endt:NotDetected]->(:NotDetectedDevice)
	DELETE endt
`

// cypher query to merge node
const cypherMergeNode = `
	MERGE (vnd:Node {id: '%s'})
	SET vnd = {id: '%s'}
`

// cypher query to merge switch
const cypherMergeSwitch = `
	MERGE (vcx:CXLswitch {id: '%s'})
	SET vcx = {id: '%s'}
`

// cypher query to delete compose edge
const cypherDeleteComposeEdge = `
	MATCH (:Node {id: '%s'})-[ecm:Compose]->()
	DELETE ecm
`

// cypher query to delete connect edge
const cypherDeleteConnectEdge = `
	MATCH (:CXLswitch {id: '%s'})-[ecn:Connect]->()
	DELETE ecn
`

// cypher query to create compose edge
const cypherCreateComposeEdge = `
	MATCH (vrs:%s {deviceID: '%s'}), (vnd:Node {id: '%s'})
	CREATE (vnd)-[:Compose]->(vrs)
`

// cypher query to create connect edge
const cypherCreateConnectEdge = `
	MATCH (vrs:%s {deviceID: '%s'}), (vcx:CXLswitch {id: '%s'})
	CREATE (vcx)-[:Connect]->(vrs)
`

// cypher query to delete node if it does'nt have at least one compose edge
const cypherDeleteNodeWithoutEdges = `
	MATCH (vnd:Node)
	OPTIONAL MATCH (vnd:Node)-[ecm:Compose]->() WITH vnd, count(ecm) AS edges
	WHERE edges = 0
	DETACH DELETE vnd
`

// cypher query to create include edge
const cypherCreateIncludeEdge = `
	MATCH (vrs:%s {deviceID: '%s'}), (vrsg:ResourceGroups {id: '%s'})
	CREATE (vrsg)-[:Include]->(vrs)
`

// cypher query to merge Unit vertex, Annotation vertex, Have edge, and delete Contain edge.
const cypherMergeUnitAndDeleteContain = `
	MERGE (vut:Unit {deviceID: '%s'})
	MERGE (vut)-[:Have]->(:Annotation {available: true})
	WITH vut
	MATCH (vut)-[ect:Contain]->()
	DELETE ect
`

// cypher query to create Unit vertex, Annotation vertex, Have edge, and Contain edge.
const cypherCreateContain = `
	MATCH (vut:Unit {deviceID: '%s'})
	MATCH %s
	CREATE %s
`

const (
	// Parts of the Cypher query to create Contain edge
	cypherCreateContainMatchParts = `(vrs%d:%s {deviceID: '%s'})`

	// Parts of the Cypher query to create Contain edge
	cypherCreateContainCreateParts = `(vut)-[:Contain]->(vrs%d)`
)

// DaprClient is an interface that defines the methods for publishing events to Dapr.
// This interface is used to enable dependency injection for testing purposes.
type DaprClient interface {
	PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...dapr.PublishEventOption) error
}

// DaprClientImpl is a wrapper around the official Dapr client.
// This is the production implementation.
type DaprClientImpl struct {
	client dapr.Client
}

// PublishEvent publishes an event to the specified pubsub and topic using the Dapr client.
func (d *DaprClientImpl) PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...dapr.PublishEventOption) error {
	return d.client.PublishEvent(ctx, pubsubName, topicName, data, opts...)
}

// DaprClientFactory is a function type that creates a new DaprClient.
// This can be replaced with a mock implementation during testing.
var DaprClientFactory = func() (DaprClient, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}
	return &DaprClientImpl{client: client}, nil
}

// RegisterDevice registers multiple device information in the configuration management database (DB).
// It starts by logging the beginning of the process and obtaining a DB connection.
// The function then begins a transaction and defers its disconnection to ensure the DB connection is properly managed.
// It retrieves lists of already registered resources, nodes, and switches to prevent duplicate registrations.
// The function reads the JSON payload from the request body, validates it, and converts it into a structured format for processing.
// If any part of the process encounters an error, the function logs the error, rolls back the transaction if necessary,
// and returns an appropriate error response to the client.
// Upon successful validation and processing of the registration data, the function attempts to register the new devices,
// committing the transaction upon success.
// Finally, it constructs a response object containing the count and IDs of the registered devices, marshals it into JSON,
// logs the response for debugging purposes, and returns it to the client with a 201 Created status.
// If the JSON marshaling fails, it logs the error and returns an error response.
func RegisterDevice(c *gin.Context) {
	common.Log.Info(fmt.Sprintf("%s[%s] start.", c.Request.URL.Path, c.Request.Method))
	funcName := "RegisterDevice"

	// Get DB connection
	cmdb := database.NewCmDb()
	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		errorDatial := "CmDbBeginTransaction error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}
	defer cmdb.CmDbDisconnection()

	// Get the list of already registered resources
	existsResources, err := getDeviceIDList(cmdb.Tx)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "getDeviceIDList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	// Get the list of already registered nodes
	existsNodes, err := getNodeList(cmdb.Tx)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "getNodeList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	// Get the list of already registered switches
	existsSwitches, err := getCxlSwitchList(cmdb.Tx)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "getCxlSwitchList error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
		return
	}

	// Read the JSON from the RequestBody and expand it into a variable in the form of an array of Maps
	registerDevieces, err := unmarshalRequestBodyForSlice(c)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "unmarshalRequestBodyForSlice error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Read the array form of Maps converted from the JSON of the RequestBody and store it in the registration information structure
	requestResources, err := validateRegisterData(registerDevieces)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "validateRegisterData error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}

	// Compare the list of already registered resources with the JSON of the RequestBody and synchronize the entire content of the RequestBody with the DB
	registerIdList, err := registerResources(cmdb.Tx, existsResources, existsNodes, existsSwitches, requestResources)
	if err != nil {
		cmdb.CmDbRollback()
		errorDatial := "registerResources error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusBadRequest, convertErrorResponse(http.StatusBadRequest, errorDatial))
		return
	}
	cmdb.CmDbCommit()

	res := map[string]any{
		"count":     len(registerIdList),
		"deviceIDs": registerIdList,
	}

	// Log output of responseBody
	logResponseBody(res)
	common.Log.Info(fmt.Sprintf("%s[%s] completed successfully.", c.Request.URL.Path, c.Request.Method))

	ctx := context.Background()
	client, err := DaprClientFactory()
	if err != nil {
		errorDatial := "publish generate error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
	}
	if err := client.PublishEvent(ctx, "configuration_manager_hwsync", "configuration_manager.hwsync.completed", nil); err != nil {
		errorDatial := "publish error"
		common.Log.Error(fmt.Sprintf("%s %s : %s", funcName, errorDatial, err.Error()), false)
		c.JSON(http.StatusInternalServerError, convertErrorResponse(http.StatusInternalServerError, errorDatial))
	}

	c.JSON(http.StatusCreated, res)
}

// getDeviceIDList retrieves a list of existing device IDs from the database.
// It logs the query being executed for debugging purposes and initializes a map to store the results.
// The function executes a Cypher query using the provided transaction and the predefined graph name.
// If the query execution fails, it returns an error indicating the failure to execute the resource list query.
// It iterates over the results of the query, extracting the device ID, resource type, and associated resource group IDs for each record.
// These details are stored in the map with the device ID as the key.
// The function sets the initial detection status of each device to true, indicating that the device has been detected.
// After processing all records, it closes the cursor and returns the map of existing devices.
// If an error occurs while processing the results, it returns an error indicating the failure to load the resource list cursor.
func getDeviceIDList(tx *sql.Tx) (map[string]existingResource, error) {
	query := getQueryResourceList()
	common.Log.Debug(fmt.Sprintf("query: %s, params: %v", query, resourceTypeList))
	res := map[string]existingResource{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectDeviceListColumnCount, query, resourceTypeList...)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}
	defer cypherCursor.Close()

	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}

		// From one record of the search result, get DeviceID, type, and associated resource group IDs
		deviceID := cmapi_repository.ExtractEntityString(row[selectDeviceListIndexDeviceID].(*age.SimpleEntity))
		resourceType := cmapi_repository.ExtractEntityString(row[selectDeviceListIndexType].(*age.SimpleEntity))
		resourceGroupIDs := cmapi_repository.ExtractEntitySlice(row[selectDeviceListIndexResourceGroupIDs].(*age.SimpleEntity))
		// The initial value of isNotDetected is "true: detected" (change to "false: not detected" when checking existence and it was detected)
		res[deviceID] = existingResource{isNotDetected: true, resourceType: hwResourceType(resourceType), resourceGroupIDs: resourceGroupIDs}
	}

	return res, nil
}

// getNodeList retrieves a list of nodes and their associated devices from the database.
// It logs the query being executed for debugging purposes and initializes a map to store the node information.
// The function executes a Cypher query to fetch the node list, handling any errors that occur during execution.
// It iterates over the query results, extracting the node ID, associated device ID, and device type for each record.
// Nodes and their devices are stored in a map, with the node ID as the key.
// If a node ID is encountered for the first time, it is added to the map along with its associated device.
// If the node ID already exists in the map, the new device is added to the node's device dictionary.
// The function ensures that each node is marked as detected by setting the isNotDetected flag to false.
// After processing all records, the cursor is closed and the map of nodes is returned.
// If an error occurs while processing the results, an appropriate error is returned.
func getNodeList(tx *sql.Tx) (map[string]existingNodeSwitch, error) {
	common.Log.Debug(fmt.Sprintf("query: %s", cypherSelectNodeList))
	res := map[string]existingNodeSwitch{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectNodeListColumnCount, cypherSelectNodeList)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}
	defer cypherCursor.Close()

	preNodeID := ""
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}

		nodeID := cmapi_repository.ExtractEntityString(row[selectNodeListIndexNodeID].(*age.SimpleEntity))
		deviceID := cmapi_repository.ExtractEntityString(row[selectNodeListIndexDeviceID].(*age.SimpleEntity))
		deviceType := cmapi_repository.ExtractEntityString(row[selectNodeListIndexType].(*age.SimpleEntity))
		if preNodeID == "" || preNodeID != nodeID {
			res[nodeID] = existingNodeSwitch{
				isNotDetected: false,
				deviceDictionary: map[string]hwResourceType{
					deviceID: hwResourceType(deviceType),
				},
			}
		} else {
			res[nodeID].deviceDictionary[deviceID] = hwResourceType(deviceType)
		}

		preNodeID = nodeID
	}

	return res, nil
}

// getCxlSwitchList retrieves a list of registered CXL switches and their associated devices from the database.
// This function is crucial for managing CXL switches and ensuring they are correctly identified along with their connected devices.
// It begins by logging the Cypher query being executed for debugging purposes and initializes a map to store the CXL switch information.
// The function then executes the Cypher query to fetch the list of CXL switches, handling any errors that occur during the execution process.
// As it iterates over the query results, it extracts the CXL switch ID, associated device ID, and device type for each record.
// The function ensures that each CXL switch and its devices are stored in a map, with the CXL switch ID serving as the key.
// If a CXL switch ID is encountered for the first time, it is added to the map along with its associated device.
// If the CXL switch ID already exists in the map, the new device is added to the switch's device dictionary.
// This process allows for a comprehensive mapping of CXL switches to their devices, facilitating easier management and access.
// After processing all records, the cursor is closed, and the map of CXL switches is returned.
// If an error occurs while processing the results, an appropriate error message is returned to indicate the failure.
func getCxlSwitchList(tx *sql.Tx) (map[string]existingNodeSwitch, error) {
	common.Log.Debug(fmt.Sprintf("query: %s", cyperSelectSwitchList))
	res := map[string]existingNodeSwitch{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectSwitchListColumnCount, cyperSelectSwitchList)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}
	defer cypherCursor.Close()

	preCxlSwitchID := ""
	for cypherCursor.Next() {
		row, err := cypherCursor.GetRow()
		if err != nil {
			common.Log.Error(err.Error())
			return nil, err
		}

		// From one record of the search result, get CXL Switch ID, the CXL Switch's associated deviceID, and type
		// It is assumed that the records retrieved by Cypher are sorted in ascending order of CXL Switch ID
		cxlSwitchID := cmapi_repository.ExtractEntityString(row[selectSwitchListIndexSwitchID].(*age.SimpleEntity))
		deviceID := cmapi_repository.ExtractEntityString(row[selectSwitchListIndexDeviceID].(*age.SimpleEntity))
		deviceType := cmapi_repository.ExtractEntityString(row[selectSwitchListIndexType].(*age.SimpleEntity))
		if deviceID == "" || deviceType == "" {
			// If there are no resources associated with the CXL Switch, create a map of CXL Switch without resource information
			res[cxlSwitchID] = existingNodeSwitch{isNotDetected: false, deviceDictionary: map[string]hwResourceType{}}
		} else {
			if preCxlSwitchID == "" || preCxlSwitchID != cxlSwitchID {
				// If the corresponding CXL Switch is to be newly registered in the map
				// The initial value of isNotDetected is "false: detected", and the switch always remains false (this may change in the future when switch information can be obtained from HW control functions)
				res[cxlSwitchID] = existingNodeSwitch{
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						deviceID: hwResourceType(deviceType),
					},
				}
			} else {
				// In the loop, if the corresponding CXL Switch is already registered in the map
				res[cxlSwitchID].deviceDictionary[deviceID] = hwResourceType(deviceType)
			}
		}

		preCxlSwitchID = cxlSwitchID
	}

	return res, nil
}

// validateRegisterData takes an array of maps representing unmarshalled request bodies and validates each for the presence and type of mandatory fields such as deviceID and type.
// This function is essential for ensuring that the data being registered meets the required format and contains all necessary information.
// It iterates through each element in the input array, checking for the existence and data type of the "deviceID" and "type" fields.
// If either the "deviceID" or "type" field is missing or if their data types are not strings, an error is returned.
// This validation step is crucial as it enforces data integrity and prevents the registration of incomplete or malformed device information.
// Upon successful validation of all elements, the function returns a pointer to a resourceRegister struct populated with the validated data.
func validateRegisterData(body []map[string]any) (*resourceRegister, error) {
	resourceRegister := resourceRegister{}

	for index := range body {
		resource := body[index]

		valDeviceID, isDeviceIDExact := resource["deviceID"]
		// The ["deviceID"] element must exist
		if !isDeviceIDExact {
			return nil, fmt.Errorf("JSON value check exist error [deviceID]. resourceIndex(%d), resource(%v)", index, resource)
		}
		// The ["deviceID"] must be of string type
		if reflect.ValueOf(valDeviceID).Kind() != reflect.String {
			return nil, fmt.Errorf("JSON value check type error [deviceID]. resourceIndex(%d), resource(%v)", index, resource)
		}
		valType, isTypeExact := resource["type"]
		// The ["type"] element must exist
		if !isTypeExact {
			return nil, fmt.Errorf("JSON value check exist error. [type]. resourceIndex(%d), resource(%v)", index, resource)
		}
		// The ["type"] must be of string type
		if reflect.ValueOf(valType).Kind() != reflect.String {
			return nil, fmt.Errorf("JSON value check type error. [type]. resourceIndex(%d), resource(%v)", index, resource)
		}
		resourceRegister.resource = append(resourceRegister.resource, resource)
	}

	return &resourceRegister, nil
}

// registerResources processes and registers multiple hardware resources in the configuration management database.
// This function orchestrates the complete registration workflow for hardware resources including CPUs, GPUs,
// memory, storage devices, and other components, managing their relationships with nodes and CXL switches.
//
// The registration process involves several key phases:
//  1. Resource Processing: For each resource in the request, merges resource and annotation vertices,
//     establishes relationships, and manages detection states.
//  2. Topology Management: Maps resources to their associated nodes and CXL switches, ensuring proper
//     hardware topology representation and removing outdated associations.
//  3. State Synchronization: Updates the detection state of existing resources and synchronizes
//     unit-level relationships between composite devices.
//  4. Graph Maintenance: Synchronizes node and CXL switch vertices, manages their edges, and performs
//     cleanup operations to maintain database integrity.
//
// Parameters:
//   - tx: Database transaction for atomic operations across all registration steps
//   - dbExistsResources: Map of existing resources indexed by device ID, used to track detection states
//   - dbExistsNodes: Map of existing nodes and their associated devices, maintaining node topology
//   - dbExistsSwitches: Map of existing CXL switches and their connected devices, maintaining switch topology
//   - requestResources: Validated resource registration data containing device information to register
//
// Returns:
//   - []string: List of device IDs that were successfully registered during this operation
//   - error: Any error encountered during the registration process, causing transaction rollback
//
// The function ensures data consistency through transaction management and maintains the integrity
// of the hardware topology graph by properly managing vertex and edge relationships.
func registerResources(
	tx *sql.Tx,
	dbExistsResources map[string]existingResource,
	dbExistsNodes map[string]existingNodeSwitch,
	dbExistsSwitches map[string]existingNodeSwitch,
	requestResources *resourceRegister,
) ([]string, error) {
	// Return list for successfully registered IDs
	registerIdList := []string{}

	for _, requestResource := range requestResources.resource {
		deviceID := requestResource["deviceID"].(string)
		resourceType := hwResourceType(requestResource["type"].(string))

		// Merge of resource Vertex and annotation Vertex
		// Also performing the following at the same time
		// - Creating Have Edge that connects resource and annotation Vertex
		// - Deleting NotDetected Edge that connects resource and NotDetectedDevice Vertex
		err := mergeResource(tx, deviceID, resourceType, requestResource, dbExistsResources)
		if err != nil {
			return nil, err
		}

		// Check if the obtained requestID exists in dbExistsResources
		updateResourcesAsDetected(dbExistsResources, deviceID, resourceType)

		// Check if the node mentioned in links exists in dbExistsNodes
		nodeID := mappingNodes(requestResource, dbExistsNodes)
		// If the specified deviceID exists in nodes other than the specified nodeID (or in all nodes if nodeID is not specified), delete the specified deviceID information from the target node
		deleteDeviceIDFromOtherNodeSwitches(deviceID, nodeID, dbExistsNodes)

		// Check if DeviceSwitchInfo information exists in existSwitchData
		switchID := mappingSwitches(requestResource, dbExistsSwitches)
		// If the specified deviceID exists in switches other than the specified switchID (or in all switches if switchID is not specified), delete the specified deviceID information from the target switch
		deleteDeviceIDFromOtherNodeSwitches(deviceID, switchID, dbExistsSwitches)

		// Set the registered resource information in the return list
		registerIdList = append(registerIdList, deviceID)
	}

	// Loop through the list in dbExistsResources where isNotDetected is true
	for deviceID, existingResource := range dbExistsResources {
		// Reflect the NotDetected state of the resource in the DB
		err := syncNotDetectedResource(tx, deviceID, existingResource)
		if err != nil {
			return nil, err
		}
	}

	for _, requestResource := range requestResources.resource {
		err := mergeUnit(tx, requestResource, dbExistsResources)
		if err != nil {
			return nil, err
		}
	}

	// Merge and logically delete node Vertex based on the information in dbExistsNodes
	for nodeID, existingNode := range dbExistsNodes {
		// Reflect the node's Vertex and Edge in the DB
		err := syncNode(tx, nodeID, existingNode)
		if err != nil {
			return nil, err
		}
	}

	// Physically delete the node Vertex (Target for deletion: Nodes that do not have any Compose Edge connected)
	// Reason for physical deletion: Since nodes without any linked resources will not be reused, physical deletion is performed to prevent unnecessary nodes from remaining.
	common.Log.Debug(fmt.Sprintf("query: %s", cypherDeleteNodeWithoutEdges))
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteNodeWithoutEdges)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	// Merge and logically delete switch Vertex based on the information in dbExistsSwitches
	for switchID, existingSwitch := range dbExistsSwitches {

		// Reflect the switch's Vertex and Edge in the DB
		err := syncSwitch(tx, switchID, existingSwitch)
		if err != nil {
			return nil, err
		}
	}

	return registerIdList, nil
}

// updateResourcesAsDetected updates the dbExistsResources map for a given deviceID and resourceType.
// If the deviceID already exists in the map, it sets the isNotDetected field to false,
// indicating the resource has been detected. If the deviceID does not exist, it creates
// a new entry with isNotDetected set to false and assigns the provided resourceType.
//
// Parameters:
//   - dbExistsResources: map of device IDs to existingResource structs, representing current resources.
//   - deviceID: the unique identifier for the device/resource to update or add.
//   - resourceType: the type of hardware resource associated with the deviceID.
func updateResourcesAsDetected(dbExistsResources map[string]existingResource, deviceID string, resourceType hwResourceType) {
	// If it exists, change isNotDetected in dbExistsResources to false: not detected
	existResData, ok := dbExistsResources[deviceID]
	if ok {
		// Exists
		existResData.isNotDetected = false
		dbExistsResources[deviceID] = existResData
	} else {
		// Does not exist, create a new entry
		dbExistsResources[deviceID] = existingResource{
			isNotDetected: false,
			resourceType:  resourceType,
		}
	}
}

// mappingNodes analyzes the links information within a given resource to determine its associated node and updates the node's resource mapping accordingly.
// This function is critical for maintaining the topology of resources and their connections within a network or system. It processes each resource, identifying its node based on the resource type and links information, and then updates the mapping of resources to nodes in the database.
//
// The function operates as follows:
// - For CPU resources, it uses the CPU's deviceID as the nodeID directly, reflecting a self-referential node association.
// - For other resource types, it extracts the nodeID from the 'deviceID' field within the first element of the 'links' array, which represents the connection to a CPU or another pivotal resource.
// - If the nodeID is already present in the dbExistsNodes map, it adds the current resource to the node's resource dictionary.
// - If the nodeID is not present, it creates a new entry in the dbExistsNodes map with the current resource.
//
// Parameters:
// - requestResource: A map representing a single resource, including its deviceID, type, and links information.
// - dbExistsNodes: A map of existing nodes, where each key is a nodeID and each value is an existingNodeSwitch struct containing node details and a dictionary of associated resources.
//
// Returns:
// - nodeID: The identifier of the node associated with the current resource. It returns an empty string if the node cannot be determined.
//
// This function is essential for constructing an accurate and dynamic representation of the network or system topology, ensuring that resources are correctly associated with their respective nodes.
func mappingNodes(
	requestResource map[string]any,
	dbExistsNodes map[string]existingNodeSwitch,
) (nodeID string) {
	deviceID := requestResource["deviceID"].(string)
	resourceType := hwResourceType(requestResource["type"].(string))

	links, ok := requestResource["links"]
	if !ok {
		return
	}
	if reflect.ValueOf(links).Kind() != reflect.Slice {
		return
	}
	linkAnyList := links.([]any)
	if len(linkAnyList) <= 0 {
		return
	}

	// Obtain nodeID
	switch resourceType {
	case CPU:
		// For CPU, use its own deviceID as the nodeID
		nodeID = deviceID
	default:
		linkMap, ok := linkAnyList[0].(map[string]any)
		if !ok {
			return
		}

		// For resources other than CPU, the deviceID written in the links of the resource is the CPU's DeviceID, which is the nodeID
		nodeID, ok = linkMap["deviceID"].(string)
		if !ok {
			return
		}
	}

	existNodeData, ok := dbExistsNodes[nodeID]
	if ok {
		// If the node being processed exists in the existing DB or request node information, add it as resource information belonging to the node
		existNodeData.deviceDictionary[deviceID] = hwResourceType(resourceType)
		dbExistsNodes[nodeID] = existNodeData
	} else {
		// For a new node, create resource information belonging to the node
		dbExistsNodes[nodeID] = existingNodeSwitch{
			isNotDetected: false,
			deviceDictionary: map[string]hwResourceType{
				deviceID: hwResourceType(resourceType),
			},
		}
	}

	return
}

// mappingSwitches processes the 'deviceSwitchInfo' from the requestResource to identify and update the switch a device is connected to.
// This function is essential for maintaining the network topology by ensuring that devices are correctly associated with their respective switches.
// It examines the 'deviceSwitchInfo' field in the requestResource map to extract the switchID. If the switchID is valid and the switch exists in the dbExistsSwitches map,
// the device is added to the switch's device dictionary. If the switch does not exist, a new entry for the switch is created in the dbExistsSwitches map.
//
// Parameters:
// - requestResource: A map representing a single resource, including its deviceID, type, and deviceSwitchInfo.
// - dbExistsSwitches: A map of existing switches, where each key is a switchID and each value is an existingNodeSwitch struct containing switch details and a dictionary of associated devices.
//
// Returns:
// - switchID: The identifier of the switch associated with the current resource. It returns an empty string if the switch cannot be determined from the 'deviceSwitchInfo'.
//
// This function plays a critical role in the dynamic management of network resources, ensuring that the association between devices and switches is accurately reflected in the system.
func mappingSwitches(
	requestResource map[string]any,
	dbExistsSwitches map[string]existingNodeSwitch,
) (switchID string) {
	deviceID := requestResource["deviceID"].(string)
	resourceType := hwResourceType(requestResource["type"].(string))

	deviceSwitchInfo, deviceSwitchInfoOk := requestResource["deviceSwitchInfo"]
	if deviceSwitchInfoOk {
		var ok bool
		switchID, ok = deviceSwitchInfo.(string)
		if ok && len(switchID) > 0 {
			existSwitchData, ok := dbExistsSwitches[switchID]
			if ok {
				// If the switch being processed exists in the existing DB or request switch information, add it as resource information belonging to the switch
				existSwitchData.deviceDictionary[deviceID] = hwResourceType(resourceType)
				dbExistsSwitches[switchID] = existSwitchData
			} else {
				// For a new switch, create resource information belonging to the switch
				dbExistsSwitches[switchID] = existingNodeSwitch{
					isNotDetected: false,
					deviceDictionary: map[string]hwResourceType{
						deviceID: hwResourceType(resourceType),
					},
				}
			}
		}
	}

	return
}

// deleteDeviceIDFromOtherNodeSwitches removes a device from all nodes or switches except for a specified one.
// This function iterates through the dbExists map, which contains the existing nodes and switches, identified by their IDs.
// If a node or switch ID does not match the excludedNodeSwitchID, the function attempts to delete the deviceID from its deviceDictionary.
// This operation ensures that a device is only associated with the specified node or switch and is removed from all others, maintaining the integrity of the network topology.
//
// Parameters:
// - deviceID: The unique identifier of the device to be removed.
// - excludedNodeSwitchID: The ID of the node or switch from which the device should not be removed.
// - dbExists: A map where the key is the node or switch ID and the value is an existingNodeSwitch struct, representing the current network topology.
//
// The function does not return any value. It directly modifies the dbExists map by removing the specified device from all nodes or switches except the specified one.
func deleteDeviceIDFromOtherNodeSwitches(
	deviceID string,
	excludedNodeSwitchID string,
	dbExists map[string]existingNodeSwitch,
) {
	for nodeSwitchID, exists := range dbExists {
		if nodeSwitchID == excludedNodeSwitchID {
			continue
		}
		delete(exists.deviceDictionary, deviceID)
	}
}

// syncNotDetectedResource updates the database to reflect the not detected state of a resource.
// This function is responsible for managing the state of resources in the database, specifically focusing on resources that are not detected.
// It performs two main operations if the resource is marked as not detected:
// 1. Deletes an existing edge between the resource vertex and the NotDetectedDevice vertex if such an edge exists.
// 2. Creates a new edge between the resource vertex and the NotDetectedDevice vertex to indicate the resource is not detected.
//
// The function uses Cypher queries to interact with the graph database, constructing queries based on the resource type and device ID.
// It logs the Cypher queries for debugging purposes and executes them using the age.ExecCypher function.
//
// Parameters:
// - tx: A *sql.Tx transaction associated with the current database operation.
// - deviceID: The unique identifier of the device associated with the resource.
// - dbExistingResource: An existingResource struct containing details about the resource, including its not detected state and resource type.
//
// Returns:
// - An error if the operation fails at any point, including errors in converting the resource type to a database label, deleting the existing edge, or creating a new edge.
//
// This function ensures that the database accurately reflects the detection state of resources, which is crucial for maintaining the integrity of the network topology.
func syncNotDetectedResource(tx *sql.Tx, deviceID string, dbExistingResource existingResource) error {
	if dbExistingResource.isNotDetected {
		resourceType := dbExistingResource.resourceType
		label, err := resourceType.convertToDBLabel()
		if err != nil {
			return err
		}
		// If the resource Vertex in the check result list and the NotDetectedDevice Vertex are already connected by an Edge, delete that Edge once
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherDeleteResourceNotdetectedEdge, label, deviceID))
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteResourceNotdetectedEdge, label, deviceID)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}

		// Connect the resource Vertex in the check result list and the NotDetectedDevice Vertex with an Edge
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherCreateResourceNotdetectedEdge, label, deviceID))
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherCreateResourceNotdetectedEdge, label, deviceID)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}
	return nil
}

// syncNode reflects the state of a node in the database by updating its vertex and edges.
// This function performs several key operations to ensure the database accurately represents the current state of a node within the network:
// 1. Deletes the NotDetected edge that connects the node vertex to the NotDetectedDevice vertex, if such an edge exists.
// 2. Merges the node vertex to ensure it exists in the database with the correct properties.
// 3. Deletes any existing Compose edges associated with the node vertex, which represent connections to other resources.
// 4. Creates new Compose edges between the node and each of its resources, based on the existingNode.deviceDictionary.
//
// These operations are performed using Cypher queries, which are constructed and executed within the function. The function logs each query for debugging purposes.
//
// Parameters:
// - tx: A *sql.Tx transaction associated with the current database operation.
// - nodeID: The unique identifier of the node being synchronized.
// - existingNode: An existingNodeSwitch struct representing the current state of the node, including its resources.
//
// Returns:
// - An error if any operation fails, including errors from deleting edges, merging the node vertex, or creating new edges.
//
// This function is crucial for maintaining the integrity and accuracy of the network topology represented in the database.
func syncNode(tx *sql.Tx, nodeID string, existingNode existingNodeSwitch) error {
	// Delete the NotDetected Edge that connects the Node Vertex and the NotDetectedDevice Vertex
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", cypherDeleteNodeNotdetectedEdge, nodeID))
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteNodeNotdetectedEdge, nodeID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Merge the Node Vertex
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherMergeNode, nodeID, nodeID))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherMergeNode, nodeID, nodeID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Delete the Compose Edge associated with the Node Vertex
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", cypherDeleteComposeEdge, nodeID))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteComposeEdge, nodeID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Connect the Node and Resource Vertices with a Compose Edge
	for deviceID, resourceType := range existingNode.deviceDictionary {
		label, err := resourceType.convertToDBLabel()
		if err != nil {
			return err
		}
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", cypherCreateComposeEdge, label, deviceID, nodeID))
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateComposeEdge, label, deviceID, nodeID)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}
	return nil
}

// syncSwitch updates the database to reflect the current state of a switch, including its connections.
// This function ensures that the database accurately represents the switch's state by performing several operations:
// 1. Deletes the NotDetected edge between the switch vertex and the NotDetectedDevice vertex, if it exists.
// 2. Merges the switch vertex to ensure it exists in the database with the correct properties.
// 3. Deletes all Connect edges associated with the switch vertex to remove outdated connections.
// 4. Creates new Connect edges between the switch and its connected devices, based on the existingSwitch.deviceDictionary.
//
// These operations use Cypher queries to interact with the graph database. The function logs each query for debugging purposes and executes them to update the database.
//
// Parameters:
// - tx: A *sql.Tx transaction associated with the current database operation.
// - switchID: The unique identifier of the switch being synchronized.
// - existingSwitch: An existingNodeSwitch struct representing the current state of the switch, including its connections.
//
// Returns:
// - An error if any operation fails, including errors from deleting edges, merging the switch vertex, or creating new edges.
//
// This function plays a crucial role in maintaining the integrity and accuracy of the network topology represented in the database.
//
// Reflect the Switch's Vertex and Edge in the DB
func syncSwitch(tx *sql.Tx, switchID string, existingSwitch existingNodeSwitch) error {
	// Delete the NotDetected Edge that connects the Switch Vertex and the NotDetectedDevice Vertex
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", cypherDeleteSwitchNotdetectedEdge, switchID))
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteSwitchNotdetectedEdge, switchID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Merge the Switch Vertex
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherMergeSwitch, switchID, switchID))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherMergeSwitch, switchID, switchID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Delete all Connect Edges associated with the Switch Vertex
	// After deletion, reattach all necessary Edges
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", cypherDeleteConnectEdge, switchID))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteConnectEdge, switchID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Connect the Node and Switch Vertices with a Connect Edge
	for deviceID, resourceType := range existingSwitch.deviceDictionary {
		label, err := resourceType.convertToDBLabel()
		if err != nil {
			return err
		}
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", cypherCreateConnectEdge, label, deviceID, switchID))
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateConnectEdge, label, deviceID, switchID)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}
	return nil
}

// mergeResource synchronizes the state of a resource in the database with its current state.
// This function performs several key operations to ensure the database accurately reflects the resource's state:
// 1. Deletes the Have edge between the Resource vertex and the Annotation vertex, if it exists.
// 2. Merges the Resource vertex with updated properties from the requestResource map.
// 3. Conditionally merges the Annotation vertex if the resource did not exist during the last hardware sync.
// 4. Registers a new Have edge between the Resource and Annotation vertices.
// 5. Deletes the NotDetected edge for the resource, if it exists.
// 6. Conditionally creates a new Include edge from the "default group" to the resource if the resource did not exist during the last hardware sync.
//
// These operations involve executing Cypher queries to interact with the graph database, and the function logs each query for debugging purposes.
//
// Parameters:
// - tx: A *sql.Tx transaction associated with the current database operation.
// - deviceID: The unique identifier of the device associated with the resource.
// - resourceType: An hwResourceType enum value representing the type of the resource.
// - requestResource: A map containing the properties of the resource to be merged into the database.
// - dbExistsResources: A map of existingResource structs representing resources that existed during the last hardware sync.
//
// Returns:
// - An error if any operation fails, including errors from deleting edges, merging vertices, or creating new edges.
//
// This function is crucial for maintaining the accuracy and integrity of the resource information stored in the database.
//
// Merge Resource Vertex and Annotation Vertex
func mergeResource(tx *sql.Tx, deviceID string, resourceType hwResourceType, requestResource map[string]any, dbExistsResources map[string]existingResource) error {
	label, err := resourceType.convertToDBLabel()
	if err != nil {
		return err
	}

	// Merge the Resource Vertex
	property, err := common.Map2CypherProperty(requestResource)
	if err != nil {
		return err
	}
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", cyperMergeResource, label, deviceID, property))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cyperMergeResource, label, deviceID, property)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	if _, ok := dbExistsResources[deviceID]; !ok {
		// For initial registration, create the Annotation Vertex
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherCreateAnnotation, label, deviceID))
		_, err := age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateAnnotation, label, deviceID)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}

	// Delete the NotDetected Edge
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s", cypherDeleteResourceNotdetectedEdge, label, deviceID))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherDeleteResourceNotdetectedEdge, label, deviceID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	if _, ok := dbExistsResources[deviceID]; !ok {
		// For initial registration, create the Include Edge with the default group
		common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", cypherCreateIncludeEdge, label, deviceID, common.DefaultGroupId))
		_, err := age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateIncludeEdge, label, deviceID, common.DefaultGroupId)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}

	return nil
}

// mergeUnit merges unit resource data from a request with existing database resources.
// It identifies non-removable device IDs from the request, creates a unit resource relation,
// and registers the unit graph in the database if the relation is not registerable.
//
// Parameters:
//   - tx: Database transaction for executing operations
//   - requestResource: Map containing unit resource data from the request
//   - dbExistsResources: Map of existing resources in the database indexed by string keys
//
// Returns:
//   - error: Any error that occurred during the merge operation, nil if successful
func mergeUnit(tx *sql.Tx, requestResource map[string]any, dbExistsResources map[string]existingResource) error {
	nonRemovableDeviceIDs := getNonRemovableDeviceIds(requestResource)

	unitResources := newUnitResources(requestResource, nonRemovableDeviceIDs)
	if unitResources.isRegisterable() {
		err := registerUnitGraph(tx, unitResources, dbExistsResources)
		if err != nil {
			return err
		}
	}

	return nil
}

// getNonRemovableDeviceIds extracts the list of non-removable device IDs from the given requestResource map.
// It expects the following nested structure in requestResource:
//   - "constraints": map[string]any (optional; absence is not treated as a warning)
//   - "constraints/nonRemovableDevices": []any (optional; absence is not treated as a warning)
//   - "constraints/nonRemovableDevices" each element: map[string]any with a "deviceID" key (string)
//
// If any part of the structure is missing or invalid (except for missing "constraints" or "nonRemovableDevices"),
// it logs a warning or debug message and returns an empty slice.
// Returns a slice of device IDs (strings) for all valid non-removable devices found.
func getNonRemovableDeviceIds(requestResource map[string]any) []string {
	constraints, ok := requestResource["constraints"]
	if !ok {
		// [Normal] The absence of the constraints element is expected, so it is treated as normal. Output a Debug log and return an empty slice
		common.Log.Debug(fmt.Sprintf("constraints field is missing. resource(%v)", requestResource))
		return []string{}
	}

	constraintsMap, ok := constraints.(map[string]any)
	if !ok {
		// [Warning] The constraints element not being a Map is unexpected, so it is treated as a Warning. Output a Warning log and return an empty slice
		common.Log.Warn(fmt.Sprintf("constraints field is not a map. resource(%v)", requestResource))
		return []string{}
	}

	nonRemovableDevices, ok := constraintsMap["nonRemovableDevices"]
	if !ok {
		// [Normal] The absence of the nonRemovableDevices element is expected, so it is treated as normal. Output a Debug log and return an empty slice
		common.Log.Debug(fmt.Sprintf("constraints/nonRemovableDevices field is missing. resource(%v)", requestResource))
		return []string{}
	}

	nonRemovableDevicesSlice, ok := nonRemovableDevices.([]any)
	if !ok {
		// [Warning] The constraints element not being a Slice is unexpected, so it is treated as a Warning. Output a Warning log and return an empty slice
		common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices field is not a list. resource(%v)", requestResource))
		return []string{}
	}

	if len(nonRemovableDevicesSlice) == 0 {
		// [Warning] The constraints element being a Slice but having 0 elements is unexpected, so it is treated as a Warning. Output a Warning log and return an empty slice
		common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices contains no elements. resource(%v)", requestResource))
		return []string{}
	}

	res := []string{}
	for i, device := range nonRemovableDevicesSlice {
		deviceMap, ok := device.(map[string]any)
		if !ok {
			// [Warning] An element under the constraints element (Slice) not being a Map is unexpected, so output a Warning log and Skip
			common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices[%d] field is not a map. resource(%v)", i, requestResource))
			continue
		}
		deviceID, ok := deviceMap["deviceID"].(string)
		if !ok {
			// [Warning] The absence of "Key:deviceID" in the Map under the constraints element (Slice) is unexpected, so output a Warning log and Skip
			common.Log.Warn(fmt.Sprintf("constraints/nonRemovableDevices[%d]/deviceID field is missing or not a string. resource(%v)", i, requestResource))
			continue
		}
		res = append(res, deviceID)
	}

	return res
}

// registerUnitGraph creates or updates a unit node in the graph database and establishes
// containment relationships with its associated resources.
//
// The function performs the following operations:
// 1. Merges a unit node and deletes existing containment relationships
// 2. Creates new containment relationships between the unit and its resources
//
// Parameters:
//   - tx: Database transaction for executing graph operations
//   - unitResources: Contains unit device ID and associated resource information
//   - dbExistsResources: Map of existing resources in the database to avoid duplicates
//
// Returns:
//   - error: nil on success, otherwise an error describing what went wrong
//
// The function uses Cypher queries to interact with the Apache AGE graph database
// and logs debug information for query execution and error details.
func registerUnitGraph(tx *sql.Tx, unitResources unitResources, dbExistsResources map[string]existingResource) error {
	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s", cypherMergeUnitAndDeleteContain, unitResources.unitDeviceID))
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherMergeUnitAndDeleteContain, unitResources.unitDeviceID)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	matches, creates, err := createContainQuery(unitResources, dbExistsResources)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	common.Log.Debug(fmt.Sprintf("query: %s, param1: %s, param2: %s, param3: %s", cypherCreateContain, unitResources.unitDeviceID, matches, creates))
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateContain, unitResources.unitDeviceID, matches, creates)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	return nil
}

// createContainQuery constructs Cypher query parts for creating "contain" relationships
// between a unit and its related devices in a Neo4j database.
//
// It takes a unitResources containing device IDs and a map of existing database resources,
// then generates MATCH and CREATE clauses for each valid related device.
//
// Parameters:
//   - unitResources: contains the list of related device IDs to process
//   - dbExistsResources: map of device ID to existing resource information
//
// Returns:
//   - string: comma-separated MATCH clauses for Cypher query
//   - string: comma-separated CREATE clauses for Cypher query
//   - error: any error encountered during label conversion
//
// The function skips devices that don't exist in dbExistsResources and logs warnings.
// If a resource type cannot be converted to a database label, an error is returned.
func createContainQuery(unitResources unitResources, dbExistsResources map[string]existingResource) (string, string, error) {
	matches := make([]string, 0, len(unitResources.resourceDeviceIDs))
	creates := make([]string, 0, len(unitResources.resourceDeviceIDs))

	for i, relatedDeviceID := range unitResources.resourceDeviceIDs {
		resource, ok := dbExistsResources[relatedDeviceID]
		if !ok {
			// This condition should not be reached unless the "nonRemovableDevices" in the resource received from HWControl contains an invalid deviceID
			common.Log.Warn(fmt.Sprintf("Resource with deviceID %s does not exist in dbExistsResources", relatedDeviceID))
			continue
		}

		label, err := resource.resourceType.convertToDBLabel()
		if err != nil {
			// This error should not occur because invalid resource types are already checked when registering resources in this API
			return "", "", err
		}
		matches = append(matches, fmt.Sprintf(cypherCreateContainMatchParts, i, label, relatedDeviceID))
		creates = append(creates, fmt.Sprintf(cypherCreateContainCreateParts, i))
	}

	return strings.Join(matches, ", "), strings.Join(creates, ", "), nil
}
