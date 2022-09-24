package parallel

import (
	"strconv"
	"testing"

	"github.com/onsi/gomega"
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

func TestDoWithState(t *testing.T) {
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
			work := make([]uint64, tt.args.workUnits)
			expected := make([]interface{}, tt.args.workUnits)
			byteArray := cpuBoundWorkFuncV3State()
			for i := 0; i < tt.args.workUnits; i++ {
				w := uint64(i)
				work[i] = w
				o, _ := cpuBoundWorkFuncV3(byteArray, w)
				expected[i] = Result[uint64, uint64]{w, o, nil}
			}

			g := gomega.NewWithT(t)

			got := DoWithState(work, cpuBoundWorkFuncV3State, cpuBoundWorkFuncV3, tt.args.workers, tt.args.bufferSize)

			g.Expect(got).To(gomega.HaveLen(tt.args.workUnits))
			g.Expect(got).To(gomega.ContainElements(expected...))
		})
	}
}
