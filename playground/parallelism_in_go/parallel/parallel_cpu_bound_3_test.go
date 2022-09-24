package parallel

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/spaolacci/murmur3"
)

func BenchmarkDoCPUBoundWorkV3(b *testing.B) {
	workUnits := 300
	maxWorkers := 24
	bufferSize := 3000
	work := make([]uint64, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = uint64(i)
	}
	for ws := 1; ws <= maxWorkers; ws++ {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DoWithState(work, cpuBoundWorkFuncV3State, cpuBoundWorkFuncV3, ws, bufferSize)
			}
		})
	}
}

func BenchmarkDoCPUBoundWorkV3Intense(b *testing.B) {
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
				DoWithState(work, cpuBoundWorkFuncV3State, cpuBoundWorkFuncV3, ws, bufferSize)
			}
		})
	}
}

func BenchmarkDoCPUBoundWorkV3Backwards(b *testing.B) {
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
				DoWithState(work, cpuBoundWorkFuncV3State, cpuBoundWorkFuncV3, ws, bufferSize)
			}
		})
	}
}

func cpuBoundWorkFuncV3State() []byte {
	return make([]byte, 8)
}

func cpuBoundWorkFuncV3(byteArray []byte, input uint64) (uint64, error) {
	const hashLoopCt = 10000
	for i := 0; i < hashLoopCt; i++ {
		binary.LittleEndian.PutUint64(byteArray, input)
		input = murmur3.Sum64(byteArray)
	}
	return input, nil
}
