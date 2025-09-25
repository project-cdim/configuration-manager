// Copyright (C) 2025 NEC Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/project-cdim/configuration-manager/controller"
	"github.com/project-cdim/configuration-manager/database"

	dapr "github.com/dapr/go-sdk/client"
)

type daprClientMock struct{}

func (m *daprClientMock) PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...dapr.PublishEventOption) error {
	return nil // No-op implementation for testing
}

// TestMain is the entry point for running tests in this package.
func TestMain(m *testing.M) {
	// read db password from environment variable
	dbPassword := os.Getenv("GRAPH_DB_PASSWORD")
	if dbPassword == "" {
		fmt.Println("GRAPH_DB_PASSWORD environment variable is not set. exiting...")
		return
	}

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"apache/age:release_PG16_1.5.0",
		postgres.WithInitScripts(filepath.Join("testdata", "container_init.sh")),
		postgres.WithDatabase("cmdb"),
		postgres.WithUsername("cmuser"),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Printf("Could not start PostgreSQL container: %v\n", err)
		return
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			fmt.Printf("Could not terminate PostgreSQL container: %v\n", err)
		}
	}()

	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Could not get PostgreSQL container host: %v\n", err)
		return
	}

	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Printf("Could not get PostgreSQL container port: %v\n", err)
		return
	}

	// Mock the database.GetSecretCmdb function to return the test container's connection details
	originalGetSecretCmdb := database.GetSecretCmdb
	database.GetSecretCmdb = func() (database.SecretCmdb, error) {
		return database.SecretCmdb{
			Host:     pgHost,
			Port:     pgPort.Port(),
			User:     "cmuser",
			Password: dbPassword,
			Dbname:   "cmdb",
		}, nil
	}
	defer func() {
		database.GetSecretCmdb = originalGetSecretCmdb
	}()

	// Mock the controller.DaprClientFactory function to return a no-op client
	originalDaprClientFactory := controller.DaprClientFactory
	controller.DaprClientFactory = func() (controller.DaprClient, error) {
		return &daprClientMock{}, nil
	}
	defer func() {
		controller.DaprClientFactory = originalDaprClientFactory
	}()

	// Run tests
	exitCode := m.Run()

	if exitCode != 0 {
		fmt.Printf("Tests failed with exit code %d\n", exitCode)
	}
}

// TestRestAPI tests all REST API endpoints for the configuration manager.
// It sets up a test engine and runs subtests for both resource and resource group operations,
// including successful retrievals, list operations, and error cases for non-existent resources.
func TestRestAPI(t *testing.T) {
	engine := SetupEngine()

	t.Run("GetResourceList", func(t *testing.T) {
		testGetResourceList(t, engine)
	})

	t.Run("GetResourceByID", func(t *testing.T) {
		testGetResourceByID(t, engine)
	})

	t.Run("GetResourceByIDNotFound", func(t *testing.T) {
		testGetResourceByIDNotFound(t, engine)
	})

	t.Run("GetResourceGroupList", func(t *testing.T) {
		testGetResourceGroupList(t, engine)
	})

	t.Run("GetResourceGroupByID", func(t *testing.T) {
		testGetResourceGroupByID(t, engine)
	})

	t.Run("GetResourceGroupByIDNotFound", func(t *testing.T) {
		testGetResourceGroupByIDNotFound(t, engine)
	})
}

// testGetResourceList tests the retrieval of a list of resources.
// Register a dummy resource to confirm that resources can be retrieved
// For simplicity, register only one resource.
func testGetResourceList(t *testing.T, engine *gin.Engine) {
	// 1. Create the test resource
	req := []map[string]any{
		{
			"deviceID": "TestRestAPI-GetResourceList-device1",
			"type":     "CPU",
		},
	}
	postApiRequest(t, engine, "/cdim/api/v1/devices", http.StatusCreated, req)

	// 2. Search the test resource
	res := getApiRequest(t, engine, "/cdim/api/v1/resources", http.StatusOK)

	// Check the response body
	var response map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	// Verify that count is 1
	assert.Equal(t, float64(1), response["count"], "Expected count to be 1")
	resources, _ := response["resources"].([]any)
	resource := resources[0].(map[string]any)
	returnedDevice, _ := resource["device"].(map[string]any)
	assert.Equal(t, "TestRestAPI-GetResourceList-device1", returnedDevice["deviceID"], "Expected deviceID to match")

	// 3. Delete the test resource
	query := "MATCH (r {deviceID: '%s'}) DETACH DELETE r"
	if err := delete(query, "TestRestAPI-GetResourceList-device1"); err != nil {
		t.Fatalf("failed to delete resource: %v", err)
	}
}

