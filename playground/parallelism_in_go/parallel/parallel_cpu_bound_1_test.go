package parallel

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/spaolacci/murmur3"
)

func BenchmarkDoCPUBoundWork(b *testing.B) {
	workUnits := 3000
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

func BenchmarkDoCPUBoundWorkMoreWorkers(b *testing.B) {
	workUnits := 3000
	maxWorkers := 6000
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

func cpuBoundWorkFunc(input string) (uint64, error) {
	const hashLoopCt = 10000
	h := murmur3.New64()
	buf := []byte(input)
	var out uint64
	for i := 0; i < hashLoopCt; i++ {
		h.Write(buf)
		out = h.Sum64()
		buf = []byte(strconv.FormatUint(out, 10))
	}
	return out, nil
}
