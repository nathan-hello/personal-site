package utils

import (
	"strings"
	"time"
)

var formats = []string{
	"2006-01-02T15:04",
	"2006-Jan-02",
	"02-Jan-2006",
	"02 Jan 2006",
	"02 January 2006",
	"2006-01-02",
}

func DateFormatString(dateStr string) string {
	dateObj := DateStringToObject(dateStr)

	if dateObj.Hour() == 0 && dateObj.Minute() == 0 {
		return dateObj.Format("02 Jan 2006")
	}

	return dateObj.Format("02 Jan 2006 15:04 MST")
}

func DateStringToObject(dateStr string) time.Time {
	var dateObj time.Time
	var err error

	dateStr = strings.TrimSpace(dateStr)
	dateStr = strings.Trim(dateStr, "\"")
	dateStr = strings.Trim(dateStr, "\\")
	dateStr = strings.TrimSpace(dateStr)

	for _, format := range formats {
		dateObj, err = time.Parse(format, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		panic(err)
	}

	return dateObj
}
