test:
	go test ./...

escape-analysis-murmur3:
	 go build -gcflags=-m=3 github.com/spaolacci/murmur3

escape-analysis-test:
	go test -gcflags=-m=3 -c ./playground/parallelism_in_go/parallel

bench:
	cd parallel && go test -run=^$$ -bench .

## IOBOUND
bench-iobound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoIOBoundWork$$ .

## CPUBOUND V1
bench-cpubound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWork$$ .

bench-cpubound-moreworkers:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkers$$ .

## CPUBOUND V2
bench-cpubound-moreworkersv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .

benchv2:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .

bench-cpuboundv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .

bench-cpuboundv2-count3:
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .	
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .	
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .		

bench-cpuboundv2-cpumem-wsl:
	# for wsl
	cd parallel && go test -cpuprofile=v2-wsl-cpu-8.out -benchmem -memprofile=v2-wsl-mem-8.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_8$$ .
	cd parallel && go test -cpuprofile=v2-wsl-cpu-12.out -benchmem -memprofile=v2-wsl-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_12$$ .
	cd parallel && go test -cpuprofile=v2-wsl-cpu-24.out -benchmem -memprofile=v2-wsl-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_24$$ .
	 
bench-cpuboundv2-cpumem-mac:
	cd parallel && go test -cpuprofile=v2-mac-cpu-12.out -benchmem -memprofile=v2-mac-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_12$$ .
	cd parallel && go test -cpuprofile=v2-mac-cpu-16.out -benchmem -memprofile=v2-mac-mem-16.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_16$$ .
	cd parallel && go test -cpuprofile=v2-mac-cpu-24.out -benchmem -memprofile=v2-mac-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_24$$ .

bench-cpuboundv2-trace-wsl:
	# for wsl
	cd parallel && go test -trace=v2-wsl-trace-8.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_8$$ .
	cd parallel && go test -trace=v2-wsl-trace-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_12$$ .
	cd parallel && go test -trace=v2-wsl-trace-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_24$$ .
	 
bench-cpuboundv2-trace-mac:
	cd parallel && go test -trace=v2-mac-trace-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_12$$ .
	cd parallel && go test -trace=v2-mac-trace-16.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_16$$ .
	cd parallel && go test -trace=v2-mac-trace-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$/workers_24$$ .

## CPUBOUND V3
bench-cpuboundv3-cpumemtrace-wsl:
	# warm up
	cd parallel && go test -count=5 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3$$/workers_24$$ .

	# run test
	cd parallel && go test -benchtime=1x -cpuprofile=v3-wsl-cpu-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_24$$ .
	cd parallel && go test -benchtime=1x -cpuprofile=v3-wsl-cpu-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_12$$ .
	cd parallel && go test -benchtime=1x -benchmem -memprofile=v3-wsl-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_24$$ .
	cd parallel && go test -benchtime=1x -benchmem -memprofile=v3-wsl-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_12$$ .
	cd parallel && go test -benchtime=1x -trace=v3-wsl-trace-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_24$$ .
	cd parallel && go test -benchtime=1x -trace=v3-wsl-trace-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Intense/workers_12$$ .

bench-cpuboundv3-count3:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3$$ .

bench-cpuboundv3-backwards-count3:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV3Backwards$$ .

## CPUBOUND V4
bench-cpuboundv4-count3:
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV4$$ .
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV4$$ .
	cd parallel && go test -timeout=2h -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV4$$ .

bench-cpuboundv4-backwards-count3:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV4Backwards$$ .


# UTILITIES
view-cpu:
	go tool pprof -http=localhost:$$PORT parallel/$$CPU

view-mem:
	go tool pprof -http=localhost:$$PORT parallel/$$MEM

view-trace:
	go tool trace -http=localhost:$$PORT parallel/$$TRACE