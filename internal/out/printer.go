package out

import (
	"fmt"

	"github.com/gosuri/uilive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type PathPrinter struct {
	noColors    bool
	tableWriter table.Writer
	liveWriter  *uilive.Writer
}

func rowPainter(row table.Row) text.Colors {
	event := row[1].(Event)
	switch event {
	case EventCreate:
		return text.Colors{text.BgHiGreen, text.FgBlack}
	case EventDelete:
		return text.Colors{text.BgHiRed, text.FgWhite}
	case EventModify:
		return text.Colors{text.BgHiYellow, text.FgBlack}
	case EventChmod:
		return text.Colors{text.BgHiCyan, text.FgBlack}
	default:
		return nil
	}
}

func NewPathPrinter(noColors bool) *PathPrinter {
	t := table.NewWriter()
	t.SetOutputMirror(nil)

	t.AppendHeader(table.Row{"TIME", "EVENT", "PATH"})
	if !noColors {
		t.SetRowPainter(rowPainter)
	}

	writer := uilive.New()
	writer.Start()

	return &PathPrinter{
		noColors:    noColors,
		tableWriter: t,
		liveWriter:  writer,
	}
}

func (p *PathPrinter) Print(newRecord Record) {
	p.tableWriter.AppendRow(newRecord.ToTableRow())
	output := p.tableWriter.Render()
	fmt.Fprintln(p.liveWriter, output)
	p.liveWriter.Flush()
}

func (p *PathPrinter) Stop() {
	p.liveWriter.Stop()
}
