package usage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type MemoryCollector struct{}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (m *MemoryCollector) Key() string {
	return "memory"
}

func (m *MemoryCollector) Sample() (Sample, error) {
	used, total, err := memoryUsage()
	if err != nil {
		return Sample{
			Key:          m.Key(),
			Name:         "Memory Usage",
			Value:        0,
			ValueInWords: err.Error(),
		}, nil
	}
	usage := (float64(used) / float64(total)) * 100
	return Sample{
		Key:          m.Key(),
		Name:         "Memory Usage",
		Value:        usage,
		ValueInWords: fmt.Sprintf("%d / %d", used, total),
	}, nil
}

func memoryUsage() (uint64, uint64, error) {
	switch runtime.GOOS {
	case "linux":
		return memoryUsageLinux()
	case "darwin":
		return memoryUsageDarwin()
	default:
		return 0, 0, errors.New("memory unavailable")
	}
}

func memoryUsageLinux() (uint64, uint64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, errors.New("memory unavailable")
	}
	defer f.Close()

	var totalKB uint64
	var availKB uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			totalKB = parseKB(line)
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			availKB = parseKB(line)
		}
	}
	if totalKB == 0 || availKB > totalKB {
		return 0, 0, errors.New("memory unavailable")
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, errors.New("memory unavailable")
	}
	total := totalKB * 1024
	used := (totalKB - availKB) * 1024
	return used, total, nil
}

func memoryUsageDarwin() (uint64, uint64, error) {
	totalOut, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return 0, 0, errors.New("memory unavailable")
	}
	pageOut, err := exec.Command("sysctl", "-n", "hw.pagesize").Output()
	if err != nil {
		return 0, 0, errors.New("memory unavailable")
	}
	vmOut, err := exec.Command("vm_stat").Output()
	if err != nil {
		return 0, 0, errors.New("memory unavailable")
	}

	total, err := strconv.ParseUint(strings.TrimSpace(string(totalOut)), 10, 64)
	if err != nil || total == 0 {
		return 0, 0, errors.New("memory unavailable")
	}
	pageSize, err := strconv.ParseUint(strings.TrimSpace(string(pageOut)), 10, 64)
	if err != nil || pageSize == 0 {
		return 0, 0, errors.New("memory unavailable")
	}

	active := uint64(0)
	wired := uint64(0)
	compressor := uint64(0)
	for _, line := range strings.Split(string(vmOut), "\n") {
		switch {
		case strings.HasPrefix(line, "Pages active"):
			active = parseDarwinPages(line)
		case strings.HasPrefix(line, "Pages wired down"):
			wired = parseDarwinPages(line)
		case strings.HasPrefix(line, "Pages occupied by compressor"):
			compressor = parseDarwinPages(line)
		}
	}
	usedPages := active + wired + compressor
	used := usedPages * pageSize
	if used > total {
		used = total
	}
	return used, total, nil
}

func parseKB(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	v, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func parseDarwinPages(line string) uint64 {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0
	}
	v := strings.TrimSpace(parts[1])
	v = strings.TrimSuffix(v, ".")
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
