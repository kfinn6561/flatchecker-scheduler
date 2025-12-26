package db

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test update request
func createTestUpdateRequest(id int, nextSearch time.Time) UpdateScheduleRequest {
	return UpdateScheduleRequest{
		Id:         id,
		NextSearch: nextSearch,
	}
}

func TestUpdateSchedules_SingleUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	nextSearch := time.Now().Add(30 * time.Minute)
	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(1, nextSearch),
	}

	mock.ExpectPrepare("UPDATE (.+)").
		ExpectExec().
		WithArgs(nextSearch, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_MultipleUpdates_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	baseTime := time.Now()
	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(1, baseTime.Add(5*time.Minute)),
		createTestUpdateRequest(2, baseTime.Add(15*time.Minute)),
		createTestUpdateRequest(3, baseTime.Add(30*time.Minute)),
		createTestUpdateRequest(4, baseTime.Add(60*time.Minute)),
		createTestUpdateRequest(5, baseTime.Add(120*time.Minute)),
	}

	prep := mock.ExpectPrepare("UPDATE (.+)")
	for _, req := range requests {
		prep.ExpectExec().
			WithArgs(req.NextSearch, req.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_EmptyRequests(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	var requests []UpdateScheduleRequest

	mock.ExpectPrepare("UPDATE (.+)")

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(1, time.Now()),
	}

	expectedErr := errors.New("failed to prepare statement")
	mock.ExpectPrepare("UPDATE (.+)").WillReturnError(expectedErr)

	err = UpdateSchedules(requests, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing statement")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_ExecuteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(1, time.Now()),
		createTestUpdateRequest(2, time.Now()),
	}

	prep := mock.ExpectPrepare("UPDATE (.+)")
	prep.ExpectExec().
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	prep.ExpectExec().
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnError(errors.New("execute failed"))

	err = UpdateSchedules(requests, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error executing statement")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_VerifyParameters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	specificTime := time.Date(2024, 12, 20, 15, 30, 0, 0, time.UTC)
	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(10, specificTime),
	}

	mock.ExpectPrepare("UPDATE (.+)").
		ExpectExec().
		WithArgs(specificTime, 10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_TimeFormatting(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	testCases := []time.Time{
		time.Now(),
		time.Now().Add(24 * time.Hour),
		time.Now().Add(-24 * time.Hour),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	prep := mock.ExpectPrepare("UPDATE (.+)")
	for i, testTime := range testCases {
		prep.ExpectExec().
			WithArgs(testTime, i+1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	var requests []UpdateScheduleRequest
	for i, testTime := range testCases {
		requests = append(requests, createTestUpdateRequest(i+1, testTime))
	}

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_StatementReuse(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(1, time.Now()),
		createTestUpdateRequest(2, time.Now()),
		createTestUpdateRequest(3, time.Now()),
	}

	// Expect single Prepare call (statement reuse)
	prep := mock.ExpectPrepare("UPDATE (.+)")
	for _, req := range requests {
		prep.ExpectExec().
			WithArgs(sqlmock.AnyArg(), req.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_LargeBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create 1000 update requests
	requests := make([]UpdateScheduleRequest, 1000)
	baseTime := time.Now()
	for i := 0; i < 1000; i++ {
		requests[i] = createTestUpdateRequest(i+1, baseTime.Add(time.Duration(i)*time.Minute))
	}

	prep := mock.ExpectPrepare("UPDATE (.+)")
	for i := 0; i < 1000; i++ {
		prep.ExpectExec().
			WithArgs(sqlmock.AnyArg(), i+1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSchedules_ParameterOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	nextSearch := time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC)
	scheduleID := 42

	requests := []UpdateScheduleRequest{
		createTestUpdateRequest(scheduleID, nextSearch),
	}

	// Verify parameters are passed in correct order: NextSearch, Id
	mock.ExpectPrepare("UPDATE (.+)").
		ExpectExec().
		WithArgs(nextSearch, scheduleID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = UpdateSchedules(requests, db)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateScheduleRequest_StructFields(t *testing.T) {
	// Verify the struct has the expected fields
	testTime := time.Now()
	request := UpdateScheduleRequest{
		Id:         123,
		NextSearch: testTime,
	}

	assert.Equal(t, 123, request.Id)
	assert.Equal(t, testTime, request.NextSearch)
}
