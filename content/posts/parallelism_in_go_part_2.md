---
title: "Parallelism in Go - Part 2"
date: 2022-09-04T09:39:01+10:00
draft: false
---

Hi there! This is a part twoer to my first post on Parallelism in Go.
In my last post I explored how goroutine workers can be used to complete IO bound work. The optimum number of goroutines to use is always aligned with the number of units of IO bound work. In this post I set out to explore another type of work which goroutine workers can be used to complete: CPU bound work.


## Theory

CPU bound work is work which uses the CPU mainly. Examples of this are adding two numbers, iterating in a loop or hashing a value. No network call, system call or disk read is being made and the CPU is being fully utilized. Its also important to note that CPU bound work does not include synchronization with other goroutines (by way of channels or mutexes).

Unlike IO bound work which waits on something else to complete, CPU bound work consumes CPU cycles. In the IO bound examples we were able to run 3000 goroutines to improve the performance of a workload with 3000 units of IO bound work. The same will not true for CPU bound work. In CPU bound work, we are limited by the number of cores on the machine. If our machine has 4 cores, it can do 4 jobs at once. If our CPU has hyperthreading enabled, then those cores can be kept twice as busy which means we have 8 virtual cores. This means it can do 8 jobs at once. From now on we will refer to the number of virtual cores as logical CPUs.

If we think about a normal application running on the operating system, the "jobs" we are referring to here are threads. Threads are "execution contexts" that are run in parallel by the CPU. However how does this work in go applications? Go maintains a pool of threads and schedules go routines on top of them. How many threads? It will spin up the same number of goroutines as the number of logical CPUs on your system. This is because spinning up more threads won't achieve any greater parallelism, since the CPU is limited by the number of logical CPUs it has. In fact, adding more threads would actually incur an additional penalty due to context switching. In a go application the threads in the thread pool are referred to as logical processors, and we say that goroutines are "scheduled on and off" these logical processors.

When thinking about go applications, goroutines are the concurrency primitive and we don't think about threads. But do the number of goroutines used for parallelising CPU bound work matter? Unlike threads, goroutines do not incur a large penalty when context switched off a logical processor, thanks to their light design. However it should still hold true, that spinning up more goroutines than logical processors, won't achieve more parallelism. We are still limited by the number of logical CPUs available.

## Hypothesis

With this in mind, I set out to demonstrate this effect with a few benchmarks. I created some cpu bound work, parallelized it using the `Do` function of the last post, and increased the goroutine workers each time to see how performance was affected. I hypothesized that:
1. as workers increased to the number of logical CPUs, performance should improve.
2. after increasing workers beyond the number of logical CPUs, performance should get worse.

The reason why I believed that performance would get worse after increasing workers beyond logical CPUs, was because even though go routines are light, spawning them still has a cost. They have to be created and scheduled onto logical processors. Since there wouldn't be any more parallel yield after passing the number of logical CPUs, I theorized the extra goroutines would just add additional cost, therefore resulting in poorer execution time.


This seemed pretty straightforward to me, but after benchmarking I found some pretty weird results...

## Attempt 1

For my first attempt I set up some cpu bound work which hashed an input. I chose hashing because it is a CPU intense operation. Ironically, I chose a fast hash algorithm which tend to be lighter on CPU! In order to generate some steam, I hashed the input 10000 times on itself. 

Similar to the benchmarking code in the previous post, the code first creates some "work", then calls the `.Do` function and passes the work to it. The test code runs several times, on each iteration increasing the number of goroutine workers used by the `.Do` function. Within the test itself, `.Do` is wrapped in a for loop bound by `b.N`. This is standard practice when benchmarking in go. `b.N` is controlled by the benchmark runner and is used to scale the work and take different benchmarks. The benchmark runner later takes the average of all the results.

