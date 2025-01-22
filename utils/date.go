package utils

import "time"

var formats = []string{
	"02/01/2006",
	"02/Jan/2006",
	"02/01/2006 15:04:05",
	"02/Jan/2006 15:04:05",
	"2006-01-02T15:04",
}

func FormatDate(dateStr string) string {
	dateObj := GetDateObject(dateStr)

	if dateObj.Hour() == 0 && dateObj.Minute() == 0 {
		return dateObj.Format("02 Jan 2006")
	}

	return dateObj.Format("02 Jan 2006 15:04 MST")
}

func GetDateObject(dateStr string) time.Time {
	var dateObj time.Time
	var err error

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
