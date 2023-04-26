/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package dblogs

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	// importing sqlite driver
	_ "modernc.org/sqlite"
)

// WebhookScanLogger handles the logic to push scan logs to db
type WebhookScanLogger struct {
	Test bool
}

// WebhookScanLog database model for log records
type WebhookScanLog struct {
	UID                string
	Request            string
	Allowed            bool
	ViolationsSummary  string
	DeniableViolations string
	CreatedAt          time.Time
}

// The file name where the DB is stored. Currently we use an SQLite DB
var dbFileName = "k8s-admission-review-logs.db"

// NewWebhookScanLogger returns a new WebhookScanLogger struct
func NewWebhookScanLogger() *WebhookScanLogger {
	return &WebhookScanLogger{}
}

// Log creates a new db record for the admission request
func (g *WebhookScanLogger) Log(WebhookScanLog WebhookScanLog) error {
	// Insert a new Log record to the DB

	db, err := g.getDbHandler()
	if err != nil {
		return err
	}
	defer db.Close()

	insertLogSQL := `INSERT INTO logs(uid, request, allowed, violations_summary, deniable_violations, created_at)
					 VALUES (?, ?, ?, ?, ?, ?)`

	statement, err := db.Prepare(insertLogSQL)
	if err != nil {
		zap.S().Errorf("failed preparing SQL statement. error: '%v'", err)
		return err
	}
	_, err = statement.Exec(WebhookScanLog.UID,
		WebhookScanLog.Request,
		WebhookScanLog.Allowed,
		WebhookScanLog.ViolationsSummary,
		WebhookScanLog.DeniableViolations,
		WebhookScanLog.CreatedAt)
	if err != nil {
		zap.S().Errorf("failed to insert a new log. error: '%v'", err)
		return err
	}

	return nil
}

// FetchLogs retrieves all the logs from the database
func (g *WebhookScanLogger) FetchLogs() ([]WebhookScanLog, error) {
	// Fetch the entire logs in the DB, ordered by created_at DESC (the most updated will be at the top)

	db, err := g.getDbHandler()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM logs ORDER BY created_at DESC")
	if err != nil {
		zap.S().Errorf("failed query logs table. error: '%v'", err)
		return nil, err
	}

	var result []WebhookScanLog
	defer row.Close()
	for row.Next() {
		var id int
		var uid string
		var request string
		var allowed bool
		var violationsSummary string
		var deniableViolations string
		var createdAt time.Time
		row.Scan(&id, &uid, &request, &allowed, &violationsSummary, &deniableViolations, &createdAt)

		result = append(result, WebhookScanLog{
			UID:                uid,
			Request:            request,
			Allowed:            allowed,
			ViolationsSummary:  violationsSummary,
			DeniableViolations: deniableViolations,
			CreatedAt:          createdAt,
		})
	}

	return result, nil
}

// GetLogURL returns a url to the UI page for reviewing the validating admission request log
func (g *WebhookScanLogger) GetLogURL(host, logUID string) string {
	// Use this as the link to show the a specific log
	return fmt.Sprintf("https://%v/k8s/webhooks/logs/%v", host, logUID)
}

func (g *WebhookScanLogger) initDBIfNeeded() error {

	// Check where the SQL file exists. If it does do nothing. Otherwise, create the DB file and the Logs table.
	if _, err := os.Stat(g.dbFilePath()); os.IsNotExist(err) {
		file, err := os.Create(g.dbFilePath())
		if err != nil {
			zap.S().Errorf("failed create db file. error: '%v'", err)
			return err
		}
		file.Close()

		db, err := sql.Open("sqlite", g.dbFilePath())
		if err != nil {
			zap.S().Errorf("failed to open sql file. error: '%v'", err)
			return err
		}
		defer db.Close()

		createLogsTableSQL := `CREATE TABLE logs (
													"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
													"uid" TEXT UNIQUE,
													"request" TEXT,
													"allowed" INTEGER,
													"violations_summary" TEXT,
													"deniable_violations" TEXT,
													"created_at" DATETIME
												  );`
		statement, err := db.Prepare(createLogsTableSQL)
		if err != nil {
			zap.S().Errorf("failed to prepare sql query to create logs table. error: '%v'", err)
			return err
		}

		if _, err := statement.Exec(); err != nil {
			zap.S().Errorf("failed to create logs table, error: '%v'", err)
			return err
		}
	}

	return nil
}

func (g *WebhookScanLogger) getDbHandler() (*sql.DB, error) {
	g.initDBIfNeeded()

	db, err := sql.Open("sqlite", g.dbFilePath())
	if err != nil {
		zap.S().Errorf("failed to open sql file. error: '%v'", err)
	}

	return db, err
}

func (g *WebhookScanLogger) dbFilePath() string {
	if g.Test {
		return "./" + dbFileName
	}
	// This is where the DB file should be located in the container (It is going to be saved in the host machine volume)
	return "/data/k8s-admission-review-logs.db"
}

// ClearDbFilePath used for Tests only - clear the DB file after the tests are done
func (g *WebhookScanLogger) ClearDbFilePath() {
	os.Remove(g.dbFilePath())
}
