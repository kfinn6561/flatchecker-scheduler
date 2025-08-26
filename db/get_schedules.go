package db

import (
	"database/sql"
	"fmt"
	"time"
)

const GET_SCHEDULE_SQL_NAME = "get_schedules.sql"

type GetSchedulesResponse struct {
	ScheduleId         int
	SearchId           int
	NextSearch         time.Time
	SearchDelayMinutes int
}

func GetSchedules(db *sql.DB) ([]GetSchedulesResponse, error) {
	stmtString, err := readSqlFile(GET_SCHEDULE_SQL_NAME)
	if err != nil {
		return nil, fmt.Errorf("error reading sql file: %v", err)
	}

	rows, err := db.Query(stmtString)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var schedules []GetSchedulesResponse

	for rows.Next() {
		var schedule GetSchedulesResponse
		if err := rows.Scan(&schedule.ScheduleId, &schedule.SearchId, &schedule.NextSearch, &schedule.SearchDelayMinutes); err != nil {
			return schedules, fmt.Errorf("error scanning row: %v", err)
		}
		schedules = append(schedules, schedule)
	}
	if err = rows.Err(); err != nil {
		return schedules, fmt.Errorf("error iterating rows: %v", err)
	}
	fmt.Println(schedules)

	return schedules, nil
}
