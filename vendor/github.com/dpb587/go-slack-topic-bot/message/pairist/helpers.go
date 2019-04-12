package pairist

import (
	"time"
)

func WorkingHours(start, stop, tzName string) func() bool {
	tz, err := time.LoadLocation(tzName)
	if err != nil {
		panic(err)
	}

	return func() bool {
		now := time.Now().In(tz)

		dt := now.Format("Mon")

		if dt == "Sat" || dt == "Sun" {
			return false
		}

		ts := now.Format("15:04")

		return ts >= start && ts < stop
	}
}
