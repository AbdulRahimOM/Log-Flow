package helper

import (
	"math"
)

func GetPriorityByFileSize(fileSize int64) uint8 {
	const (
		maxPriority = 10
		minPriority = 1
		minSize     = 10 * 1024              // 1kB
		maxSize     = 1 * 1024 * 1024 * 1024 // 1GB
	)

	var (
		logMinSize = math.Log(float64(minSize))
		logMaxSize = math.Log(float64(maxSize))
	)

	if fileSize <= minSize {
		return maxPriority
	}
	if fileSize >= maxSize {
		return minPriority
	}

	scale := (math.Log(float64(fileSize)) - logMinSize) / (logMaxSize - logMinSize)
	priority := maxPriority - scale*(maxPriority-minPriority)
	return uint8(math.Round(priority))
}
