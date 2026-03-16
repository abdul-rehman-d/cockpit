package usage

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type CPUSample struct {
	idle  uint64
	total uint64
}

type CPUCollector struct {
	prev CPUSample
}

func NewCPUCollector() *CPUCollector {
	now, _ := ReadCPUSample()
	return &CPUCollector{
		prev: now,
	}
}

func (c *CPUCollector) Key() string {
	return "cpu"
}

func (c *CPUCollector) Sample() (Sample, error) {
	sample, err := ReadCPUSample()
	if err != nil {
		return Sample{
			Key:          c.Key(),
			Name:         "CPU Usage",
			Value:        0,
			ValueInWords: err.Error(),
		}, nil
	}
	if c.prev.total == 0 || sample.total <= c.prev.total || sample.idle < c.prev.idle {
		return Sample{
			Key:          c.Key(),
			Name:         "CPU Usage",
			Value:        0,
			ValueInWords: "cpu warming up",
		}, nil
	}
	totalDelta := float64(sample.total - c.prev.total)
	idleDelta := float64(sample.idle - c.prev.idle)
	usage := (1.0 - idleDelta/totalDelta) * 100
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	c.prev = sample
	return Sample{
		Key:          c.Key(),
		Name:         "CPU Usage",
		Value:        usage,
		ValueInWords: strconv.FormatFloat(usage, 'f', 2, 64) + "%",
	}, nil
}

func ReadCPUSample() (CPUSample, error) {
	switch runtime.GOOS {
	case "linux":
		return readCPUSampleLinux()
	case "darwin":
		return readCPUSampleDarwin()
	default:
		return CPUSample{}, errors.New("cpu unavailable")
	}
}

func readCPUSampleLinux() (CPUSample, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return CPUSample{}, errors.New("cpu unavailable")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return CPUSample{}, errors.New("cpu unavailable")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return CPUSample{}, errors.New("cpu unavailable")
	}

	vals := make([]uint64, 0, len(fields)-1)
	for _, f := range fields[1:] {
		v, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return CPUSample{}, errors.New("cpu unavailable")
		}
		vals = append(vals, v)
	}

	idle := vals[3]
	if len(vals) > 4 {
		idle += vals[4]
	}
	total := uint64(0)
	for _, v := range vals {
		total += v
	}
	if total == 0 {
		return CPUSample{}, errors.New("cpu unavailable")
	}
	return CPUSample{idle: idle, total: total}, nil
}

func readCPUSampleDarwin() (CPUSample, error) {
	out, err := exec.Command("ps", "-A", "-o", "%cpu").Output()
	if err != nil {
		return CPUSample{}, errors.New("cpu unavailable")
	}

	lines := strings.Split(string(out), "\n")
	sum := 0.0
	for i, line := range lines {
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		sum += v
	}
	if sum < 0 {
		sum = 0
	}
	cores := float64(runtime.NumCPU())
	if cores < 1 {
		cores = 1
	}
	usage := sum / cores
	if usage > 100 {
		usage = 100
	}

	const scale = 1000.0
	total := uint64(100 * scale)
	idle := uint64((100 - usage) * scale)
	return CPUSample{idle: idle, total: total}, nil
}
