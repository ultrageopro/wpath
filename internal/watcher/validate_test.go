package watcher

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ultrageopro/wpath/config"
	"github.com/ultrageopro/wpath/internal/out"
)

func TestValidateByTime(t *testing.T) {
	t.Parallel()

	_, err := validateByTime(out.Record{}, nil)
	assert.Error(t, err)

	current := time.Now()
	_, err = validateByTime(out.Record{Time: current}, &current)
	assert.NoError(t, err)

	ok, err := validateByTime(out.Record{Time: current.Add(time.Hour)}, &current)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = validateByTime(out.Record{Time: current.Add(-time.Hour)}, &current)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestValidateByRegex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rec     out.Record
		re      *regexp.Regexp
		wantOK  bool
		wantErr bool
	}{
		{
			name:    "nil regex -> error",
			rec:     out.Record{Path: "test"},
			re:      nil,
			wantErr: true,
		},
		{
			name:   "exact match",
			rec:    out.Record{Path: "test"},
			re:     regexp.MustCompile(`^test$`),
			wantOK: true,
		},
		{
			name:   "not matched",
			rec:    out.Record{Path: "test"},
			re:     regexp.MustCompile(`^test2$`),
			wantOK: false,
		},
		{
			name:   "empty pattern matches anything (including empty)",
			rec:    out.Record{Path: "anything"},
			re:     regexp.MustCompile(``),
			wantOK: true,
		},
		{
			name:   "prefix match",
			rec:    out.Record{Path: "test/file.txt"},
			re:     regexp.MustCompile(`^test`),
			wantOK: true,
		},
		{
			name:   "case-insensitive",
			rec:    out.Record{Path: "TeSt"},
			re:     regexp.MustCompile(`(?i)^test$`),
			wantOK: true,
		},
		{
			name:   "special chars dot literal",
			rec:    out.Record{Path: "file.name"},
			re:     regexp.MustCompile(`^file\.name$`),
			wantOK: true,
		},
		{
			name:   "unicode path",
			rec:    out.Record{Path: "путь/файл.txt"},
			re:     regexp.MustCompile(`файл\.txt$`),
			wantOK: true,
		},
		{
			name:   "empty path no match",
			rec:    out.Record{Path: ""},
			re:     regexp.MustCompile(`^.+$`),
			wantOK: false,
		},
		{
			name:   "empty path with empty pattern -> true",
			rec:    out.Record{Path: ""},
			re:     regexp.MustCompile(``),
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ok, err := validateByRegex(tt.rec, tt.re)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOK, ok)
		})
	}
}

func TestValidateRecord(t *testing.T) {
	t.Parallel()

	now := time.Now()
	before, after := now.Add(-time.Hour), now.Add(time.Hour)

	tests := []struct {
		name string
		rec  out.Record
		args config.Args
		want bool
	}{
		{
			name: "filter by regex",
			rec:  out.Record{Path: "test"},
			args: config.Args{FilterRE: regexp.MustCompile(`^test$`)},
			want: true,
		},
		{
			name: "filter by time",
			rec:  out.Record{Path: "test", Time: now},
			args: config.Args{SinceT: &before},
			want: true,
		},
		{
			name: "filter by regex and time",
			rec:  out.Record{Path: "test", Time: now},
			args: config.Args{FilterRE: regexp.MustCompile(`^test$`), SinceT: &before},
			want: true,
		},
		{
			name: "filter by regex and time but time is too old",
			rec:  out.Record{Path: "test", Time: now},
			args: config.Args{FilterRE: regexp.MustCompile(`^test$`), SinceT: &after},
			want: false,
		},
		{
			name: "no filter",
			rec:  out.Record{Path: "test"},
			args: config.Args{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, validateRecord(tt.rec, tt.args))
		})
	}
}
