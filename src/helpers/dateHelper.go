package helpers

import (
	"time"
)

func GetDayMonthYearFrom(date string) (string, error) {
	desiredFormat := "02.01.2006"
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", nil
	}

	formattedDate := t.Format(desiredFormat)
	return formattedDate, nil
}
