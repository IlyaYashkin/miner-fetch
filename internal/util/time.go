package util

import (
	"fmt"
	"time"
)

func FormatDuration(seconds int) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute

	return fmt.Sprintf("%d д. %d ч. %d м.", days, hours, minutes)
}
