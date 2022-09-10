package parallel

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/spaolacci/murmur3"
)

func TestDo(t *testing.T) {
	type args struct {
		workUnits  int
		workers    int
		bufferSize int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "basic",
			args: args{
				workUnits:  3,
				workers:    2,
				bufferSize: 2,
			},
		},
		{
			name: "basic 200 units",
			args: args{
				workUnits:  200,
				workers:    2,
				bufferSize: 2,
			},
		},
		{
			name: "basic 500 units",
			args: args{
				workUnits:  500,
				workers:    2,
				bufferSize: 2,
			},
		},
		{
			name: "10 workers 500 units",
			args: args{
				workUnits:  500,
				workers:    10,
				bufferSize: 2,
			},
		},
		{
			name: "10 workers 500 units 100 buffer size",
			args: args{
				workUnits:  500,
				workers:    10,
				bufferSize: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// generate work and expected outputs
			work := make([]string, tt.args.workUnits)
			expected := make([]interface{}, tt.args.workUnits)
			for i := 0; i < tt.args.workUnits; i++ {
				w := "work" + strconv.Itoa(i)
				work[i] = w
				o, _ := cpuBoundWorkFunc(w)
				expected[i] = Result[string, uint64]{w, o, nil}
			}

			g := gomega.NewWithT(t)

			got := Do(work, cpuBoundWorkFunc, tt.args.workers, tt.args.bufferSize)

			g.Expect(got).To(gomega.HaveLen(tt.args.workUnits))
			g.Expect(got).To(gomega.ContainElements(expected...))
		})
	}
}

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

func BenchmarkDoCPUBoundWorkV2Part1(b *testing.B) {
	workUnits := 30000
	maxWorkers := 12
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

func BenchmarkDoCPUBoundWorkV2Part2(b *testing.B) {
	workUnits := 30000
	maxWorkers := 24
	bufferSize := 3000
	work := make([]string, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = "work" + strconv.Itoa(i)
	}
	for ws := 13; ws <= maxWorkers; ws++ {
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

func ioBoundWorkFunc(n uint64) (uint64, error) {
	time.Sleep(1 * time.Millisecond)
	return 0, nil
}