// testGetResourceList tests the retrieval of a list of resources.
// Register a dummy resource to confirm that resources can be retrieved
// For simplicity, register only one resource.
func testGetResourceByID(t *testing.T, engine *gin.Engine) {
	// 1. Create the test resource
	req := []map[string]any{
		{
			"deviceID": "TestRestAPI-GetResource-device1",
			"type":     "CPU",
		},
	}
	postApiRequest(t, engine, "/cdim/api/v1/devices", http.StatusCreated, req)

	// 2. Search the test resource
	res := getApiRequest(t, engine, "/cdim/api/v1/resources/TestRestAPI-GetResource-device1", http.StatusOK)

	var response map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	returnedDevice, _ := response["device"].(map[string]any)
	assert.Equal(t, "TestRestAPI-GetResource-device1", returnedDevice["deviceID"], "Expected deviceID to match")

	// 3. Delete the test resource
	query := "MATCH (r {deviceID: '%s'}) DETACH DELETE r"
	if err := delete(query, "TestRestAPI-GetResource-device1"); err != nil {
		t.Fatalf("failed to delete resource: %v", err)
	}
}

// testGetResourceByIDNotFound tests the retrieval of a resource by ID
// Test for when a resource with the specified ID does not exist
func testGetResourceByIDNotFound(t *testing.T, engine *gin.Engine) {
	// 1. Search the non-existent resource
	res := getApiRequest(t, engine, "/cdim/api/v1/resources/unknown", http.StatusNotFound)

	// Check the response body
	var response map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	assert.Equal(t, "notFound", response["code"], "Expected error code to match")
}

// testGetResourceGroupList tests the retrieval of a list of resource groups.
// It creates a test resource group, retrieves the list, and verifies the response.
// Finally, it cleans up by deleting the created resource group.
func testGetResourceGroupList(t *testing.T, engine *gin.Engine) {
	// 1. Create the test resource group
	req := map[string]any{
		"name":        `TestRestAPI-GetResourceGroupList-group1`,
		"description": `This group contains "double quotes", 'single quotes', \n and <>&.`,
	}
	postApiRequest(t, engine, "/cdim/api/v1/resource-groups", http.StatusCreated, req)

	// 2. Search the test group
	res := getApiRequest(t, engine, "/cdim/api/v1/resource-groups", http.StatusOK)

	// Check the response body
	var response map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	// Verify that count is 2 (1 default group + 1 created group)
	assert.Equal(t, float64(2), response["count"], "Expected count to be 2")
	resourceGroups, _ := response["resourceGroups"].([]any)
	resourceGroup := resourceGroups[1].(map[string]any)
	assert.Equal(t, `This group contains "double quotes", 'single quotes', \n and <>&.`, resourceGroup["description"], "Expected description to match")

	// 3. Delete the test group
	query := "MATCH (vrsg:ResourceGroup {id: '%s'}) DETACH DELETE vrsg"
	if err := delete(query, resourceGroup["id"]); err != nil {
		t.Fatalf("failed to delete ResourceGroup: %v", err)
	}
}

// testGetResourceGroupByID tests the retrieval of a resource group by ID.
// It creates a test resource group, retrieves it by ID, verifies the response,
// and finally cleans up by deleting the created resource group.
func testGetResourceGroupByID(t *testing.T, engine *gin.Engine) {
	// 1. Create the test resource group
	req := map[string]any{
		"name":        `TestRestAPI-GetResourceGroupByID-group1`,
		"description": `This group contains "double quotes", 'single quotes', \n and <>&.`,
	}
	res := postApiRequest(t, engine, "/cdim/api/v1/resource-groups", http.StatusCreated, req)

	// Extract the group ID from the response
	var createResponse map[string]any
	err := json.NewDecoder(res.Body).Decode(&createResponse)
	assert.NoError(t, err, "failed to decode response")
	groupID, ok := createResponse["id"].(string)
	assert.True(t, ok, "Group ID not found in response")

	// 2. Search the test group
	res = getApiRequest(t, engine, fmt.Sprintf("/cdim/api/v1/resource-groups/%s", groupID), http.StatusOK)

	// Check the response body
	var response map[string]any
	err = json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	assert.Equal(t, `This group contains "double quotes", 'single quotes', \n and <>&.`, response["description"], "Expected description to match")

	// 3. Delete the test group
	query := "MATCH (vrsg:ResourceGroup {id: '%s'}) DETACH DELETE vrsg"
	if err := delete(query, groupID); err != nil {
		t.Fatalf("failed to delete ResourceGroup: %v", err)
	}
}

