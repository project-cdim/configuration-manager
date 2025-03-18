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

// resourceTypeList is a list of resource types.
var resourceTypeList = [...]string{
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
MATCH (vrs: %s)
WHERE exists(vrs.deviceID) AND exists(vrs.type)
OPTIONAL MATCH (vrsg)-[ein:Include]->(vrs)
RETURN vrs.deviceID, vrs.type, COLLECT(vrsg.id)`

const queryResourceList_unionall string = `
UNION ALL`

// getQueryResourceList generates a Cypher query to retrieve a list of resources.
// It iterates over a list of resource types, formatting each into a part of the query using a predefined template.
// These parts are then joined together using a UNION ALL clause to combine the results from different resource types into a single list.
func getQueryResourceList() string {
	items := []string{}
	for _, resourceType := range resourceTypeList {
		items = append(items, fmt.Sprintf(queryResourceList_match_return, resourceType))
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
	MATCH (vnd: Node)-[ecm:Compose]->(vrs)
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
	MATCH (vcx: CXLswitch)
	OPTIONAL MATCH (vcx)-[ecn:Connect]->(vrs)
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
	MERGE (vrs: %s {deviceID: '%s'})
	SET vrs = %s
`

const (
	mergeColumnCount  = 0
	deleteColumnCount = 0
)

// cypher query to create annotation vertex and have edge
const cypherCreateAnnotation = `
	MATCH (vrs: %s {deviceID: '%s'})
	CREATE (van:Annotation {available: true})
	CREATE (vrs)-[ehv: Have]->(van)
`

// cypher query to delete notDetected edge from resource vertex
const cypherDeleteResourceNotdetectedEdge = `
	MATCH (vrs: %s {deviceID: '%s'})-[endt: NotDetected]->(vndd: NotDetectedDevice) 
	DELETE endt
`

// cypher query to create notDetected edge from resource vertex
const cypherCreateResourceNotdetectedEdge = `
	MATCH (vrs: %s {deviceID: '%s'}), (vndd: NotDetectedDevice)
	CREATE (vrs)-[endt: NotDetected]->(vndd)
`

// cypher query to delete notDetected edge from node vertex
const cypherDeleteNodeNotdetectedEdge = `
	MATCH (vnd: Node {id: '%s'})-[endt: NotDetected]->(vndd: NotDetectedDevice) 
	DELETE endt
`

// cypher query to delete notDetected edge from switch vertex
const cypherDeleteSwitchNotdetectedEdge = `
	MATCH (vcx: CXLswitch {id: '%s'})-[endt: NotDetected]->(vndd: NotDetectedDevice) 
	DELETE endt
`

// cypher query to merge node
const cypherMergeNode = `
	MERGE (vnd: Node {id: '%s'})
	SET vnd = {id: '%s'}
`

// cypher query to merge switch
const cypherMergeSwitch = `
	MERGE (vcx: CXLswitch {id: '%s'})
	SET vcx = {id: '%s'}
`

// cypher query to delete compose edge
const cypherDeleteComposeEdge = `
	MATCH (vnd: Node {id: '%s'})-[ecm: Compose]->(vrs)
	DELETE ecm
`

// cypher query to delete connect edge
const cypherDeleteConnectEdge = `
	MATCH (vcx: CXLswitch {id: '%s'})-[ecn: Connect]->(vrs)
	DELETE ecn
`

// cypher query to create compose edge
const cypherCreateComposeEdge = `
	MATCH (vrs: %s {deviceID: '%s'}), (vnd: Node {id: '%s'})
	CREATE (vnd)-[ecm: Compose]->(vrs)
`

// cypher query to create connect edge
const cypherCreateConnectEdge = `
	MATCH (vrs: %s {deviceID: '%s'}), (vcx: CXLswitch {id: '%s'})
	CREATE (vcx)-[ecn: Connect]->(vrs)
`

// cypher query to delete node if it does'nt have at least one compose edge
const cypherDeleteNodeWithoutEdges = `
	MATCH (vnd: Node)
	OPTIONAL MATCH (vnd: Node)-[ecm: Compose]->(vrs) WITH vnd, count(ecm) AS edges
	WHERE edges = 0
	DETACH DELETE vnd
`

// cypher query to create include edge
const cypherCreateIncludeEdge = `
	MATCH (vrs: %s {deviceID: '%s'}), (vrsg: ResourceGroups {id: '%s'})
	CREATE (vrsg)-[ein: Include]->(vrs)
`

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
	common.Log.Debug(getQueryResourceList())
	res := map[string]existingResource{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectDeviceListColumnCount, getQueryResourceList())
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

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
	cypherCursor.Close()

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
	common.Log.Debug(cypherSelectNodeList)
	res := map[string]existingNodeSwitch{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectNodeListColumnCount, cypherSelectNodeList)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

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
	cypherCursor.Close()

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
	common.Log.Debug(cyperSelectSwitchList)
	res := map[string]existingNodeSwitch{}
	cypherCursor, err := age.ExecCypher(tx, database.GRAPH_NAME, selectSwitchListColumnCount, cyperSelectSwitchList)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

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
	cypherCursor.Close()

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

// registerResources compares the list of registered resources with the JSON of the RequestBody,
// marking existing ones while performing registration and update for resources and annotation Vertexes through Merge process.
// This function is pivotal for maintaining the integrity and up-to-dateness of the resource database.
// It takes a transaction object, maps of existing resources, nodes, and switches, and a struct containing request resources as inputs.
// The function iterates over each request resource, performing a series of checks and operations:
// - Merges resource and annotation Vertexes, creating or deleting edges as necessary.
// - Checks and updates the mapping of resources to ensure they are correctly associated with nodes and switches.
// - Deletes device IDs from nodes or switches where they no longer belong.
// - Synchronizes the state of resources, nodes, and switches with the database, reflecting any changes made.
// - Physically deletes node Vertexes that are no longer connected to any edges, to prevent clutter in the database.
// The function returns a list of successfully registered device IDs or an error if any operation fails.
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
		mappingResources(dbExistsResources, deviceID)

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
	common.Log.Debug(cypherDeleteNodeWithoutEdges)
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

// mappingResources updates the detection status of a device in the existing resources map.
// This function checks if a device, identified by its deviceID, exists in the map of existing resources.
// If the device is found, it updates the device's detection status to indicate that the device is not "not detected" (i.e., it has been detected).
// This is crucial for maintaining the accuracy of the resource tracking system, ensuring that devices are correctly marked as detected when they are present in the request.
//
// Parameters:
// - dbExistsResources: A map of existing resources where the key is the deviceID and the value is the existingResource struct.
// - deviceID: The unique identifier of the device being checked.
//
// The function does not return any value. It directly modifies the dbExistsResources map by updating the detection status of the specified device.
func mappingResources(dbExistsResources map[string]existingResource, deviceID string) {
	// If it exists, change isNotDetected in dbExistsResources to false: not detected
	existResData, deviceIDOk := dbExistsResources[deviceID]
	if deviceIDOk {
		// Exists
		existResData.isNotDetected = false
		dbExistsResources[deviceID] = existResData
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
		cypherNotDetectedEdgeDelete := fmt.Sprintf(cypherDeleteResourceNotdetectedEdge, label, deviceID)
		common.Log.Debug(cypherNotDetectedEdgeDelete)
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherNotDetectedEdgeDelete)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}

		// Connect the resource Vertex in the check result list and the NotDetectedDevice Vertex with an Edge
		cypherNotDetectedEdgeCreate := fmt.Sprintf(cypherCreateResourceNotdetectedEdge, label, deviceID)
		common.Log.Debug(cypherNotDetectedEdgeCreate)
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherNotDetectedEdgeCreate)
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
	cypherNodeNotDetectedEdgeDelete := fmt.Sprintf(cypherDeleteNodeNotdetectedEdge, nodeID)
	common.Log.Debug(cypherNodeNotDetectedEdgeDelete)
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherNodeNotDetectedEdgeDelete)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Merge the Node Vertex
	cypherNodeCreate := fmt.Sprintf(cypherMergeNode, nodeID, nodeID)
	common.Log.Debug(cypherNodeCreate)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherNodeCreate)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Delete the Compose Edge associated with the Node Vertex
	cypherComposeDelete := fmt.Sprintf(cypherDeleteComposeEdge, nodeID)
	common.Log.Debug(cypherComposeDelete)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherComposeDelete)
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
		cypherComposeEdgeCreate := fmt.Sprintf(cypherCreateComposeEdge, label, deviceID, nodeID)
		common.Log.Debug(cypherComposeEdgeCreate)
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherComposeEdgeCreate)
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
	cypherSwitchNotDetectedEdgeDelete := fmt.Sprintf(cypherDeleteSwitchNotdetectedEdge, switchID)
	common.Log.Debug(cypherSwitchNotDetectedEdgeDelete)
	_, err := age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherSwitchNotDetectedEdgeDelete)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Merge the Switch Vertex
	cypherSwitchCreate := fmt.Sprintf(cypherMergeSwitch, switchID, switchID)
	common.Log.Debug(cypherSwitchCreate)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherSwitchCreate)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	// Delete all Connect Edges associated with the Switch Vertex
	// After deletion, reattach all necessary Edges
	cypherConnectDelete := fmt.Sprintf(cypherDeleteConnectEdge, switchID)
	common.Log.Debug(cypherConnectDelete)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherConnectDelete)
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
		cypherConnectEdgeCreate := fmt.Sprintf(cypherCreateConnectEdge, label, deviceID, switchID)
		common.Log.Debug(cypherConnectEdgeCreate)
		_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherConnectEdgeCreate)
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
	cypherResourceMerge := fmt.Sprintf(cyperMergeResource, label, deviceID, property)
	common.Log.Debug(cypherResourceMerge)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherResourceMerge)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	if _, ok := dbExistsResources[deviceID]; ok {
		// If the resource existed during the last HW sync, do not update the Annotation Vertex
	} else {
		// If the resource did not exist during the last HW sync, create the Annotation Vertex
		cypherAnnotationCreate := fmt.Sprintf(cypherCreateAnnotation, label, deviceID)
		common.Log.Debug(cypherAnnotationCreate)
		_, err := age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherAnnotationCreate)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}

	// Delete the NotDetected Edge
	cypherNotDetectedEdgeDelete := fmt.Sprintf(cypherDeleteResourceNotdetectedEdge, label, deviceID)
	common.Log.Debug(cypherNotDetectedEdgeDelete)
	_, err = age.ExecCypher(tx, database.GRAPH_NAME, deleteColumnCount, cypherNotDetectedEdgeDelete)
	if err != nil {
		common.Log.Error(err.Error())
		return err
	}

	label, err = resourceType.convertToDBLabel()
	if err != nil {
		return err
	}

	if _, ok := dbExistsResources[deviceID]; ok {
		// If the resource existed during the last HW sync, do nothing (since it's already connected to the resource group and there are no changes)
	} else {
		// If the resource did not exist during the last HW sync, create a new Include Edge from "default group" to "resource"
		cypherCreateIncludeEdge := fmt.Sprintf(cypherCreateIncludeEdge, label, deviceID, common.DefaultGroupId)
		common.Log.Debug(cypherCreateIncludeEdge)
		_, err := age.ExecCypher(tx, database.GRAPH_NAME, mergeColumnCount, cypherCreateIncludeEdge)
		if err != nil {
			common.Log.Error(err.Error())
			return err
		}
	}

	return nil
}
