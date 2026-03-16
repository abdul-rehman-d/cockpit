package usage

import (
	"errors"
	"fmt"
	"syscall"
)

type StorageCollector struct{}

func NewStorageCollector() *StorageCollector {
	return &StorageCollector{}
}

func (m *StorageCollector) Key() string {
	return "memory"
}

func (m *StorageCollector) Sample() (Sample, error) {
	used, total, err := storageUsage("/")
	if err != nil {
		return Sample{
			Key:          m.Key(),
			Name:         "Storage Usage",
			Value:        0,
			ValueInWords: err.Error(),
		}, nil
	}
	usage := float64(used) / float64(total)
	return Sample{
		Key:          m.Key(),
		Name:         "Storage Usage",
		Value:        usage,
		ValueInWords: fmt.Sprintf("%s used of %s", formatBytes(used), formatBytes(total)),
	}, nil
}

func storageUsage(path string) (uint64, uint64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, errors.New("storage unavailable")
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	if total == 0 || free > total {
		return 0, 0, errors.New("storage unavailable")
	}
	used := total - free
	return used, total, nil
}
