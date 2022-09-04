test:
	go test ./...

bench-cpubound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWork$$ .

bench-cpubound2:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoCPUBoundWork2$$ .

bench-iobound:
	cd parallel && go test -run=^$$ -bench ^BenchmarkDoIOBoundWork$$ .

bench:
	cd parallel && go test -bench .