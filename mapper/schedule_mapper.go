package mapper

import (
	"flatchecker-scheduler/db"
	"flatchecker-scheduler/pubsublib"
)

func MapSchedules(dbSchedules []db.GetSchedulesResponse) []pubsublib.ScheduledSearchesMessage {
	pubsubSchedules := make([]pubsublib.ScheduledSearchesMessage, len(dbSchedules))

	for i, dbSchedule := range dbSchedules {
		pubsubSchedules[i] = mapSchedule(dbSchedule)
	}

	return pubsubSchedules
}

func mapSchedule(schedule db.GetSchedulesResponse) pubsublib.ScheduledSearchesMessage {
	return pubsublib.ScheduledSearchesMessage{
		ScheduleId: schedule.ScheduleId,
		SearchId:   schedule.SearchId,
	}
}
