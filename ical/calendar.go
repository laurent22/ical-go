package ical

type Calendar struct {
	Items []CalendarItem
}

func (this *Calendar) Serialize() string {
	serializer := calSerializer{
		calendar: this,
		buffer: new(strBuffer),
	}
	return serializer.serialize()
}

func (this *Calendar) ToICS() string {
	return this.Serialize()
}
