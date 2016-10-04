package ical

import (
  "time"
  "strings"
)

type calSerializer struct {
	calendar *Calendar
	buffer *strBuffer
}

func (this *calSerializer) serialize() string {
	this.serializeCalendar()
	return strings.TrimSpace(this.buffer.String())
}

func (this *calSerializer) serializeCalendar() {
	this.begin()
	this.version()
	this.items()
	this.end()
}

func (this *calSerializer) begin() {
	this.buffer.Write("BEGIN:VCALENDAR\n")
}

func (this *calSerializer) end() {
	this.buffer.Write("END:VCALENDAR\n")
}

func (this *calSerializer) version() {
	this.buffer.Write("VERSION:2.0\n")
}

func (this *calSerializer) items() {
	for _, item := range this.calendar.Items {
		item.serializeWithBuffer(this.buffer)
	}
}

type calItemSerializer struct {
	item *CalendarItem
	buffer *strBuffer
}

const (
	itemSerializerTimeFormat = "20060102T150405Z"
)

func (this *calItemSerializer) serialize() string {
	this.serializeItem()
	return strings.TrimSpace(this.buffer.String())
}

func (this *calItemSerializer) serializeItem() {
	this.begin()
	this.uid()
	this.created()
	this.lastModified()
	this.dtstart()
	this.dtend()
	this.summary()
	this.location()
	this.end()
}

func (this *calItemSerializer) begin() {
	this.buffer.Write("BEGIN:VEVENT\n")
}

func (this *calItemSerializer) end() {
	this.buffer.Write("END:VEVENT\n")
}

func (this *calItemSerializer) uid() {
	this.serializeStringProp("UID", this.item.Id)
}

func (this *calItemSerializer) summary() {
	this.serializeStringProp("SUMMARY", this.item.Summary)
}

func (this *calItemSerializer) location() {
	this.serializeStringProp("LOCATION", this.item.Location)
}

func (this *calItemSerializer) dtstart() {
	this.serializeTimeProp("DTSTART", this.item.StartAtUTC())
}

func (this *calItemSerializer) dtend() {
	this.serializeTimeProp("DTEND", this.item.EndAtUTC())
}

func (this *calItemSerializer) created() {
	this.serializeTimeProp("CREATED", this.item.CreatedAtUTC)
}

func (this *calItemSerializer) lastModified() {
	this.serializeTimeProp("LAST-MODIFIED", this.item.ModifiedAtUTC)
}

func (this *calItemSerializer) serializeStringProp(name, value string) {
	if value != "" {
		escapedValue := escapeTextType(value)
		this.buffer.Write("%s:%s\n", name, escapedValue)
	}
}

func (this *calItemSerializer) serializeTimeProp(name string, value *time.Time) {
	if value != nil {
		this.buffer.Write("%s:%s\n", name, value.Format(itemSerializerTimeFormat))
	}
}
