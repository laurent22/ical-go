package ical

import "testing"

func TestCorruptParameters(t *testing.T) {
	data :=
`
BEGIN:VCALENDAR
BEGIN:VEVENT
UID;:123
DTSTART:20100101T120003Z
DTEND:20100101T120004Z
SUMMARY:Foo Bar
END:VEVENT
END:VCALENDAR
`
	_, err := ParseCalendar(data)

	errMsg := "Invalid parameter format: UID"
	if err.Error() != errMsg {
		t.Error("wrong error message, got: " + err.Error() + " but expected: " + errMsg )
	}
}

func TestMissingEnd(t *testing.T) {
data :=
`
BEGIN:VCALENDAR
BEGIN:VEVENT
UID:123
END:VCALENDAR
`
	_, err := ParseCalendar(data)

	errMsg := "Missing END tag"
	if err.Error() != errMsg {
		t.Error("wrong error message, got: " + err.Error() + " but expected: " + errMsg )
	}
}
