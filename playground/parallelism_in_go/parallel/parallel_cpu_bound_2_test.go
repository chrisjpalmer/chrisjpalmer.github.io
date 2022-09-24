package parallel

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkDoCPUBoundWorkV2(b *testing.B) {
	workUnits := 30000
	maxWorkers := 24
	bufferSize := 3000
	work := make([]string, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = "work" + strconv.Itoa(i)
	}
	for ws := 1; ws <= maxWorkers; ws++ {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFunc, ws, bufferSize)
			}
		})
	}
}

func BenchmarkDoCPUBoundWorkMoreWorkersV2(b *testing.B) {
	workUnits := 30000
	maxWorkers := 1000
	bufferSize := 3000
	work := make([]string, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = "work" + strconv.Itoa(i)
	}
	for ws := 0; ws <= maxWorkers; ws += 100 {
		_ws := ws
		if _ws == 0 {
			_ws = 1
		}
		b.Run(fmt.Sprintf("workers %d", _ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFunc, _ws, bufferSize)
			}
		})
	}
}
