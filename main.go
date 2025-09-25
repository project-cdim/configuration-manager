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

package main

import (
	logger "github.com/project-cdim/cdim-go-logger"

	"github.com/project-cdim/configuration-manager/common"
	"github.com/project-cdim/configuration-manager/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	logger_common "github.com/project-cdim/cdim-go-logger/common"
)

// v1 route base url
const urlBaseV1 = "/" + common.ProjectName + "/api/v1"

// Audit Trail Logger
var log, _ = logger.New(logger_common.Option{Tag: logger_common.TAG_TRAIL})

func main() {
	engine := SetupEngine()
	engine.Run(":8080")
}

// SetupEngine initializes and returns a new instance of the gin Engine. This function configures
// the engine with essential middleware, including a custom logging middleware for audit trails,
// and CORS support using the default configuration. It also sets up a versioned API route group
// (v1) and defines routes for various operations such as retrieving resource lists, individual
// resources, nodes, CXL switches, and racks from a configuration management database. Additionally,
// it includes routes for searching resources based on conditions, registering devices, and updating
// resource annotations.
//
// Returns:
// - A pointer to the configured gin Engine instance, ready to handle incoming HTTP requests.
func SetupEngine() *gin.Engine {
	// Create an instance of the gin Engine
	engine := gin.Default()

	// Add custom middleware to output audit logs to the gin Engine
	engine.Use(logMiddleware())

	engine.Use(cors.New(cors.Config{
		// Allowed Methods
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		// Allowed origins
		AllowOrigins: []string{
			"*",
		},
		// Allowed HTTP request headers
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}))

	// v1 route group
	v1 := engine.Group(urlBaseV1)
	{
		// Retrieve a list of all resources from the configuration management database
		v1.GET("/resources", controller.GetResourceList)

		// Retrieve a specific resource from the configuration management database
		v1.GET("/resources/:id", controller.GetResource)

		// Retrieve a list of resources available for configuration design from the configuration management database
		v1.GET("/resources/available", controller.GetAvailableResourceList)

		// Retrieve a list of resources unused for configuration design from the configuration management database
		v1.GET("/resources/unused", controller.GetUnusedResourceList)

		// Update the additional information of a specific resource
		v1.PUT("/resources/:id/annotation", controller.UpdateAnnotation)

		// Retrieve a list of all resource groups from the configuration management database
		v1.GET("/resource-groups", controller.GetGroupList)

		// Register a new resource group in the configuration management database
		v1.POST("/resource-groups", controller.CreateGroup)

		// Retrieve a specific resource group from the configuration management database
		v1.GET("/resource-groups/:id", controller.GetGroup)

		// Update the information of a specific resource group
		v1.PUT("/resource-groups/:id", controller.UpdateGroup)

		// Delete a specific resource group from the configuration management database
		v1.DELETE("/resource-groups/:id", controller.DeleteGroup)

		// Update the resource group to which the resource belongs
		v1.PUT("/resources/:id/resource-groups", controller.AssignResourceToGroup)

		// Retrieve a list of all nodes and their associated resources from the configuration management database
		v1.GET("/nodes", controller.GetNodeList)

		// Retrieve a specific node and its associated resources from the configuration management database
		v1.GET("/nodes/:id", controller.GetNode)

		// Retrieve a list of all CXL switches from the configuration management database
		v1.GET("/cxlswitches", controller.GetCxlSwitchList)

		// Retrieve a specific CXL switch from the configuration management database
		v1.GET("/cxlswitches/:id", controller.GetCxlSwitch)

		// Retrieve a specific rack from the configuration management database
		v1.GET("/racks/:id", controller.GetRack)

		// Register multiple device information in the configuration management database
		v1.POST("/devices", controller.RegisterDevice)
	}

	return engine
}

// logMiddleware creates a custom middleware for use with the Gin web framework. This middleware
// is designed to log the start and end of API requests for monitoring and debugging purposes.
// The logging is performed in two stages:
//   - Before the execution of the API handler, it logs the start of the request, including the
//     HTTP method and the request URL path.
//   - After the API handler has executed, it logs the end of the request along with the HTTP
//     status code of the response.
//
// This middleware is useful for tracking API usage patterns and identifying potential issues
// with request handling.
//
// Returns:
// - A Gin middleware function (gin.HandlerFunc) that can be used to log request start and end times.
func logMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.TrailReq(ctx.Request.Method, ctx.Request.URL.Path, "-", "request start.")
		ctx.Next()
		log.TrailRes(ctx.Writer.Status(), "response end.")
	}
}
