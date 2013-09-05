package ical

import (
	"strings"
	"time"	
	"strconv"	
	"errors"
	"log"
)

// ~/Library/Calendars/05C38275-CBA9-4321-8457-2DC326799AE8.calendar/Events

// BEGIN:VCALENDAR
// VERSION:2.0
// CALSCALE:GREGORIAN
// BEGIN:VTODO
// STATUS:COMPLETED
// CREATED:20130812T125803Z
// UID:4E00D9DC-5AF8-43DA-B543-EACB8182CA97
// SUMMARY:dghshdfjhf
// COMPLETED:20130812T125805Z
// DTSTAMP:20130812T125803Z
// PRIORITY:1
// PERCENT-COMPLETE:100
// SEQUENCE:0
// DESCRIPTION:erezrare\nrateztaz\ntrois
// END:VTODO
// END:VCALENDAR


// STATUS: NEEDS-ACTION / COMPLETE
// PERCENT-COMPLETE - enlever / 100



// DTSTART;TZID=Asia/Harbin:20130817T000000
// DUE;TZID=Asia/Harbin:20130817T000000
// BEGIN:VALARM
// X-WR-ALARMUID:FEEACF1F-D546-462C-A3E2-F9706F44BB33
// UID:FEEACF1F-D546-462C-A3E2-F9706F44BB33
// TRIGGER;VALUE=DATE-TIME:20130816T160000Z
// DESCRIPTION:Avertissement événement
// ACTION:DISPLAY
// END:VALARM


// TODO: TIME ZONE

func EncodeDateProperty(name string, t time.Time) string {
	var output string
	zone, _ := t.Zone()
	if zone != "UTC" && zone != "" {
		output = ";TZID=" + zone + ":" + t.Format("20060102T150405")
	} else {
		output = ":" + t.Format("20060102T150405") + "Z"
	}
	return name + output
}

func EscapeTextType(input string) string {
	output := strings.Replace(input, "\\", "\\\\", -1)
	output = strings.Replace(output, ";", "\\;", -1)
	output = strings.Replace(output, ",", "\\,", -1)
	output = strings.Replace(output, "\n", "\\n", -1)	
	return output
}

type Node struct {
	Name string
	Value string
	Type int // 1 = Object, 0 = Name/Value
	Parameters map[string]string
	Children []*Node
}

func (this *Node) Parameter(name string, defaultValue string) string {
	if len(this.Parameters) <= 0 { return defaultValue }
	v, ok := this.Parameters[name]
	if !ok { return defaultValue }
	return v
} 

func (this *Node) ChildrenByName(name string) []*Node {
	var output []*Node
	for _, child := range this.Children {
		if child.Name == name {
			output = append(output, child)
		}
	}
	return output
}

