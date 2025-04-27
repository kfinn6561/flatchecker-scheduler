select ID as ScheduleId,
    SearchId
FROM Schedule
WHERE NextSearch < NOW()