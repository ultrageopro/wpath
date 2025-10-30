package watcher

import (
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/notEpsilon/go-pair"
	"github.com/ultrageopro/wpath/config"
	"github.com/ultrageopro/wpath/internal/out"
)

type Processor struct {
	Watcher    *fsnotify.Watcher
	Printer    *out.PathPrinter
	logger     *log.Logger
	Args       config.Args
	Operations chan pair.Pair[string, out.Event]
}

func NewProcessor(printer *out.PathPrinter, args config.Args) (*Processor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Processor{
		Watcher:    watcher,
		Printer:    printer,
		Operations: make(chan pair.Pair[string, out.Event]),
		logger:     log.Default(),
		Args:       args,
	}, nil
}

func (p *Processor) updateDirs(path string) error {
	return filepath.WalkDir(
		path,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				p.logger.Printf("walkdir error: %v", err)
				return err
			}
			if d.IsDir() {
				p.Watcher.Add(path)
			}
			return nil
		},
	)
}

func (p *Processor) watch() {
	defer close(p.Operations)
	for {
		select {
		case event, ok := <-p.Watcher.Events:
			if !ok {
				return
			}
			if event.Op.Has(fsnotify.Create) {
				p.Operations <- pair.Pair[string, out.Event]{First: event.Name, Second: out.EventCreate}
			} else if event.Op.Has(fsnotify.Remove) {
				p.Operations <- pair.Pair[string, out.Event]{First: event.Name, Second: out.EventDelete}
			} else if event.Op.Has(fsnotify.Write) {
				p.Operations <- pair.Pair[string, out.Event]{First: event.Name, Second: out.EventModify}
			} else if event.Op.Has(fsnotify.Chmod) {
				p.Operations <- pair.Pair[string, out.Event]{First: event.Name, Second: out.EventChmod}
			}
		case err, ok := <-p.Watcher.Errors:
			if !ok {
				return
			}
			p.logger.Printf("watcher error: %v", err)
		}
	}
}

func (p *Processor) processEvent(event pair.Pair[string, out.Event], mu *sync.Mutex, path string) {
	mu.Lock()
	r := out.NewRecord(
		time.Now(),
		event.Second,
		event.First,
	)
	if validateRecord(r, p.Args) {
		p.Printer.Print(r)
	}
	p.updateDirs(path)
	mu.Unlock()
}

func (p *Processor) Watch(path string, mu *sync.Mutex) error {
	err := p.Watcher.Add(path)
	if err != nil {
		return err
	}
	defer p.Watcher.Close()

	go p.watch()
	for event := range p.Operations {
		p.processEvent(event, mu, path)
	}
	return err
}
