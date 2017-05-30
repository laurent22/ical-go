package ical

import (
	"testing"
)

func TestCalendarSerialize(t *testing.T) {
	calendar := new(Calendar)

	// test calendar w/o items

	expected := "BEGIN:VCALENDAR\nVERSION:2.0\nEND:VCALENDAR"
	output := calendar.Serialize()

	if output != expected {
		t.Error("\nExpected calendar serialization to be:\n", expected, "\n\nbut got:\n", output)
	}

	// test calendar with items

	calendar.Items = append(calendar.Items, CalendarEvent{Summary: "Foo"})

	expected = "BEGIN:VCALENDAR\nVERSION:2.0\nBEGIN:VEVENT\nSUMMARY:Foo\nEND:VEVENT\nEND:VCALENDAR"
	output = calendar.Serialize()

	if output != expected {
		t.Error("\nExpected calendar serialization to be:\n", expected, "\n\nbut got:\n", output)
	}
}
