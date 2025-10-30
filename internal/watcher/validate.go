package watcher

import (
	"regexp"
	"time"

	"git.culab.ru/course-projects/path/config"
	"git.culab.ru/course-projects/path/internal/out"
)

func validateRecord(record out.Record, args config.Args) (bool, error) {
	valid := true
	if args.FilterRE != nil {
		valid, _ = validateByRegex(record, args.FilterRE)
	}
	if args.SinceT != nil {
		valid, _ = validateByTime(record, args.SinceT)
	}
	return valid, nil
}

func validateByRegex(record out.Record, regexp *regexp.Regexp) (bool, error) {
	if regexp.MatchString(record.Path) {
		return true, nil
	}
	return false, nil
}

func validateByTime(record out.Record, since *time.Time) (bool, error) {
	if record.Time.Before(*since) {
		return false, nil
	}
	return true, nil
}
