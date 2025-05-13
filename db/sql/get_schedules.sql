select ID as ScheduleId,
    SearchId,
    NextSearch,
    SearchDelayMinutes
FROM Schedule
WHERE NextSearch < NOW()