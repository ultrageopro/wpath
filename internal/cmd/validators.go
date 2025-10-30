package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func validatePath(flagPath string) error {
	if flagPath == "" {
		return errors.New("--path is required")
	}

	abs, err := filepath.Abs(flagPath)
	if err != nil {
		return fmt.Errorf("resolve --path: %w", err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("stat --path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("--path must be a directory: %s", abs)
	}

	return nil
}

func validateFilter(flagFilterName string) error {
	if flagFilterName != "" {
		re, err := regexp.Compile(flagFilterName)
		if err != nil {
			return fmt.Errorf("--filter-name: invalid regexp: %w", err)
		}
		Args.FilterRE = re
	}
	return nil
}

func validateSince(flagSinceStr string) error {
	if flagSinceStr != "" {
		if t, err := time.Parse(time.RFC3339, flagSinceStr); err == nil {
			Args.SinceT = &t
		} else if sec, err := strconv.ParseInt(flagSinceStr, 10, 64); err == nil {
			t := time.Unix(sec, 0).UTC()
			Args.SinceT = &t
		} else {
			return fmt.Errorf("--since: use RFC3339 (e.g. 2025-10-30T12:34:56Z) or UNIX seconds")
		}
	}
	return nil
}
