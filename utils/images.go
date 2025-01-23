package utils

import (
	"fmt"
	"math"
	"time"
)

func GetBlogImageDir(d time.Time) string {
	year := d.Year()
	return fmt.Sprintf("./public/images/covers/%d", year)
}

func FormatSize(bytes int64) string {

	if bytes == 0 {
		return "0 Bytes"

	}
	const k = 1024
	const decimals = 2
	sizes := []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	exponent := math.Floor(math.Log(float64(bytes)) / math.Log(k))

	real := float64(bytes) / math.Pow(k, exponent)

	truncated := math.Trunc(real*100) / 100

	return fmt.Sprintf("%d %s", int(truncated), sizes[int(exponent)])

}
