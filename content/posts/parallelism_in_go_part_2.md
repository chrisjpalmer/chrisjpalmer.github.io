---
title: "Parallelism in Go - Part 2"
date: 2022-09-04T09:39:01+10:00
draft: false
---

Hi there! This is a part twoer to my first post on Parallelism in Go.
In my last post I explored how goroutine workers can be used to complete IO bound work. The optimum number of goroutines to use is always aligned with the number of units of IO bound work. In this post I set out to explore another time of work which goroutine workers can be used to complete: CPU bound work.


## Hypothesis

CPU bound work is work which uses the CPU mainly. Examples of this are adding two numbers, iterating in a loop or hashing a value. No network call, system call or disk read is being made and the CPU is being fully utilized. Its also important to note that CPU bound work does not include synchronization with other goroutines (by way of channels or mutexes).

Unlike IO bound work which waits on something else to complete, CPU bound work consumes CPU cycles. In the IO bound examples we were able to run 3000 goroutines to improve the performance of a workload with 3000 units of IO bound work. The same is not true for CPU bound work. In CPU bound work, we are limited by the number of cores on the machine. If our machine has 4 cores, it can do 4 jobs at once. If our CPU has hyperthreading enabled, then those cores can be kept twice as busy which means we have 8 virtual cores which means it can do 8 jobs at once. These virtual cores are called logical CPUs for simplicity.

If we think about a normal application running on the operating system, the "jobs" we are referring to here are threads. Threads are "execution contexts" that are run in parallel by the CPU. However how does this work in go applications? Go maintains a pool of threads and schedules go routines on top of them. How many threads? It will spin up the same number of goroutines as the number of logical CPUs on your system. This is because spinning up more threads won't achieve any greater parallelism, since the CPU is limited by the number of logical CPUs it has. In fact, adding more threads would actually incur an additional penalty due to context switching. In a go application the threads in the thread pool are referred to as logical processors, and we say that goroutines are "scheduled on and off" these logical processors.

When thinking about go applications, goroutines are the concurrency primitive and we don't think about threads. But do the number of goroutines used for parallelising CPU bound work matter? Unlike threads, goroutines do not incur a large penalty when context switched off a logical processor, thanks to their light design. However it should still hold true, that spinning up more goroutines than logical processors, won't achieve more parallelism. We are still limited by the number of logical CPUs available.

With this in mind, I set out to demonstrate this effect with a few benchmarks. I created some cpu bound work, parallelized using the `Do` function of the last post, and increased the goroutine workers each time to see how performance was affected. I hypothesized that as you increase workers to the number of logical CPUs, that performance would improve. I also hypothesized that after increasing workers beyond that number, performance would stay the same or eventually get worse. 

This seemed pretty straightforward to me, but after benchmarking I found some pretty weird results...

## Attempt 1

For my first attempt I set up some cpu bound work which hashed an input in a tight loop 10000 times. I chose hashing because it is a CPU intense operation. Since murmur3 is quite performant, I looped it 10000 times to generate some steam. 

Similar to the benchmarking code in the previous post, I first generate some workload, and then call the `.Do` function increasing the number of workers each time. The `Do` function is surrounded by another for loop which takes into account `b.N`. Its important when writing benchmarks to run the target code `b.N` times so that the go benchmarking runtime can control the number of iterations. It does this to take multiple samples and then average out the results.

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

For this attempt I was going to run the benchmarks in 3 environments and compare:
- WSL (Windows Subsystem Linux) - 6 cores
- Windows - 6 cores
- Mac - 8 cores

My expectation was that performance would increase as you increased workers up to the number of logical CPUs (which on the Mac was 16 and on the Windows and WSL was 12). I was then hoping to see performance decrease a little after that number. The reason why I believed that performance would decrease after workers > logical CPUs is because I am aware that even though goroutines are lightweight, they still incur a penalty for being scheduled on and off logical processors. I was expecting to see the execution time of the benchmark jump up as a result of this penalty.

