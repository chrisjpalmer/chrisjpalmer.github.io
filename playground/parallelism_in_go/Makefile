test:
	go test ./...

bench-cpubound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWork$$ .

bench-cpubound-moreworkers:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkers$$ .

bench-cpuboundv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .

bench-cpuboundv2-count3:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1$$ .	
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2$$ .	
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1$$ .	
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2$$ .	
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1$$ .	
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2$$ .	

bench-cpuboundv2-cpumem-wsl:
	# for wsl
	cd parallel && go test -cpuprofile=wsl-cpu-8.out -benchmem -memprofile=wsl-mem-8.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_8$$ .
	cd parallel && go test -cpuprofile=wsl-cpu-12.out -benchmem -memprofile=wsl-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_12$$ .
	cd parallel && go test -cpuprofile=wsl-cpu-24.out -benchmem -memprofile=wsl-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_24$$ .
	 
bench-cpuboundv2-cpumem-mac:
	cd parallel && go test -cpuprofile=mac-cpu-12.out -benchmem -memprofile=mac-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_12$$ .
	cd parallel && go test -cpuprofile=mac-cpu-16.out -benchmem -memprofile=mac-mem-16.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_16$$ .
	cd parallel && go test -cpuprofile=mac-cpu-24.out -benchmem -memprofile=mac-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_24$$ .

bench-cpuboundv2-trace-wsl:
	# for wsl
	cd parallel && go test -trace=wsl-trace-8.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_8$$ .
	cd parallel && go test -trace=wsl-trace-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_12$$ .
	cd parallel && go test -trace=wsl-trace-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_24$$ .
	 
bench-cpuboundv2-trace-mac:
	cd parallel && go test -trace=mac-trace-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part1/workers_12$$ .
	cd parallel && go test -trace=mac-trace-16.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_16$$ .
	cd parallel && go test -trace=mac-trace-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2Part2/workers_24$$ .

bench-cpubound-moreworkersv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .

benchv2:
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .
	cd parallel && go test -count=3 -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .
	
bench-iobound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoIOBoundWork$$ .

bench:
	cd parallel && go test -run=^$$ -bench .

view-cpu:
	go tool pprof -http=localhost:$$PORT parallel/$$CPU

view-mem:
	go tool pprof -http=localhost:$$PORT parallel/$$MEM

view-trace:
	go tool trace -http=localhost:$$PORT parallel/$$TRACE