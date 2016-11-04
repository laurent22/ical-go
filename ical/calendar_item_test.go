package ical

import (
  "strings"
  "testing"
  "time"
)

func TestStartAndEndAtUTC(t *testing.T) {
  item := CalendarItem{}

  if item.StartAtUTC() != nil {
    t.Error("StartAtUTC should have been nil")
  }
  if item.EndAtUTC() != nil {
    t.Error("EndAtUTC should have been nil")
  }

  tUTC := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
  item.StartAt = &tUTC
  item.EndAt   = &tUTC
  startTime := *(item.StartAtUTC())
  endTime   := *(item.EndAtUTC())

  if startTime != tUTC {
    t.Error("StartAtUTC should have been", tUTC, ", but was", startTime)
  }
  if endTime != tUTC {
    t.Error("EndAtUTC should have been", tUTC, ", but was", endTime)
  }

  tUTC = time.Date(2010, time.March, 8, 2, 0, 0, 0, time.UTC)
  nyk, err := time.LoadLocation("America/New_York")
  if err != nil {
    panic(err)
  }
  tNYK := tUTC.In(nyk)
  item.StartAt = &tNYK
  item.EndAt = &tNYK
  startTime = *(item.StartAtUTC())
  endTime   = *(item.EndAtUTC())

  if startTime != tUTC {
    t.Error("StartAtUTC should have been", tUTC, ", but was", startTime)
  }
  if endTime != tUTC {
    t.Error("EndAtUTC should have been", tUTC, ", but was", endTime)
  }
}

func TestCalendarItemSerialize(t *testing.T) {
  ny, err := time.LoadLocation("America/New_York")
  if err != nil {
    panic(err)
  }

  createdAt := time.Date(2010, time.January, 1, 12, 0, 1, 0, time.UTC)
  modifiedAt := createdAt.Add(time.Second)
  startsAt    := createdAt.Add(time.Second * 2).In(ny)
  endsAt      := createdAt.Add(time.Second * 3).In(ny)

  item := CalendarItem {
    Id: "123",
    CreatedAtUTC: &createdAt,
    ModifiedAtUTC: &modifiedAt,
    StartAt: &startsAt,
    EndAt: &endsAt,
    Summary: "Foo Bar",
    Location: "Berlin\nGermany",
    Description: "Lorem\nIpsum",
  }

  // expects that DTSTART and DTEND be in UTC (Z)
  // expects that string values (LOCATION for example) be escaped
  expected := `
BEGIN:VEVENT
UID:123
CREATED:20100101T120001Z
LAST-MODIFIED:20100101T120002Z
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
DESCRIPTION:Lorem\nIpsum
LOCATION:Berlin\nGermany
END:VEVENT`

  output := item.Serialize()
  if output != strings.TrimSpace(expected) {
    t.Error("Expected calendar item serialization to be:\n", expected, "\n\nbut got:\n", output)
  }
}
