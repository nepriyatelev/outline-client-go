package outline

import (
	"fmt"
	"time"
)

func formatDuration(d time.Duration) string {
	// Determine the sign and work with the absolute value
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}

	// Hours
	h := int64(d.Hours())
	if h != 0 {
		return fmt.Sprintf("%s%dh", sign, h)
	}
	// Minutes
	m := int64(d.Minutes())
	if m != 0 {
		return fmt.Sprintf("%s%dm", sign, m)
	}
	// Seconds
	s := int64(d.Seconds())
	if s != 0 {
		return fmt.Sprintf("%s%ds", sign, s)
	}
	// If still zero, ignore the sign and return "0s"
	return "0s"
}
