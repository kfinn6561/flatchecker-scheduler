package db

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSchedules_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Set up expected rows
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, time.Now(), 30).
		AddRow(2, 20, time.Now().Add(time.Hour), 60).
		AddRow(3, 30, time.Now().Add(2*time.Hour), 120)

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Len(t, schedules, 3)
	assert.Equal(t, 1, schedules[0].ScheduleId)
	assert.Equal(t, 10, schedules[0].SearchId)
	assert.Equal(t, 30, schedules[0].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"})
	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Empty(t, schedules)
	assert.NotNil(t, schedules)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("database connection lost")
	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnError(expectedErr)

	_, err = GetSchedules(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error querying database")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create rows with incompatible types
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow("not-an-int", 10, time.Now(), 30)

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	_, err = GetSchedules(db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error scanning row")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, time.Now(), 30).
		RowError(0, errors.New("row iteration error"))

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	_, err = GetSchedules(db)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_TimeParsingPastDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	pastTime := time.Now().Add(-24 * time.Hour)
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, pastTime, 30)

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)
	assert.True(t, schedules[0].NextSearch.Before(time.Now()))

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_TimeParsingFutureDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	futureTime := time.Now().Add(48 * time.Hour)
	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, futureTime, 30)

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Len(t, schedules, 1)
	assert.True(t, schedules[0].NextSearch.After(time.Now()))

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_LargeResultSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"})
	for i := 1; i <= 1000; i++ {
		rows.AddRow(i, i*100, time.Now(), 30)
	}

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Len(t, schedules, 1000)
	assert.Equal(t, 1, schedules[0].ScheduleId)
	assert.Equal(t, 1000, schedules[999].ScheduleId)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_VerifyQueryContent(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"})
	// Expect query to match pattern
	mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnRows(rows)

	_, err = GetSchedules(db)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedules_ConnectionClosed(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	// Close connection before query
	db.Close()

	_, err = GetSchedules(db)
	assert.Error(t, err)
}

func TestGetSchedules_VariousDelayValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ScheduleId", "SearchId", "NextSearch", "SearchDelayMinutes"}).
		AddRow(1, 10, time.Now(), 5).
		AddRow(2, 20, time.Now(), 0).
		AddRow(3, 30, time.Now(), 1440). // 1 day
		AddRow(4, 40, time.Now(), 10080) // 1 week

	mock.ExpectQuery("(?i)SELECT (.+)").WillReturnRows(rows)

	schedules, err := GetSchedules(db)
	require.NoError(t, err)
	assert.Len(t, schedules, 4)
	assert.Equal(t, 5, schedules[0].SearchDelayMinutes)
	assert.Equal(t, 0, schedules[1].SearchDelayMinutes)
	assert.Equal(t, 1440, schedules[2].SearchDelayMinutes)
	assert.Equal(t, 10080, schedules[3].SearchDelayMinutes)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetSchedulesResponse_StructFields(t *testing.T) {
	// Verify the struct has the expected fields
	schedule := GetSchedulesResponse{
		ScheduleId:         123,
		SearchId:           456,
		NextSearch:         time.Now(),
		SearchDelayMinutes: 30,
	}

	assert.Equal(t, 123, schedule.ScheduleId)
	assert.Equal(t, 456, schedule.SearchId)
	assert.NotZero(t, schedule.NextSearch)
	assert.Equal(t, 30, schedule.SearchDelayMinutes)
}