I should also note at this point, that I was paying close attention to Dave Chetney's guide on benchmarks. He always suggests to run multiple instances of your benchmarks and then average them out. He also advises that you run some heavy benchmark prior to the real one because CPUs are sometimes lazy and don't perform until you give them a really hefty workload. I did both of these things when gathering my results. Every benchmark run 3 times before moving onto the next number of workers. Additionally I ran the entire benchmark command 3 times. In the end I had 9 results for each worker. Using some spreadsheeting, I took the averages and graphed the results.

![](/images/cpubound1.png)

![](/images/cpubound2.png)

The results surprised me a little. Yes execution time decreased as you increased goroutine workers. However I was expecting performance to decrease after workers surpassed the number of logical CPUs on each machine. In fact the opposite happened: as workers surpassed logical CPUs performance improved! What's more is that the result for the Mac was showing that the best result was when workers was 8 which was half the number of lofical processors! This was a very weird result and I was puzzled by it.

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

I ran the tests again but this time just on my Mac and WSL environments:

![](/images/cpubound3.png)
![](/images/cpubound4.png)

The results were better. This time I saw that the execution time for the Mac tests didn't increase after 8 workers. However I was still puzzled why performance was getting better after surpassing the number of logical CPUs. Specifically the workers = 24 result was 240ms faster than than the workers = 12 result (1.12% faster).

At this point I thought it might be a good idea to look at CPU and Memory profiles as well as execution traces.

```
go test -cpuprofile=wsl-cpu-12.out -benchmem -memprofile=wsl-mem-12.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_12$$ .
go test -cpuprofile=wsl-cpu-24.out -benchmem -memprofile=wsl-mem-24.out -run=^$$ -bench ^BenchmarkDoCPUBoundWorkV2/workers_24$$ .
```

After comparing several CPU traces I didn't find anything particularly interesting.
However the memory trace was telling me something. There were huge byte array allocations in my code. Over the course of the test the application allocated over 16Gb of memory. Obviously 16Gb wasn't live the whole time (coz it would have been cleaned up by the GC)... but it led me to think my results could be being impacted by the GC.

