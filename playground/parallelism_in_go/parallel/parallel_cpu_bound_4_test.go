package parallel

import (
	"fmt"
	"testing"
)

func BenchmarkDoCPUBoundWorkV4(b *testing.B) {
	workUnits := 300
	maxWorkers := 1000
	bufferSize := 3000
	work := make([]uint64, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = uint64(i)
	}
	for ws := 1; ws <= maxWorkers; ws++ {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFuncV4, ws, bufferSize)
			}
		})
	}
}

func BenchmarkDoCPUBoundWorkV4Intense(b *testing.B) {
	workUnits := 300000
	maxWorkers := 24
	bufferSize := 3000
	work := make([]uint64, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = uint64(i)
	}
	for ws := 1; ws <= maxWorkers; ws++ {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFuncV4, ws, bufferSize)
			}
		})
	}
}

func BenchmarkDoCPUBoundWorkV4Backwards(b *testing.B) {
	workUnits := 300
	maxWorkers := 24
	bufferSize := 3000
	work := make([]uint64, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = uint64(i)
	}
	for ws := maxWorkers; ws > 0; ws-- {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFuncV4, ws, bufferSize)
			}
		})
	}
}

func cpuBoundWorkFuncV4(input uint64) (uint64, error) {
	var i uint64
	var x uint64
	for ; i < 1000000; i++ {
		x = i % 2
	}
	return x, nil
}
