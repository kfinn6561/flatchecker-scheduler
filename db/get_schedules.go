package db

import "database/sql"

const GET_SCHEDULE_SQL_NAME = "get_schedules.sql"

type GetSchedulesResponse struct {
	ScheduleId int
	SearchId   int
}

func GetSchedules(db *sql.DB) ([]GetSchedulesResponse, error) {
	stmtString, err := readSqlFile(GET_SCHEDULE_SQL_NAME)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(stmtString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []GetSchedulesResponse

	for rows.Next() {
		var schedule GetSchedulesResponse
		if err := rows.Scan(&schedule.ScheduleId, &schedule.SearchId); err != nil {
			return schedules, err
		}
		schedules = append(schedules, schedule)
	}
	if err = rows.Err(); err != nil {
		return schedules, err
	}

	return schedules, nil
}
