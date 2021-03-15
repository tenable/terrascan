package httpserver

import (
	"database/sql"
	"os"
	"time"

	"go.uber.org/zap"
)

// WebhookScanLogger handles the logic to push scan logs to db
type WebhookScanLogger struct {
	test bool
}

type webhookScanLog struct {
	UID                string
	Request            string
	Allowed            bool
	ViolationsSummary  string
	DeniableViolations string
	CreatedAt          time.Time
}

// The file name where the DB is stored. Currently we use an SQLite DB
var dbFileName = "k8s-admission-review-logs.db"

func (g *WebhookScanLogger) log(webhookScanLog webhookScanLog) error {
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
	_, err = statement.Exec(webhookScanLog.UID,
		webhookScanLog.Request,
		webhookScanLog.Allowed,
		webhookScanLog.ViolationsSummary,
		webhookScanLog.DeniableViolations,
		webhookScanLog.CreatedAt)
	if err != nil {
		zap.S().Errorf("failed to insert a new log. error: '%v'", err)
		return err
	}

	return nil
}

func (g *WebhookScanLogger) fetchLogs() ([]webhookScanLog, error) {
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

	var result []webhookScanLog
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

		result = append(result, webhookScanLog{
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

func (g *WebhookScanLogger) fetchLogByID(logUID string) (*webhookScanLog, error) {
	// Fetch a specific log by its request UID

	db, err := g.getDbHandler()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM logs WHERE uid=?", logUID)
	if err != nil {
		zap.S().Errorf("failed query logs table. error: '%v'", err)
		return nil, err
	}
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

		return &webhookScanLog{
			UID:                uid,
			Request:            request,
			Allowed:            allowed,
			ViolationsSummary:  violationsSummary,
			DeniableViolations: deniableViolations,
			CreatedAt:          createdAt,
		}, nil
	}

	return &webhookScanLog{}, nil
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

		db, err := sql.Open("sqlite3", g.dbFilePath())
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
			zap.S().Errorf("failed to create logs table. error: '%v'", err)
			return err
		}
		statement.Exec()
	}

	return nil
}

func (g *WebhookScanLogger) getDbHandler() (*sql.DB, error) {
	g.initDBIfNeeded()

	db, err := sql.Open("sqlite3", g.dbFilePath())
	if err != nil {
		zap.S().Errorf("failed to open sql file. error: '%v'", err)
	}

	return db, err
}

func (g *WebhookScanLogger) dbFilePath() string {
	if g.test {
		return "./" + dbFileName
	}
	// This is where the DB file should be located in the container (It is going to be saved in the host machine volume)
	return "/data/k8s-admission-review-logs.db"
}

// Used for Tests only - clear the DB file after the tests are done
func (g *WebhookScanLogger) clearDbFilePath() {
	os.Remove(g.dbFilePath())
}
