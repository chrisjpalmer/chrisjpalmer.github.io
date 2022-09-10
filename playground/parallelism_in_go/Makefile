test:
	go test ./...

bench-cpubound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWork$$ .

bench-cpubound-moreworkers:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkers$$ .

bench-cpuboundv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2$$ .

bench-cpubound-moreworkersv2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWorkMoreWorkersV2$$ .

bench-iobound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoIOBoundWork$$ .

bench:
	cd parallel && go test -run=^$$ -bench .