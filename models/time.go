package models

import (
	"strings"
	"time"
)

// Time is custom time format used in all responses
type Time struct {
	time.Time
}

// List of known time formats
var (
	ctLayouts      = []string{"2006-01-02", "2006-01-02 15:04:05"}
	ctZonedLayouts = []string{"2006-01-02T15:04:05-0700"}
)

// UnmarshalJSON parses JSON time string with custom time formats
func (t *Time) UnmarshalJSON(b []byte) error {
	pTime, err := parseTime(string(b))
	if err != nil {
		return err
	}

	t.Time = pTime
	return nil
}

// UnmarshalCSV converts CSV string field internal date
func (t *Time) UnmarshalCSV(s string) error {
	pTime, err := parseTime(s)
	if err != nil {
		return err
	}

	t.Time = pTime
	return nil
}

func parseTime(ts string) (time.Time, error) {
	var (
		pTime time.Time
		err   error
	)

	s := strings.TrimSpace(strings.Trim(ts, "\""))
	if s == "" || s == "null" {
		return pTime, nil
	}

	// Load IST location.
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return pTime, err
	}

	// Iterate through zoneless layouts and assign zone as IST.
	for _, l := range ctLayouts {
		pTime, err = time.ParseInLocation(l, s, loc)
		if err == nil && !pTime.IsZero() {
			break
		}
	}

	// If pattern not found then iterate and parse layouts with zone.
	if pTime.IsZero() {
		for _, l := range ctZonedLayouts {
			pTime, err = time.Parse(l, s)
			if err == nil && !pTime.IsZero() {
				break
			}
		}
	}

	return pTime, err
}