// testGetResourceGroupByIDNotFound tests the retrieval of a resource group by ID
// Test for when a resource group with the specified ID does not exist
func testGetResourceGroupByIDNotFound(t *testing.T, engine *gin.Engine) {
	// 1. Search the non-existent group
	res := getApiRequest(t, engine, "/cdim/api/v1/resource-groups/unknown", http.StatusNotFound)

	// Check the response body
	var response map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error during JSON unmarshal")
	assert.Equal(t, "notFound", response["code"], "Expected error code to match")
}

// postApiRequest is a test helper function that sends a POST request to the specified URL
// with the given body encoded as JSON and verifies the response meets expectations.
//
// Parameters:
//   - t: testing.T instance for test assertions and error reporting
//   - engine: gin.Engine instance to handle the HTTP request
//   - url: target URL path for the POST request
//   - wantCode: expected HTTP status code for response validation
//   - body: request payload that will be JSON-encoded and sent in the request body
//
// The function automatically sets the Content-Type header to "application/json" and
// performs the following validations:
//   - Asserts that the response status code matches wantCode
//   - Asserts that the response Content-Type header starts with "application/json"
//
// Returns:
//   - *httptest.ResponseRecorder: the response recorder containing the HTTP response
//
// The function will call t.Fatalf() if JSON encoding of the request body fails.
func postApiRequest(t *testing.T, engine *gin.Engine, url string, wantCode int, body any) *httptest.ResponseRecorder {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		t.Fatalf("failed to encode request body: %v", err)
	}

	req := httptest.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	// Verify the results
	assert.Equal(t, wantCode, res.Code, fmt.Sprintf("Expected status code %d", wantCode))
	assert.True(t, strings.HasPrefix(res.Header().Get("Content-Type"), "application/json"), "Expected Content-Type to start with application/json")

	return res
}

// getApiRequest performs a GET request to the specified URL using the provided Gin engine
// and validates the response status code and content type.
//
// Parameters:
//   - t: testing instance for assertions
//   - engine: Gin engine to handle the HTTP request
//   - url: target URL path for the GET request
//   - wantCode: expected HTTP status code for validation
//
// Returns:
//   - *httptest.ResponseRecorder: recorder containing the HTTP response
//
// The function automatically asserts that:
//   - Response status code matches wantCode
//   - Content-Type header starts with "application/json"
func getApiRequest(t *testing.T, engine *gin.Engine, url string, wantCode int) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	// Verify the results
	assert.Equal(t, wantCode, res.Code, fmt.Sprintf("Expected status code %d", wantCode))
	assert.True(t, strings.HasPrefix(res.Header().Get("Content-Type"), "application/json"), "Expected Content-Type to start with application/json")

	return res
}

// delete executes a Cypher delete query against the Apache AGE database.
// Directly call the Apache AGE golang driver to delete resources.
// Here, we provide a simple implementation that just executes the query as is.
//
// Parameters:
//   - query: The Cypher delete query string to execute
//   - args: Optional variadic arguments to be used in the query
//
// Returns:
//   - error: Returns an error if transaction begin, query execution, or commit fails
//
// The function automatically handles transaction management and ensures proper
// connection cleanup through deferred disconnection.
func delete(query string, args ...any) error {
	cmdb := database.NewCmDb()
	err := cmdb.CmDbBeginTransaction()
	if err != nil {
		return err
	}
	defer cmdb.CmDbDisconnection()

	_, err = cmdb.CmDbExecCypher(0, query, args...)
	if err != nil {
		return err
	}

	if err := cmdb.CmDbCommit(); err != nil {
		return err
	}

	return nil
}
