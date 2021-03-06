/*
 *   Copyright 2015 Benoit LETONDOR
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var database *sql.DB

const (
	CONN_STRING = `host=/var/run/postgresql sslmode=disable user=postgres`
)

func Init() (*sql.DB, error) {
	// Init Mysql DB
	dbLink, err := sql.Open("postgres", CONN_STRING)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = dbLink.Ping()
	if err != nil {
		return nil, err
	}

	// Set up global var
	database = dbLink

	return database, nil
}
