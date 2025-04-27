select ID as ScheduleId,
    SearchId
FROM Schedule
WHERE LastSearched < NOW()