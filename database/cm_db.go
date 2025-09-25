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

package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/project-cdim/configuration-manager/common"

	"github.com/apache/age/drivers/golang/age"
	_ "github.com/lib/pq"
)

// Database connection information
const (
	dsnTemplate string = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"                            // DSN connection info to Database
	GRAPH_NAME  string = "cdim_graph"                                                                               // Graph DB name
	SECRETS_URL string = "http://localhost:3500/v1.0/secrets/" + common.ProjectName + "-configuration-manager/cmdb" // secret URL
)

// SecretCmdb is cmdb's secret information.
type SecretCmdb struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// Structure with DB handler and transaction as fields
// - DB handler: nil when not connected, non-nil when connected
// - Transaction: nil when no transaction is active, non-nil when a transaction is active
type CmDb struct {
	Db *sql.DB // DB handler. nil when not connected.
	Tx *sql.Tx // Transaction. nil when no transaction is active.
}

// NewCmDb creates a new instance of CmDb with both Db and Tx fields initialized as nil.
// This function is useful for setting up a new database connection handler with no active connection or transaction.
func NewCmDb() CmDb {
	return CmDb{
		Db: nil,
		Tx: nil,
	}
}

// CmDbConnection connects to the graph database using the provided secret information.
// It attempts to establish a connection to a PostgreSQL database and then loads the graph DB extension.
// If any step fails, it returns an error with a specific error code and message.
//
// - Retrieves secret information necessary for the connection.
// - Formats the DSN (Data Source Name) string for PostgreSQL connection.
// - Attempts to open a connection to the PostgreSQL database.
// - Loads the graph DB extension.
// - Returns nil on success or an error if any step fails.
func (g *CmDb) CmDbConnection() error {
	var err error

	// Retrieve secret information
	secretCmdb, err := GetSecretCmdb()
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf(dsnTemplate, secretCmdb.Host, secretCmdb.Port, secretCmdb.User, secretCmdb.Password, secretCmdb.Dbname)

	// Connect to PostgreSQL and load the graph DB extension
	if g.Db, err = sql.Open("postgres", dsn); err != nil {
		common.Log.Error(err.Error())
		return err
	}

	if _, err = age.GetReady(g.Db, GRAPH_NAME); err != nil {
		common.Log.Error(err.Error())
		if err = g.Db.Close(); err != nil {
			common.Log.Error(err.Error())
		}
		return err
	}

	return nil
}

// CmDbBeginTransaction initiates a new database transaction.
// It checks if the database connection is already established; if not, it connects to the Graph DB first.
// After ensuring the database connection, it attempts to start a new transaction.
// If any step fails, it returns an error with a specific error code and message.
//
// - Checks if the database connection exists, connects if necessary.
// - Attempts to begin a new transaction.
// - Returns nil on success or an error if the connection or transaction creation fails.
//
// CmDbBeginTransaction create a transaction.
// If not connected to Graph DB, create transaction after connecting to Graph DB.
func (g *CmDb) CmDbBeginTransaction() error {
	var err error

	if g.Db == nil {
		if err = g.CmDbConnection(); err != nil {
			return err
		}
	}

	if g.Tx, err = g.Db.Begin(); err != nil {
		common.Log.Error(err.Error())
		if err = g.Db.Close(); err != nil {
			common.Log.Error(err.Error())
		}
		return err
	}

	return nil
}

// CmDbCommit commits the current transaction.
// It finalizes the transaction by committing any changes made during the transaction to the database.
// Before committing, it checks if a transaction is active; if not, it returns an error indicating the transaction is invalid.
// If the commit operation fails, it returns an error with a specific error code and message.
// After attempting to commit, it clears the current transaction regardless of success or failure.
//
// - Checks if a transaction is currently active.
// - Attempts to commit the current transaction.
// - Clears the current transaction.
// - Returns nil on success or an error if the transaction is invalid or the commit fails.
func (g *CmDb) CmDbCommit() error {
	var err error
	defer g.clearTransaction()

	if g.Tx == nil {
		common.Log.Error("transaction is invalid.")
		return errors.New("transaction is invalid")
	}

	if err = g.Tx.Commit(); err != nil {
		common.Log.Error(err.Error())
		return err
	}

	return nil
}