This led me to another aside on the GC where I ended up reading this [wonderful article](https://tip.golang.org/doc/gc-guide) on how the GC works. I was searching for some information which might support a new running theory I had: increased workers beyond logical processors, although incurring a scheduler penalty, might actually be benefiting the GC in some way. I still can't be sure whether the GC was the problem with these tests but I did note that the GC mark phase cannot complete until a goroutine is put to sleep. In the case where workers = logical processors, all workers were being kept as busy as possible so they probably didn't want to sleep. In cases like this the GC can issue a penalty to those workers by introducing a write barrier OR also making that worker do something called gcAssist. Gc Assist is when a goroutine is forced by the GC to stop what its doing and participate in the mark phase (mark phase is the phase when the GC discovers all the live allocations). It was possible because all my goroutine workers were so busy that they were actually incurring GC penalties and perhaps a larger number of workers made it easier for the GC to mark memory coz they would sleep more regularly.

When I viewed the execution traces, I definitely saw that for workers = 24, go routines were frequently switched off logical processors where as for workers = 12, goroutines could stay on for much longer. Could this be a contributing factor... ? I never dug deep enough to find out. 

However it got me thinking, what if I eliminated the memory factor altogether and made this test as CPU bound as possible.

## Experiment 3

My theory was that either allocations or the GC (or both) were contributing to my tests somehow, so I wanted to build a better work function which minimized their effects. I did some analysis on the work function:


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

I found:
1. on line 9, I was allocating a string on every iteration of the loop
2. on line 4, I was allocating a byte array every time `cpuBoundWorkFunc` was called

Whats more is after doing some escape analysis, I also found that the buffer being passed to `h.Write(buf)` was escaping to the heap:

```sh
go build -gcflags=-m=3 github.com/spaolacci/murmur3
go test -gcflags=-m=3 -c ./playground/parallelism_in_go/parallel_cpu_bound_2_test.go
```

This was apparently happening due to the internals of murmur3. Specifically somewhere internally a range was being taken from this byte slice. This is grounds for Go to allocate the byte array to the heap so it was escaping to the heap.

I wanted to eliminate heap allocations so I dreamed a new `Do` function that could maintain a shared state per worker goroutine:

```go
//parallel.go

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
```

Then I created a new work function and a function to initialize the shared state:

```go
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
```

My new work function would obviously not produce the same output, but it was more or less doing the same thing and this time avoiding allocations on every iteration of the loop. I hoped this would eliminate the noise factor which I hypothesized was coming from the GC.

Here were my results:

![](/images/cpubound5.png)

![](/images/cpubound6.png)

I noticed right away this strange peak at workers = 2 than hadn't cropped up before. I found this weird since it didn't make sense that an increase in parallelism was also increasing execution time! Oh well, it certainly revealed that this work function had very different characteristics to the last one!

Without fail however, I continued to see execution time trend downwards after workers > logical processors! It couldn't be alluded!

Once again I turned to cpu profiles, memory profiles and execution traces.

This time the memory trace was very boring. Hardly any memory was being allocated for the whole test. However when comparing CPU profiles between workers = 12 and workers = 24, I couldn't see any obvious reason why workers = 24 was faster than the other. I was so baffled by this result I started to wonder whether somehow ramping the workers up to 24 was actually "warming" up the CPU better, making the workers = 24 test run faster. To be sure, I actually ran the tests backwards:

```go
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
```

Astonishingly I got the same results. For some unknown reason, workers = 24 just performed better.

I started to look into the CPU profiles of individual lines of the murmur hash function. This is where I got another hunch. It seemed that between 12 and 24 worker tests particular lines of code in murmur3 just performed better for workers = 24. I wondered if this might have something to do with alignment of memory and cache lines. Its impossible to know because I just don't have that much in-depth knowledge. However it gave me one lust hunch... perhaps murmur3 is just such a complex function itself, that its creating situations where for no apparent reason a higher number of workers favour it. This led me to one last test...

### Experiment 4

No more murmur.. this time I made the simplest work function conceivable with 0 memory allocations:

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

This function does nothing more than operate in a tight loop. The only reason I added lines like `x = i % 2` was to prevent the go compiler from optimising out my variables. I thought go probably doesn't know how to optimize the result of `x` in this situation so this would be a good way to make the CPU do some work.

I ran the test and these were my results:

![](/images/cpubound8.png)
![](/images/cpubound9.png)
![](/images/cpubound10.png)

The test results showed something different this time. Instead of the execution time trending down as workers increased beyond 12, it was trending up. Still the optimium number of worker go routines was not 12, however at least I was starting to see that additional go routines was incurring a schedule penalty as predicted.

In fact the results were even more apparent when I ran it all the way up to 100 workers:
![](/images/cpubound11.png)

This result got me pretty excited. So I ran with even more workers and performed 3 test runs to make sure I wasn't running into any noise from the machine:

![](/images/cpubound12.png)

.. and more workers:

![](/images/cpubound13.png)

Okay finally I could see what I was looking for! Taking the trendline each time, I found that there was a 0.0005ms (500ns) scheduler penalty per goroutine.
As for the magic number 12, I could not find it. In repeats of the test it just didn't show up:

![](/images/cpubound14.png)
![](/images/cpubound15.png)
![](/images/cpubound16.png)

However what I did notice is that there seems to be random noise in all the results. For example, comparing 3 individual back to back runs of the test, the results don't exactly align with each other:
![](/images/cpubound17.png)

This suggested to me that perhaps its impossible to see the 500ns penalty between workers = 12 and workers = 13. Perhaps that just wasn't going to possible given that my personal computer always has some degree of noise.

## Conclusion

Although I coulodn't find the magic number workers = logical processors, I can see with a reasonable degree of confidence that workers = logical processors is the point at which you get the most reasonable optimization with CPU bound work.s

