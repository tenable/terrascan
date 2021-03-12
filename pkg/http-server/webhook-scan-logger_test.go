package httpserver

import (
	"testing"
	"time"
)

func TestLogs(t *testing.T) {
	t.Run("Log a new webhook scan", func(t *testing.T) {
		var logger = WebhookScanLogger{
			test: true,
		}
		defer logger.clearDbFilePath()

		fetchedLogs, err := logger.fetchLogs()
		if len(fetchedLogs) > 0 {
			t.Errorf("At the beginning no logs should exist. Got: '%v'", len(fetchedLogs))
		}

		if err != nil {
			t.Errorf("Got error")
		}

		var log = webhookScanLog{
			UID:                "myUID",
			Request:            "MyRequest",
			CreatedAt:          time.Now(),
			Allowed:            true,
			DeniableViolations: "MyViolations",
			ViolationsSummary:  "ViolationsSummary",
		}

		logger.log(log)

		fetchedLogs, err = logger.fetchLogs()
		if err != nil {
			t.Errorf("Got error")
		}

		if len(fetchedLogs) != 1 {
			t.Errorf("A new log should be returned. Got: '%v' logs", len(fetchedLogs))
		}

		myFetchLog, err := logger.fetchLogById(log.UID)
		if err != nil {
			t.Errorf("Got error")
		}

		if len(myFetchLog.UID) < 1 {
			t.Errorf("Log with ID: '%v' is not returned by fetchLogById", log.UID)
		}

		if myFetchLog.UID != log.UID {
			t.Errorf("Wrong UID. Expected '%v', Got: '%v'", log.UID, myFetchLog.UID)

		}

		if myFetchLog.Allowed != log.Allowed {
			t.Errorf("Wrong Allowed. Expected '%v', Got: '%v'", log.Allowed, myFetchLog.Allowed)
		}

		if myFetchLog.ViolationsSummary != log.ViolationsSummary {
			t.Errorf("Wrong ViolationsSummary. Expected '%v', Got: '%v'", log.ViolationsSummary, myFetchLog.ViolationsSummary)
		}

		if myFetchLog.Request != log.Request {
			t.Errorf("Wrong Request. Expected '%v', Got: '%v'", log.Request, myFetchLog.Request)
		}

		if myFetchLog.DeniableViolations != log.DeniableViolations {
			t.Errorf("Wrong DeniableViolations. Expected '%v', Got: '%v'", log.DeniableViolations, myFetchLog.DeniableViolations)
		}

		if myFetchLog.CreatedAt.Unix() != log.CreatedAt.Unix() {
			t.Errorf("Wrong CreatedAt. Expected '%v', Got: '%v'", log.CreatedAt, myFetchLog.CreatedAt)
		}
	})
}
