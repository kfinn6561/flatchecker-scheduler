package db

import (
	"database/sql"
	"fmt"
	"time"
)

const UPDATE_SCHEDULE_SQL_FILE_NAME = "update_schedule.sql"

type UpdateScheduleRequest struct {
	Id         int
	NextSearch time.Time
}

func UpdateSchedules(requests []UpdateScheduleRequest, dbConn *sql.DB) error {
	stmtString, err := readSqlFile(UPDATE_SCHEDULE_SQL_FILE_NAME)
	if err != nil {
		return fmt.Errorf("error reading sql file: %v", err)
	}

	stmt, err := dbConn.Prepare(stmtString)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}

	for _, request := range requests {
		_, err := stmt.Exec(request.NextSearch, request.Id)
		if err != nil {
			return fmt.Errorf("error executing statement: %v", err)
		}
	}

	return nil
}
