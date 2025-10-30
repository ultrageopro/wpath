package config

import (
	"regexp"
	"time"
)

type Args struct {
	FlagPath       string
	FlagFilterName string
	FlagSinceStr   string
	FlagNoColor    bool

	FilterRE *regexp.Regexp
	SinceT   *time.Time
}
