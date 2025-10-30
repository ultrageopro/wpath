package out

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Event string

const (
	EventCreate Event = "CREATED"
	EventDelete Event = "DELETED"
	EventModify Event = "MODIFIED"
	EventChmod  Event = "CHMOD"
)

type Record struct {
	Time  time.Time
	Event Event
	Path  string
}

func (r Record) ToTableRow() table.Row {
	return table.Row{
		r.Time.Format(time.RFC3339),
		r.Event,
		r.Path,
	}
}

func NewRecord(time time.Time, event Event, path string) Record {
	return Record{
		Time:  time,
		Event: event,
		Path:  path,
	}
}
