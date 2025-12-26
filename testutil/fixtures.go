package testutil

import (
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/pubsublib"
	"time"
)

// CreateTestSchedule creates a test GetSchedulesResponse for testing
func CreateTestSchedule(scheduleID, searchID, delayMinutes int) db.GetSchedulesResponse {
	return db.GetSchedulesResponse{
		ScheduleId:         scheduleID,
		SearchId:           searchID,
		NextSearch:         time.Now(),
		SearchDelayMinutes: delayMinutes,
	}
}

// CreateTestScheduleWithTime creates a test GetSchedulesResponse with a specific time
func CreateTestScheduleWithTime(scheduleID, searchID, delayMinutes int, nextSearch time.Time) db.GetSchedulesResponse {
	return db.GetSchedulesResponse{
		ScheduleId:         scheduleID,
		SearchId:           searchID,
		NextSearch:         nextSearch,
		SearchDelayMinutes: delayMinutes,
	}
}

// CreateTestUpdateRequest creates a test UpdateScheduleRequest
func CreateTestUpdateRequest(id int, nextSearch time.Time) db.UpdateScheduleRequest {
	return db.UpdateScheduleRequest{
		Id:         id,
		NextSearch: nextSearch,
	}
}

// CreateTestPubSubMessage creates a test ScheduledSearchesMessage
func CreateTestPubSubMessage(scheduleID, searchID int) pubsublib.ScheduledSearchesMessage {
	return pubsublib.ScheduledSearchesMessage{
		ScheduleId: scheduleID,
		SearchId:   searchID,
	}
}

// CreateTestConfig creates a test configuration map
func CreateTestConfig(environment string) map[string]string {
	config := make(map[string]string)
	config["environment"] = environment
	config["user-name"] = "testuser"
	config["db-name"] = "testdb"

	if environment == "dev" {
		config["ip-address"] = "127.0.0.1"
	} else if environment == "prod" {
		config["connection-name"] = "project:region:instance"
	}

	return config
}
