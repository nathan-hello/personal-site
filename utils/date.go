package utils

import "time"

func FormatDate(dateStr string) string {
	formats := []string{
		"02/01/2006",    
		"02/Jan/2006",   
		"02/01/2006 15:04:05",  
		"02/Jan/2006 15:04:05", 
	}

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

	if dateObj.Hour() == 0 && dateObj.Minute() == 0 {
		return dateObj.Format("02/01/2006")
	}

	return dateObj.Format("02/01/2006 15:04:05")
}