```go
func BenchmarkDoCPUBoundWork(b *testing.B) {
	workUnits := 3000
	maxWorkers := 24
	bufferSize := 3000

	// create 3000 units of work
	work := make([]string, workUnits)
	for i := 0; i < workUnits; i++ {
		work[i] = "work" + strconv.Itoa(i)
	}

	// run a benchmark, increase workers by 1 and run the next benchmark...
	// repeat until we reach 24 workers
	for ws := 1; ws <= maxWorkers; ws++ {
		b.Run(fmt.Sprintf("workers %d", ws), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Do(work, cpuBoundWorkFunc, ws, bufferSize)
			}
		})
	}
}

// cpuBoundWorkFunc - does some heavy duty cpu bound work
// hashes the input 10000 times on itself and returns the result
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
```

For this attempt I ran the the benchmarks in 3 environments to get some comparison:
- WSL (Windows Subsystem Linux) - 6 cores
- Windows - 6 cores
- Mac - 8 cores

My expectation was that performance would increase as you increased workers up to the number of logical CPUs (which on the Mac was 16 and on the Windows and WSL was 12). I was then hoping to see performance decrease after that number.

I should also note at this point, that I was paying close attention to [Dave Cheney's guide on benchmarks](https://dave.cheney.net/tag/benchmarking). He always suggests to run multiple instances of your benchmarks and then average them out. He also advises that you run some heavy benchmark prior to the real one because CPUs are sometimes lazy and don't perform until you give them a really hefty workload. I did both of these things when gathering my results. Every benchmark run was 3 times before moving onto the next number of workers. Additionally I ran the entire benchmark command 3 times. In the end I had 9 results for each number of workers that I wanted to test.

### Results

![](/images/cpubound1.png)

This graph showed what I was expecting to see. On the y-axis is execution time and on the x-axis is the number of goroutine workers used by the `Do` function to complete the work. In each test the amount of work completed was the same. As goroutine workers increased to the number of logical CPUs on each environment (windows & wsl = 12, mac = 16), execution time decreased.

![](/images/cpubound2.png)

However upon zooming into the above graph, something weird was happening. Instead of performance getting worse after passing logical CPUs for each environment, performance improved! This was not what expected at all. I hypothesized that increasing goroutines beyond logical CPUs would incur a performance penalty!

## Attempt 2

I wanted to be sure that I wasn't seeing these results due to the "lazy CPU" effect mentioned earlier. I theorized that perhaps the test wasn't really pushing the CPU to its limits and only when more work was created (by increasing the goroutines), the CPU actually decided to work harder. To ensure this wasn't the case, I increased the work units and ran the test again, hoping to see my original hypothesis which was that performance should not get better once the workers surpassed the number of logical CPUs.

```diff
func BenchmarkDoCPUBoundWorkV2(b *testing.B) {
-   workUnits := 3000
+   workUnits := 30000
	...
	// all the same
}
```

### Results

I ran the tests again but this time just on my Mac and WSL environments:

![](/images/cpubound3.png)
![](/images/cpubound4.png)

The results were better. This time I saw that the execution time for the Mac tests didn't increase after 8 workers. However I was still puzzled why performance was getting better after surpassing the number of logical CPUs. Specifically the workers = 24 result was 240ms faster than than the workers = 12 result (1.12% speed up).

### CPU & Memory Profiles

At this point I thought it might be a good idea to look at CPU and Memory profiles.

```sh
go test -cpuprofile=wsl-cpu-12.out -run=^$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_12$ .
go test -benchmem -memprofile=wsl-mem-12.out -run=^$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_12$ .
go test -cpuprofile=wsl-cpu-24.out -run=^$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_24$ .
go test -benchmem -memprofile=wsl-mem-24.out -run=^$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_24$ .
```

CPU traces revealed that both processes were almost identical. Some extra time was being spent in some abstract runtime functions. I am not smart enough to work out why. I looked at the memory traces and noticed a huge amount of memory being allocated on the heap throughout the test. This was not at all surprising, after all I was trying to generate steam. But it got me thinking... my test could just be too noisy for me to really see the result I was after. It was just possible that I was testing too many mechanics at once that were all interfering with each other.


I did end up on a lovely aside studying the [Go GC](https://tip.golang.org/doc/gc-guide) though. It didn't help me pinpoint the exact cause but it left we with a sense that "the GC is a complex beast" and that I needed to remove it as much as possible if I really wanted to see the result I was after.

## Attempt 3

### Escape analysis

For my third attempt I sought to avoid all heap allocations in the work function.


```go
/*1*/ func cpuBoundWorkFunc(input string) (uint64, error) {
/*2*/ 	const hashLoopCt = 10000
/*3*/ 	h := murmur3.New64()
/*4*/ 	buf := []byte(input)
/*5*/ 	var out uint64
/*6*/ 	for i := 0; i < hashLoopCt; i++ {
/*7*/ 		h.Write(buf)
/*8*/ 		out = h.Sum64()
/*9*/ 		buf = []byte(strconv.FormatUint(out, 10))
/*10*/ 	}
/*11*/ 	return out, nil
/*12*/ }
```

Without much analysis, there were some obvious ones:
1. For every piece of work, a byte array was being allocated on line 4.
2. For every iteration of the loop, `strconv.FormatUint` was allocating a string on line 9

I changed the above function to this:

```go
/*1*/ func cpuBoundWorkFunc(input uint64) (uint64, error) {
/*2*/ 	const hashLoopCt = 10000
/*3*/ 	barr := make([]byte, 8)
/*4*/ 	for i := 0; i < hashLoopCt; i++ {
/*5*/ 		binary.LittleEndian.PutUint64(barr, input)
/*6*/ 		input = murmur3.Sum64(barr)
/*7*/ 	}
/*8*/ 	return input, nil
/*9*/ }
```

Although I hadn't eliminated the allocation per work (line 3), I wasn't too concerned because this allocation looked like it wouldn't "escape to the heap". Instead this change avoided an allocation on each iteration of the loop by encoding the hash output into an existing byte array rather than creating a new one.

To be sure I ran some escape analysis using the go compiler:

```sh
go test -gcflags=-m=3 -c ./playground/parallelism_in_go/parallel_cpu_bound_2_test.go
```

Surprisingly the escape analysis showed that `barr` was escaping to the heap! But how? My understanding was that `barr` what temporarily be read and appended to the internal hash buffer, long term references would not be needed! It seemed that `murmur3.Sum64` was preventing this value from being allocated to the stack. I did some escape analysis on murmur3 to understand why:

```sh
go build -gcflags=-m=3 github.com/spaolacci/murmur3
```

Why? At a certain point in the murmur3 hash code, a slice of the byte slice was being taken `barr[:n-1]`. Unfortunately in this situation, the go compiler cannot predict how the memory is going to be used so to be on the safe side, it allocates the memory to the heap :(. Consequently my new work function was still allocating an 8 byte array to the heap for every unit of work.

### Getting rid of heap allocations

I really wanted to be rid of allocations so I came up with a solution by creating a variation of my original `Do` function. I created `DoWithState` which had an extra parameter for a function that created shared state for each goroutine worker. It was to be invoked when the goroutine worker was first spawned. I used this to create my 8 byte array once per goroutine worker and therefore avoid a heap allocation for every piece of work processed.

```go
// parallel.go

func DoWithState[I any, O any, S any](work []I, stateFunc StateFunc[S], workFunc WorkFuncWithState[S, I, O], workers int, bufferSize int) []Result[I, O] {
	...
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			s := stateFunc() // create a shared state when the worker goroutine starts.
			defer wg.Done()
			for w := range workC {
				o, err := workFunc(s, w) // pass that shared state to the work function
				resultC <- Result[I, O]{w, o, err}
			}
		}()
	}
}

// parallel_cpu_bound_3_test.go
func cpuBoundWorkFuncV3State() []byte {
	return make([]byte, 8)
}

func cpuBoundWorkFuncV3(barr []byte, input uint64) (uint64, error) {
	const hashLoopCt = 10000
	for i := 0; i < hashLoopCt; i++ {
		binary.LittleEndian.PutUint64(barr, input)
		input = murmur3.Sum64(barr)
	}
	return input, nil
}
```

### Results

With this new work function, I reran my test and here were my results *(this time just on my WSL environment because I was getting lazy)*:

![](/images/cpubound5.png)

![](/images/cpubound6.png)

There was a strange increase in execution at workers = 2 which seemed odd. However apart from this my test results resembled the previous ones in many ways. Just like the previous results, execution trended down while workers was lower than logical CPUs. As workers passed the number of logical CPUs, execution time continued to trend downwards. I was fairly puzzled by the result so I studied execution and memory profiles again. Just like last time, nothing stood out except that more time was spent in runtime functions when workers was lower. It did occur to me though, that although my test might be free from heap allocations, murmur3 was quite a complex function itself and could be introducing noise too. Perhaps murmur3 was introducing complex mechanics that were all interfering with each other at the same time. This led me to one final test where I sought to eliminate all noise!

## Attempt 4

For my final attempt I decided to eliminate murmur3 completely. This way, there would be no unexpected side effects created from the complexities within the hash function. I created a simple work function whose job was to count to 1 million!

```go
func cpuBoundWorkFuncV4(input uint64) (uint64, error) {
	var i uint64
	var x uint64
	for ; i < 1000000; i++ {
		x = i % 2
	}
	return x, nil
}
```

Why the `x = i % 2`? I wanted to create a little more work than just iterating in a loop but I couldn't simply divide two constant numbers together. The go compiler is notoriously good at optimising functions like this and sometimes will even calculate the result for you. To allude the go compiler, I made the result of `x` dependent on `i` so it was hard to optimize!

### Results

I ran the test and these were my results:

![](/images/cpubound9.png)
![](/images/cpubound10.png)
![](/images/cpubound11.png)

Like before, execution time decreased as workers approach logical CPUs, and some higher values of workers out performed the result when workers was equal to logical CPUs. However in this test I noticed that execution time was trending up after workers surpassed logical CPUs (as seen in the last graph with up to 100 workers). This led me to believe I was onto something.


So I ran it up to 180 workers:

![](/images/cpubound12.png)

And then all the way to 1000 workers:

![](/images/cpubound13.png)

Finally I could see what I was looking for! Plotting the trendline revealed a 0.0005ms (500ns) increase in execution per goroutine added. Whether this was from scheduling costs or just the cost to boot the goroutine, I don't know. What we can deduce is that this is an extremely small and therefore insignificant penalty for spinning up a goroutine. I was admitedly expecting something much worse but in the end only found a meer 500ns penalty!

## Lessons Learnt

1. Execution time improves as your goroutine workers approach the available number of logical CPUs. This make sense because your workload can be parallelized onto all the available cores. 
2. There is a penalty for creating more go routines than logical CPUs BUT the penalty is insignificant and you shouldn't worry about it. 500ns is nothing.
3. Function complexity contributes to noise in an application which makes it worthless to try and isolate factors like the number of goroutines.


This last point rings home for me after a few years of working on production issues and optimising code. Often you will be asked to "make service x" faster. Sometimes you find that certain operations are slow, or the GC is thrashing your application. You can try to tweaking the GC settings and/or allocating more resources to solve these issues. However at a certain point this approach won't work any more. Instead you have to go the source: reduce allocations and optimize functions that burn lots of CPU. This whole exercise felt like a mirror of this truth.

From personal experience I have also learned that optimising code comes with trade offs. Optimised code will generally make it...
1. harder to train new developers
2. harder to follow changing business requirements
3. harder to debug/maintain

Optimisation can and should be done, but it comes with a price. Its all about judging whether that price is worth it.

## Conclusion

For CPU bound work, spawning as many goroutines as there are logical CPUs will allow you to use all the available resources. `runtime.NumCPU()` can be used find out the number of logical CPUs at runtime. If you happen to go beyond that, it doesn't really matter :)