// CmDbRollback rolls back the current transaction.
// It undoes any changes made during the transaction, ensuring the database state remains consistent.
// Before rolling back, it checks if a transaction is active; if not, it returns an error indicating the transaction is invalid.
// If the rollback operation fails, it returns an error with a specific error code and message.
// After attempting to roll back, it clears the current transaction regardless of success or failure.
//
// - Checks if a transaction is currently active.
// - Attempts to roll back the current transaction.
// - Clears the current transaction.
// - Returns nil on success or an error if the transaction is invalid or the rollback fails.
func (g *CmDb) CmDbRollback() error {
	var err error
	defer g.clearTransaction()

	if g.Tx == nil {
		common.Log.Error("transaction is invalid.")
		return errors.New("transaction is invalid")
	}

	if err = g.Tx.Rollback(); err != nil {
		common.Log.Error(err.Error())
		return err
	}

	return nil
}

// CmDbDisconnection safely disconnects from the Graph DB.
// If there is an active transaction, it first attempts to roll it back to ensure data consistency before disconnecting.
// This method ensures that the database connection is cleanly closed, preventing potential data leaks or inconsistencies.
// It returns an error if either the rollback or disconnection fails, with a specific error code and message for troubleshooting.
//
// - Checks if there is an active transaction and attempts to roll it back.
// - Closes the database connection.
// - Returns nil on successful disconnection or an error if any step fails.
func (g *CmDb) CmDbDisconnection() error {
	var err error
	defer g.clearConnection()

	if g.Tx != nil {
		g.CmDbRollback()
	}

	if err = g.Db.Close(); err != nil {
		common.Log.Error(err.Error())
		return err
	}

	return nil
}

// CmDbExecCypher executes a Cypher query and returns the resulting Cursor.
// This function is designed to execute a specified Cypher query within an active transaction context.
// It requires the number of columns expected in the result and the Cypher query string as parameters.
// If there is no active transaction, it returns an error indicating the transaction is invalid.
// On successful execution, it returns a CypherCursor which can be used to iterate over the results.
//
// - Requires an active transaction to execute the Cypher query.
// - Returns a CypherCursor on success or an error if the execution fails or if there is no active transaction.
func (g *CmDb) CmDbExecCypher(columnCount int, cypher string, args ...any) (*age.CypherCursor, error) {
	if g.Tx == nil {
		common.Log.Error("Cypher was not executed due to invalid transaction.")
		return nil, errors.New("transaction is invalid")
	}

	cypherCursor, err := age.ExecCypher(g.Tx, GRAPH_NAME, columnCount, cypher, args...)
	if err != nil {
		common.Log.Error(err.Error())
		return nil, err
	}

	return cypherCursor, nil
}

// clearTransaction sets the transaction field of the structure to nil.
// This method is used internally to ensure that the transaction field is properly reset after a transaction is committed or rolled back.
// It checks if the transaction field is not nil before setting it to nil to avoid unnecessary operations.
//
// - Ensures the transaction field is reset to nil after use.
func (g *CmDb) clearTransaction() {
	if g.Tx != nil {
		g.Tx = nil
	}
}

// clearConnection resets the DB handler field to nil within the structure.
// If the transaction field is not nil, it also resets the transaction field to nil.
// This method is used internally to ensure that both the DB handler and transaction fields are properly reset after disconnection or when cleaning up resources.
// It ensures that the structure is in a clean state, preventing potential issues related to stale connections or transactions.
//
// - Resets the DB handler field to nil.
// - Also resets the transaction field to nil if it is not already nil.
func (g *CmDb) clearConnection() {
	g.clearTransaction()
	if g.Db != nil {
		g.Db = nil
	}
}

// getSecretCmdbImpl retrieves all items under the secret name "cmdb" from the Dapr secret store "cmdb".
// This function constructs the request URL for the Dapr secret store, sends an HTTP GET request to that URL, and processes the response.
// If the request is successful, it reads the response body, unmarshals the JSON into a struct of type secretCmdb, and returns this struct.
// In case of any errors (e.g., network issues, errors during the unmarshalling process), it returns an empty secretCmdb struct along with the error.
//
// The function performs the following steps:
// 1. Constructs the request URL for the Dapr secret store.
// 2. Sends an HTTP GET request to the constructed URL.
// 3. Reads the response body upon a successful request.
// 4. Unmarshals the JSON response body into a secretCmdb struct.
// 5. Returns the unmarshalled secretCmdb struct or an error if any step fails.
func getSecretCmdbImpl() (SecretCmdb, error) {
	resp, err := http.Get(SECRETS_URL)
	if err != nil {
		common.Log.Error(err.Error())
		return SecretCmdb{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Log.Error(err.Error())
		return SecretCmdb{}, err
	}
	var cmdb SecretCmdb
	err = json.Unmarshal(body, &cmdb)
	if err != nil {
		common.Log.Error(err.Error())
		return SecretCmdb{}, err
	}

	return cmdb, nil
}

var GetSecretCmdb = getSecretCmdbImpl
