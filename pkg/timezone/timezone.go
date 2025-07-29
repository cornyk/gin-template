package timezone

import (
	"time"
)

func InitTimezone(location string) {
	if location == "" {
		location = "UTC"
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		panic("Cannot loading timezone: " + err.Error())
	}
	time.Local = loc
}
