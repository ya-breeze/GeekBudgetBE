package utils

import (
	"time"
)

type Granularity string

const (
	GranularityMonth Granularity = "month"
	GranularityYear  Granularity = "year"
)

func RoundToGranularity(date time.Time, granularity Granularity, roundUp bool,
) time.Time {
	var rounded time.Time
	if granularity == GranularityMonth {
		rounded = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	} else if granularity == GranularityYear {
		rounded = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
	}

	if !roundUp {
		return rounded
	}

	// If calculating ceiling (roundUp) and validation passes - return rounded
	// Check if date is already rounded (equal to rounded).
	// Note: Equal compares time instant. We constructed rounded with same Location.
	if date.Equal(rounded) {
		return rounded
	}

	add := 1
	if granularity == GranularityMonth {
		return rounded.AddDate(0, add, 0)
	} else if granularity == GranularityYear {
		return rounded.AddDate(add, 0, 0)
	}
	return date
}
