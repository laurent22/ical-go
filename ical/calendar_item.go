package ical

import (
	"time"
)

type CalendarItem struct {
	Id string
	Summary string
	Location string
	CreatedAtUTC *time.Time
	ModifiedAtUTC *time.Time
	StartAt *time.Time
	EndAt *time.Time
}

func (this *CalendarItem) StartAtUTC() *time.Time {
	return inUTC(this.StartAt)
}

func (this *CalendarItem) EndAtUTC() *time.Time {
	return inUTC(this.EndAt)
}

func (this *CalendarItem) Serialize() string {
	buffer := new(strBuffer)
	return this.serializeWithBuffer(buffer)
}

func (this *CalendarItem) AsICS() string {
	calendar := Calendar {
		Items: []CalendarItem{*this},
	}

	return calendar.Serialize()
}

func (this *CalendarItem) serializeWithBuffer(buffer *strBuffer) string {
	serializer := calItemSerializer{
		item: this,
		buffer: buffer,
	}
	return serializer.serialize()
}

func inUTC(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	tUTC := t.UTC()
	return &tUTC
}
