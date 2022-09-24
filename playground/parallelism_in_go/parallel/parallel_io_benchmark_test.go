package parallel

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkDoIOBoundWork(b *testing.B) {
	workUnits := 3000
	maxWorkers := 6000
	bufferSize := 3000
	work := make([]uint64, workUnits)
	for ws := 0; ws <= maxWorkers; ws += 100 {
		_ws := ws
		if _ws == 0 {
			_ws = 1
		}
		b.Run(fmt.Sprintf("workers %d", _ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, ioBoundWorkFunc, _ws, bufferSize)
			}
		})
	}
}

func ioBoundWorkFunc(n uint64) (uint64, error) {
	time.Sleep(1 * time.Millisecond)
	return 0, nil
}
