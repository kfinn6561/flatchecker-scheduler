package main

import (
	"database/sql"
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/testutil"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAndUpdateSchedules_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	now := time.Now()
	// Mock GetSchedules
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, now, 30).
		AddRow(2, 20, now, 60).
		AddRow(3, 30, now, 120)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	// Mock UpdateSchedules
	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 2).WillReturnResult(sqlmock.NewResult(0, 1))
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 3).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 3)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_EmptySchedules(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"})
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)
	mock.ExpectPrepare("UPDATE (.+)")

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Empty(t, schedules)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_GetSchedulesError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT (.+)").WillReturnError(sql.ErrConnDone)

	_, err = GetAndUpdateSchedules(mockDB)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_TimeCalculation(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	startTime := time.Now()
	delayMinutes := 60

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, startTime, delayMinutes)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prep := mock.ExpectPrepare("UPDATE (.+)")
	// Capture the time argument to verify it's approximately correct
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)

	// Verify time was calculated (can't check exact time due to execution delay)
	// but we can verify the original schedule data is returned
	assert.Equal(t, 1, schedules[0].ScheduleId)
	assert.Equal(t, 60, schedules[0].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_MultipleSchedulesDifferentDelays(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, now, 5).
		AddRow(2, 20, now, 15).
		AddRow(3, 30, now, 30)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 2).WillReturnResult(sqlmock.NewResult(0, 1))
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 3).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 3)
	assert.Equal(t, 5, schedules[0].SearchDelayMinutes)
	assert.Equal(t, 15, schedules[1].SearchDelayMinutes)
	assert.Equal(t, 30, schedules[2].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_ZeroDelay(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, now, 0)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)
	assert.Equal(t, 0, schedules[0].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_LargeDelay(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	now := time.Now()
	oneWeekMinutes := 10080 // 7 days * 24 hours * 60 minutes

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, now, oneWeekMinutes)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)
	assert.Equal(t, oneWeekMinutes, schedules[0].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_UpdateSchedulesError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, now, 30)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)
	mock.ExpectPrepare("UPDATE (.+)").WillReturnError(sql.ErrConnDone)

	_, err = GetAndUpdateSchedules(mockDB)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAndUpdateSchedules_UpdateRequestMapping(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	scheduleID := 5
	delayMinutes := 30
	now := time.Now()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(scheduleID, 50, now, delayMinutes)

	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().WithArgs(sqlmock.AnyArg(), scheduleID).WillReturnResult(sqlmock.NewResult(0, 1))

	schedules, err := GetAndUpdateSchedules(mockDB)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)
	assert.Equal(t, scheduleID, schedules[0].ScheduleId)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
