package mapper

import (
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapSchedules_SingleSchedule(t *testing.T) {
	// Test mapping a single schedule
	dbSchedules := []db.GetSchedulesResponse{
		testutil.CreateTestSchedule(100, 200, 30),
	}

	result := MapSchedules(dbSchedules)

	assert.Len(t, result, 1)
	assert.Equal(t, 100, result[0].ScheduleId)
	assert.Equal(t, 200, result[0].SearchId)
}

func TestMapSchedules_MultipleSchedules(t *testing.T) {
	// Test mapping multiple schedules
	dbSchedules := []db.GetSchedulesResponse{
		testutil.CreateTestSchedule(1, 10, 5),
		testutil.CreateTestSchedule(2, 20, 15),
		testutil.CreateTestSchedule(3, 30, 30),
		testutil.CreateTestSchedule(4, 40, 60),
		testutil.CreateTestSchedule(5, 50, 120),
	}

	result := MapSchedules(dbSchedules)

	assert.Len(t, result, 5)
	for i, schedule := range result {
		assert.Equal(t, i+1, schedule.ScheduleId)
		assert.Equal(t, (i+1)*10, schedule.SearchId)
	}
}

func TestMapSchedules_EmptySlice(t *testing.T) {
	// Test mapping an empty slice
	var dbSchedules []db.GetSchedulesResponse

	result := MapSchedules(dbSchedules)

	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestMapSchedules_FieldMapping(t *testing.T) {
	// Test that only ScheduleId and SearchId are mapped
	futureTime := time.Now().Add(24 * time.Hour)
	dbSchedules := []db.GetSchedulesResponse{
		testutil.CreateTestScheduleWithTime(100, 200, 30, futureTime),
	}

	result := MapSchedules(dbSchedules)

	assert.Len(t, result, 1)
	assert.Equal(t, 100, result[0].ScheduleId)
	assert.Equal(t, 200, result[0].SearchId)
	// Verify NextSearch and SearchDelayMinutes are not in the result
	// (implicitly tested by the struct having only ScheduleId and SearchId)
}

func TestMapSchedules_OrderPreserved(t *testing.T) {
	// Test that order is preserved
	dbSchedules := []db.GetSchedulesResponse{
		testutil.CreateTestSchedule(1, 10, 5),
		testutil.CreateTestSchedule(2, 20, 15),
		testutil.CreateTestSchedule(3, 30, 30),
	}

	result := MapSchedules(dbSchedules)

	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].ScheduleId)
	assert.Equal(t, 2, result[1].ScheduleId)
	assert.Equal(t, 3, result[2].ScheduleId)
}

func TestMapSchedule_AllFields(t *testing.T) {
	// Test mapping a single schedule with all fields
	dbSchedule := testutil.CreateTestSchedule(999, 888, 45)

	result := mapSchedule(dbSchedule)

	assert.Equal(t, 999, result.ScheduleId)
	assert.Equal(t, 888, result.SearchId)
}

func TestMapSchedule_ZeroValues(t *testing.T) {
	// Test mapping with zero values (edge case)
	dbSchedule := testutil.CreateTestSchedule(0, 0, 0)

	result := mapSchedule(dbSchedule)

	assert.Equal(t, 0, result.ScheduleId)
	assert.Equal(t, 0, result.SearchId)
}

func TestMapSchedules_LargeDataSet(t *testing.T) {
	// Test mapping a large dataset
	dbSchedules := make([]db.GetSchedulesResponse, 1000)
	for i := 0; i < 1000; i++ {
		dbSchedules[i] = testutil.CreateTestSchedule(i+1, (i+1)*100, 30)
	}

	result := MapSchedules(dbSchedules)

	assert.Len(t, result, 1000)
	// Spot check first, middle, and last elements
	assert.Equal(t, 1, result[0].ScheduleId)
	assert.Equal(t, 100, result[0].SearchId)
	assert.Equal(t, 500, result[499].ScheduleId)
	assert.Equal(t, 50000, result[499].SearchId)
	assert.Equal(t, 1000, result[999].ScheduleId)
	assert.Equal(t, 100000, result[999].SearchId)
}
