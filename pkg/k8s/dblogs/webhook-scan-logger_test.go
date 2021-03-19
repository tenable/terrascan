/*
    Copyright (C) 2020 Accurics, Inc.

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
	"testing"
	"time"
)

func TestLogs(t *testing.T) {

	var logger = WebhookScanLogger{
		Test: true,
	}
	defer logger.ClearDbFilePath()

	// insert a new db record
	var log = WebhookScanLog{
		UID:                "myUID",
		Request:            "MyRequest",
		CreatedAt:          time.Now(),
		Allowed:            true,
		DeniableViolations: "MyViolations",
		ViolationsSummary:  "ViolationsSummary",
	}

	t.Run("initialize db", func(t *testing.T) {

		// no logs exist in db, should return 0 logs
		fetchedLogs, err := logger.FetchLogs()
		if len(fetchedLogs) > 0 {
			t.Errorf("no logs should exist in db; got: '%v' logs", len(fetchedLogs))
		}
		if err != nil {
			t.Errorf("unexpected error: '%v'", err)
		}
	})

	t.Run("insert db record", func(t *testing.T) {
		if err := logger.Log(log); err != nil {
			t.Errorf("unexpected error: '%v'", err)
		}
	})

	// test case 3: fetch db record by log ID
	myFetchLog, err := logger.FetchLogByID(log.UID)
	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}

	t.Run("one db record", func(t *testing.T) {
		fetchedLogs, err := logger.FetchLogs()
		if err != nil {
			t.Errorf("unexpected error: '%v'", err)
		}
		if len(fetchedLogs) != 1 {
			t.Errorf("db has one log, got: '%v' logs", len(fetchedLogs))
		}
	})

	t.Run("fetch record by id", func(t *testing.T) {

		if len(myFetchLog.UID) < 1 {
			t.Errorf("Log with ID: '%v' is not returned by fetchLogByID", log.UID)
		}

		if myFetchLog.UID != log.UID {
			t.Errorf("Wrong UID. Expected '%v', Got: '%v'", log.UID, myFetchLog.UID)

		}
	})

	t.Run("verify allowed", func(t *testing.T) {
		if myFetchLog.Allowed != log.Allowed {
			t.Errorf("Wrong Allowed. Expected '%v', Got: '%v'", log.Allowed, myFetchLog.Allowed)
		}
	})

	t.Run("verify violations summary", func(t *testing.T) {
		if myFetchLog.ViolationsSummary != log.ViolationsSummary {
			t.Errorf("Wrong ViolationsSummary. Expected '%v', Got: '%v'", log.ViolationsSummary, myFetchLog.ViolationsSummary)
		}
	})

	t.Run("verify request", func(t *testing.T) {
		if myFetchLog.Request != log.Request {
			t.Errorf("Wrong Request. Expected '%v', Got: '%v'", log.Request, myFetchLog.Request)
		}
	})

	t.Run("verify deniable violations", func(t *testing.T) {
		if myFetchLog.DeniableViolations != log.DeniableViolations {
			t.Errorf("Wrong DeniableViolations. Expected '%v', Got: '%v'", log.DeniableViolations, myFetchLog.DeniableViolations)
		}
	})

	t.Run("verify timestamp", func(t *testing.T) {
		if myFetchLog.CreatedAt.Unix() != log.CreatedAt.Unix() {
			t.Errorf("Wrong CreatedAt. Expected '%v', Got: '%v'", log.CreatedAt, myFetchLog.CreatedAt)
		}
	})
}
