package usage

import "fmt"

func formatBytes(size uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	value := float64(size)
	unitIndex := 0

	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%d %s", size, units[unitIndex])
	}

	if value >= 100 {
		return fmt.Sprintf("%.0f %s", value, units[unitIndex])
	}
	if value >= 10 {
		return fmt.Sprintf("%.1f %s", value, units[unitIndex])
	}

	return fmt.Sprintf("%.2f %s", value, units[unitIndex])
}
