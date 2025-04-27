package main

import (
	"database/sql"
	"flatchecker-scheduler/db"
	"time"
)

func GetAndUpdateSchedules(dbConn *sql.DB) ([]db.GetSchedulesResponse, error) {
	schedules, err := db.GetSchedules(dbConn)
	if err != nil {
		return nil, err
	}

	requests := make([]db.UpdateScheduleRequest, len(schedules))
	for i, schedule := range schedules {
		request := db.UpdateScheduleRequest{
			Id:         schedule.ScheduleId,
			NextSearch: schedule.NextSearch.Add(time.Minute * time.Duration(schedule.SearchDelayMinutes)),
		}
		requests[i] = request
	}

	err= db.UpdateSchedules(requests, dbConn)
	if err!=nil{
		return nil, err
	}
	return schedules, nil
}