func (this *Node) ChildByName(name string) *Node {
	for _, child := range this.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (this *Node) PropString(name string, defaultValue string) string {
	for _, child := range this.Children {
		if child.Name == name {
			return child.Value
		}
	}
	return defaultValue
}

func (this *Node) PropDate(name string, defaultValue time.Time) time.Time {
	// 20130816T160000Z
	// Mon Jan 2 15:04:05 -0700 MST 2006
	// DUE;TZID=Asia/Harbin:20130817T000000
	
	node := this.ChildByName(name)
	if node == nil { return defaultValue }
	tzid := node.Parameter("TZID", "")
	var output time.Time
	var err error
	if tzid != "" {
		loc, err := time.LoadLocation(tzid)
		if err != nil { panic(err) }
		output, err = time.ParseInLocation("20060102T150405", node.Value, loc)		
	} else {
		output, err = time.Parse("20060102T150405Z", node.Value)
	}
	
	if err != nil { panic(err) }
	return output
}

func (this *Node) PropInt(name string, defaultValue int) int {
	n := this.PropString(name, "")
	if n == "" { return defaultValue }
	output, err := strconv.Atoi(n)
	if err != nil { panic(err) }
	return output
}

func (this *Node) String() string {
	s := ""
	if this.Type == 1 {
		s += "===== " + this.Name
		s += "\n"
	} else {
		s += this.Name
		s += ":" + this.Value
		s += "\n"		
	}
	for _, child := range this.Children {
		s += child.String()
	}
	if this.Type == 1 {
		s += "===== /" + this.Name
		s += "\n"
	}
	return s
}

func UnescapeTextType(s string) string {
	s = strings.Replace(s, "\\;", ";", -1)
	s = strings.Replace(s, "\\,", ",", -1)
	s = strings.Replace(s, "\\n", "\n", -1)	
	s = strings.Replace(s, "\\\\", "\\", -1)	
	return s
}

func ParseTextType(lines []string, lineIndex int) (string, int) {
	line := lines[lineIndex]
	colonIndex := strings.Index(line, ":")
	output := strings.TrimSpace(line[colonIndex+1:len(line)])
	lineIndex++
	for {
		line := lines[lineIndex]
		if line == "" || line[0] != ' ' {
			return UnescapeTextType(output), lineIndex
		}
		output += line[1:len(line)]
		lineIndex++
	}
	return UnescapeTextType(output), lineIndex
}

func EncodeTextType(s string) string {
	output := ""
	s = EscapeTextType(s)
	lineLength := 0
	for _, c := range s {
		if lineLength + len(string(c)) > 75 {
			output += "\n "
			lineLength = 1
		}
		output += string(c)
		lineLength += len(string(c))
	}
	return output
}

// BEGIN:VCALENDAR
// VERSION:2.0
// CALSCALE:GREGORIAN
// BEGIN:VTODO
// STATUS:COMPLETED
// CREATED:20130812T125803Z
// UID:4E00D9DC-5AF8-43DA-B543-EACB8182CA97
// SUMMARY:dghshdfjhf
// COMPLETED:20130812T125805Z
// DTSTAMP:20130812T125803Z
// PRIORITY:1
// PERCENT-COMPLETE:100
// SEQUENCE:0
// DESCRIPTION:erezrare\nrateztaz\ntrois
// END:VTODO
// END:VCALENDAR

func TodoFromNode(node *Node) Todo {
	if node.Name != "VTODO" { panic("Node is not a VTODO") }

	var todo Todo
	todo.SetId(node.PropString("UID", ""))
	todo.SetSummary(node.PropString("SUMMARY", ""))
	todo.SetDescription(node.PropString("DESCRIPTION", ""))
	todo.SetDueDate(node.PropDate("DUE", time.Time{}))
	//todo.SetAlarmDate(this.TimestampBytesToTime(reminderDate))
	todo.SetCreatedDate(node.PropDate("CREATED", time.Time{}))
	todo.SetModifiedDate(node.PropDate("DTSTAMP", time.Time{}))
	todo.SetPriority(node.PropInt("PRIORITY", 0))
	todo.SetPercentComplete(node.PropInt("PERCENT-COMPLETE", 0))
	return todo	
}

func ParseCalendarNode(lines []string, lineIndex int) (*Node, bool, error, int) {
	line := strings.TrimSpace(lines[lineIndex])
	_ = log.Println
	colonIndex := strings.Index(line, ":")
	if colonIndex <= 0 {
		return nil, false, errors.New("Invalid value/pair: " + line), lineIndex + 1
	}
	name := line[0:colonIndex]
	splitted := strings.Split(name, ";")
	var parameters map[string]string
	if len(splitted) >= 2 {
		name = splitted[0]
		parameters = make(map[string]string)
		for i := 1; i < len(splitted); i++ {
			p := strings.Split(splitted[i], "=")
			if len(p) != 2 { panic("Invalid parameter format: " + name) } 
			parameters[p[0]] = p[1]
		}	
	}
	value := line[colonIndex+1:len(line)]
	
	if name == "BEGIN" {
		node := new(Node)
		node.Name = value
		node.Type = 1
		lineIndex = lineIndex + 1
		for {
			child, finished, _, newLineIndex := ParseCalendarNode(lines, lineIndex)
			if finished {
				return node, false, nil, newLineIndex
			} else {
				if child != nil {
					node.Children = append(node.Children, child)
				}
				lineIndex = newLineIndex
			}
		}
	} else if name == "END" {
		return nil, true, nil, lineIndex + 1
	} else {
		node := new(Node)
		node.Name = name
		if name == "DESCRIPTION" || name == "SUMMARY" {
			text, newLineIndex := ParseTextType(lines, lineIndex)
			node.Value = text
			node.Parameters = parameters
			return node, false, nil, newLineIndex
		} else {
			node.Value = value
			node.Parameters = parameters
			return node, false, nil, lineIndex + 1
		}
	}
	
	panic("Unreachable")
	return nil, false, nil, lineIndex + 1
}

type Calendar struct {
	Items []CalendarItem
}

type CalendarItem struct {
	id string
	summary string
	description string
	priority int // 0..9 (O -> none, 1 -> highest, 9 -> lowest)
	percentComplete int
	createdDate time.Time
	modifiedDate time.Time
	completedDate time.Time
	startDate time.Time
	alarmDate time.Time
	sequence int
}

func (this *CalendarItem) SetId(v string) { this.id = v }
func (this *CalendarItem) Id() string {	return this.id }
func (this *CalendarItem) SetSummary(v string) { this.summary = v }
func (this *CalendarItem) Summary() string { return this.summary }
func (this *CalendarItem) SetDescription(v string) { this.description = v }
func (this *CalendarItem) Description() string { return this.description }
func (this *CalendarItem) SetPriority(v int) { this.priority = v }
func (this *CalendarItem) Priority() int { return this.priority }
func (this *CalendarItem) SetPercentComplete(v int) { this.percentComplete = v }
func (this *CalendarItem) PercentComplete() int { return this.percentComplete }
func (this *CalendarItem) SetCreatedDate(v time.Time) { this.createdDate = v }
func (this *CalendarItem) CreatedDate() time.Time { return this.createdDate }
func (this *CalendarItem) SetModifiedDate(v time.Time) { this.modifiedDate = v }
func (this *CalendarItem) ModifiedDate() time.Time { return this.modifiedDate }
func (this *CalendarItem) SetCompletedDate(v time.Time) { this.completedDate = v }
func (this *CalendarItem) CompletedDate() time.Time { return this.completedDate }
func (this *CalendarItem) SetStartDate(v time.Time) { this.startDate = v }
func (this *CalendarItem) StartDate() time.Time { return this.startDate }
func (this *CalendarItem) SetAlarmDate(v time.Time) { this.alarmDate = v }
func (this *CalendarItem) AlarmDate() time.Time { return this.alarmDate }
func (this *CalendarItem) SetSequence(v int) { this.sequence = v }
func (this *CalendarItem) Sequence() int { return this.sequence }

type Todo struct {
	CalendarItem
	dueDate time.Time
}

func (this *Todo) SetDueDate(v time.Time) { this.dueDate = v }
func (this *Todo) DueDate() time.Time { return this.dueDate }

func (this *Todo) ICalString(target string) string {
	s := "BEGIN:VTODO\n"
	
	if target == "macTodo" {
		status := "NEEDS-ACTION"
		if this.PercentComplete() == 100 {
			status = "COMPLETED"
		}
		s += "STATUS:" + status + "\n"
	}
	
	s += EncodeDateProperty("CREATED", this.CreatedDate()) + "\n"
	s += "UID:" + this.Id() + "\n"
	s += "SUMMARY:" + EscapeTextType(this.Summary()) + "\n"
	if this.PercentComplete() == 100 && !this.CompletedDate().IsZero() {
		s += EncodeDateProperty("COMPLETED", this.CompletedDate()) + "\n"
	}
	s += EncodeDateProperty("DTSTAMP", this.ModifiedDate()) + "\n"
	if this.Priority() != 0 {
		s += "PRIORITY:" + strconv.Itoa(this.Priority()) + "\n"
	}
	if this.PercentComplete() != 0 {
		s += "PERCENT-COMPLETE:" + strconv.Itoa(this.PercentComplete()) + "\n"
	}
	if target == "macTodo" {
		s += "SEQUENCE:" + strconv.Itoa(this.Sequence()) + "\n"
	}
	if this.Description() != "" {
		s += "DESCRIPTION:" + EncodeTextType(this.Description()) + "\n"
	}
	
	s += "END:VTODO\n"
	
	return s
}