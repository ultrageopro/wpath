package watcher

import (
	"fmt"
	"regexp"
	"time"

	"git.culab.ru/course-projects/path/config"
	"git.culab.ru/course-projects/path/internal/out"
)

func validateRecord(record out.Record, args config.Args) bool {
	valid := true
	if args.FilterRE != nil {
		valid, _ = validateByRegex(record, args.FilterRE)
	}
	if args.SinceT != nil {
		valid, _ = validateByTime(record, args.SinceT)
	}
	return valid
}

func validateByRegex(record out.Record, regexp *regexp.Regexp) (bool, error) {
	if regexp == nil {
		return true, fmt.Errorf("regexp is nil")
	}

	if regexp.MatchString(record.Path) {
		return true, nil
	}
	return false, nil
}

func validateByTime(record out.Record, since *time.Time) (bool, error) {
	if since == nil {
		return true, fmt.Errorf("since is nil")
	}

	if record.Time.Before(*since) {
		return false, nil
	}
	return true, nil
}
