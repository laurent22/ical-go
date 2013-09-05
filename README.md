# ical-go

iCal package for Go (Golang)

## Installation

    go get https://github.com/laurent22/ical-go

## Status

Currently, the package doesn't support the full iCal specification. Todo are partially supported and Events not really.

The most useful function in the package is `ParseCalendarNode` which parses a VCALENDAR string, unwrap and unfold lines, etc. and put all this into a usable structure (a collection of nodes with name, value, type, etc.). From that, it's easy to support whatever feature is required.

## License

MIT
