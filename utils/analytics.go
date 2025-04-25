package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func HttpAnalytic(now time.Time, remoteAddr string, statusCode int, method string, path string, startTime time.Time, json string) {
	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(Env().LOG_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	json = strings.ReplaceAll(json, " ", "\\x20")

	logLine := fmt.Sprintf("%s %s %d %s %s %s %s\n",
		now.Format(time.DateTime),
		remoteAddr,
		statusCode,
		method,
		path,
		time.Since(now),
		json,
	)

	fmt.Print(logLine)

	if _, err := file.WriteString(logLine); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
