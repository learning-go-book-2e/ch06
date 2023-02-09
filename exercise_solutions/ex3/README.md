# Exercise 3

## Question
Write a program that builds a `[]Person` with 10,000,000 entries (they could all be the same names and ages). See how long it takes to run. Change the value of `GOGC` and see how that effects the time it takes for the program to complete. Set the environment variable `GODEBUG=gctrace=1` to see when garbage collections happen and see how changing `GOGC` changes the number of garbage collections. What happens if you create the slice with a capacity of 10,000,000?

## Solution

The program itself is very simple, but we can use it to understand a lot about garbage collection and how to optimize performance:

```go
func main() {
	var people []Person
	//people := make([]Person, 0, 10_000_000)
	for i := 0; i < 10_000_000; i++ {
		people = append(people, MakePerson("Fred", "Williamson", 25))
	}
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func MakePerson(firstName, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
```

(All times are from an M1 Apple Silicon Mac with 16GB of RAM)

First we'll see the timings when you don't pre-allocate space in the slice.

```shell
$ time ./ex3
./ex3  1.83s user 0.22s system 197% cpu 1.036 total
$ time GOGC=50 ./ex3
GOGC=50 ./ex3  1.70s user 0.19s system 186% cpu 1.020 total
$ time GOGC=200 ./ex3
GOGC=200 ./ex3  0.97s user 0.24s system 155% cpu 0.779 total
$ time GOGC=1000 ./ex3
GOGC=1000 ./ex3  0.48s user 0.17s system 141% cpu 0.459 total
$ time GOGC=off ./ex3 
GOGC=off ./ex3  0.14s user 0.21s system 91% cpu 0.376 total
```

You can look at the total time see that while there's not much difference between `GOGC=100` (the default) and `GOGC=50`, setting `GOGC=1000` makes the program over twice as fast and turning off GC entirely makes it faster still. Let's re-run with `GODEBUG=gctrace=1` and see what that tells us:

(You can find the definition of the fields printed by `GODEBUG=gctrace=1` at https://pkg.go.dev/runtime#hdr-Environment_Variables)

```shell
$ time GODEBUG=gctrace=1 ./ex3  
gc 1 @0.001s 6%: 0.013+1.7+0.005 ms clock, 0.10+0/1.8/0.051+0.042 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.003s 8%: 0.005+3.3+0.003 ms clock, 0.040+0/3.3/0.013+0.030 ms cpu, 5->5->1 MB, 5 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.007s 10%: 0.008+3.5+0.003 ms clock, 0.069+0/3.5/0.022+0.025 ms cpu, 4->4->2 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.011s 10%: 0.021+3.8+0.002 ms clock, 0.16+0/3.8/0.001+0.022 ms cpu, 5->5->2 MB, 5 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 5 @0.015s 11%: 0.018+4.8+0.003 ms clock, 0.14+0/4.8/0.003+0.024 ms cpu, 6->6->3 MB, 6 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 6 @0.020s 11%: 0.018+5.3+0.001 ms clock, 0.14+0/5.3/0+0.014 ms cpu, 8->8->4 MB, 8 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 7 @0.026s 11%: 0.017+7.0+0.003 ms clock, 0.14+0/7.0/0+0.030 ms cpu, 10->10->5 MB, 10 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 8 @0.033s 11%: 0.018+7.8+0.002 ms clock, 0.14+0/7.8/0+0.023 ms cpu, 12->12->7 MB, 12 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 9 @0.041s 11%: 0.018+7.7+0.003 ms clock, 0.14+0/7.7/0+0.031 ms cpu, 15->15->8 MB, 15 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 10 @0.050s 11%: 0.017+9.9+0.003 ms clock, 0.14+0/9.9/0.002+0.024 ms cpu, 19->19->10 MB, 19 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 11 @0.060s 11%: 0.017+11+0.003 ms clock, 0.14+0/11/0.006+0.024 ms cpu, 24->24->13 MB, 24 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 12 @0.071s 11%: 0.020+11+0.003 ms clock, 0.16+0/11/0+0.028 ms cpu, 30->30->17 MB, 30 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 13 @0.084s 11%: 0.018+13+0.004 ms clock, 0.14+0/13/0+0.033 ms cpu, 38->38->21 MB, 38 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 14 @0.098s 11%: 0.016+18+0.003 ms clock, 0.13+0/18/0+0.024 ms cpu, 47->47->26 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 15 @0.118s 11%: 0.015+20+0.003 ms clock, 0.12+0/20/0+0.028 ms cpu, 59->59->33 MB, 59 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 16 @0.142s 11%: 0.033+26+0.003 ms clock, 0.27+0/26/0.044+0.027 ms cpu, 74->74->41 MB, 74 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 17 @0.169s 11%: 0.015+34+0.003 ms clock, 0.12+0/34/0+0.027 ms cpu, 93->93->51 MB, 93 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 18 @0.206s 12%: 0.018+31+0.005 ms clock, 0.14+0/38/0+0.046 ms cpu, 116->116->64 MB, 116 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 19 @0.243s 12%: 0.017+42+0.003 ms clock, 0.14+0/49/0.005+0.025 ms cpu, 145->145->80 MB, 145 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 20 @0.286s 12%: 0.018+53+0.002 ms clock, 0.14+7.2/61/0.003+0.019 ms cpu, 181->181->101 MB, 182 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 21 @0.343s 13%: 0.014+55+0.003 ms clock, 0.11+5.1/61/23+0.025 ms cpu, 227->227->126 MB, 227 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 22 @0.406s 13%: 0.022+67+0.007 ms clock, 0.18+4.4/72/27+0.060 ms cpu, 284->284->157 MB, 284 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 23 @0.476s 13%: 0.018+95+0.002 ms clock, 0.14+4.3/103/39+0.022 ms cpu, 355->355->197 MB, 355 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 24 @0.578s 13%: 0.016+103+0.003 ms clock, 0.12+7.2/112/43+0.030 ms cpu, 444->444->246 MB, 444 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 25 @0.695s 13%: 0.022+103+6.9 ms clock, 0.18+0/114/59+55 ms cpu, 555->555->308 MB, 555 MB goal, 0 MB stacks, 0 MB globals, 8 P
GODEBUG=gctrace=1 ./ex3  1.69s user 0.18s system 185% cpu 1.013 total
$ time GOGC=50 GODEBUG=gctrace=1 ./ex3 
gc 1 @0.000s 6%: 0.009+0.90+0.005 ms clock, 0.073+0/0.98/0.026+0.041 ms cpu, 2->2->0 MB, 2 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.002s 9%: 0.023+1.6+0.004 ms clock, 0.18+0/1.7/0.021+0.032 ms cpu, 2->2->0 MB, 2 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.004s 10%: 0.022+1.9+0.003 ms clock, 0.17+0/2.0/0.030+0.026 ms cpu, 2->2->1 MB, 2 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.006s 11%: 0.020+2.4+0.002 ms clock, 0.16+0/2.4/0+0.020 ms cpu, 2->2->1 MB, 2 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 5 @0.009s 11%: 0.025+2.9+0.002 ms clock, 0.20+0/3.0/0+0.021 ms cpu, 3->3->1 MB, 3 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 6 @0.012s 11%: 0.023+3.5+0.006 ms clock, 0.19+0/3.5/0+0.054 ms cpu, 4->4->2 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 7 @0.016s 11%: 0.022+3.8+0.002 ms clock, 0.17+0/3.8/0+0.017 ms cpu, 5->5->2 MB, 5 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 8 @0.020s 11%: 0.023+4.3+0.002 ms clock, 0.18+0/4.3/0+0.017 ms cpu, 6->6->3 MB, 6 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 9 @0.024s 11%: 0.021+5.7+0.001 ms clock, 0.16+0/5.7/0+0.015 ms cpu, 8->8->4 MB, 8 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 10 @0.031s 11%: 0.025+5.6+0.003 ms clock, 0.20+0/5.6/0+0.030 ms cpu, 10->10->5 MB, 10 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 11 @0.037s 11%: 0.026+7.0+0.002 ms clock, 0.21+0/7.0/0.026+0.019 ms cpu, 12->12->7 MB, 12 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 12 @0.044s 11%: 0.018+8.6+0.002 ms clock, 0.14+0/8.7/0+0.021 ms cpu, 15->15->8 MB, 15 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 13 @0.054s 11%: 0.022+9.7+0.003 ms clock, 0.17+0/9.7/0+0.024 ms cpu, 19->19->10 MB, 19 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 14 @0.064s 11%: 0.020+9.6+0.006 ms clock, 0.16+0/9.6/0+0.048 ms cpu, 24->24->13 MB, 24 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 15 @0.075s 11%: 0.019+12+0.003 ms clock, 0.15+0/12/0+0.024 ms cpu, 30->30->17 MB, 30 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 16 @0.088s 11%: 0.024+13+0.003 ms clock, 0.19+0/13/0.010+0.027 ms cpu, 38->38->21 MB, 38 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 17 @0.103s 11%: 0.016+16+0.003 ms clock, 0.12+0/16/0+0.028 ms cpu, 47->47->26 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 18 @0.120s 11%: 0.016+22+0.003 ms clock, 0.13+0/22/0+0.027 ms cpu, 59->59->33 MB, 59 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 19 @0.144s 11%: 0.017+25+0.003 ms clock, 0.13+0/25/0+0.024 ms cpu, 74->74->41 MB, 74 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 20 @0.173s 11%: 0.015+31+0.003 ms clock, 0.12+0/31/0+0.024 ms cpu, 93->93->51 MB, 93 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 21 @0.205s 12%: 0.023+36+0.003 ms clock, 0.18+0/43/0.007+0.028 ms cpu, 116->116->64 MB, 116 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 22 @0.248s 12%: 0.020+37+0.003 ms clock, 0.16+2.8/39/13+0.024 ms cpu, 145->145->80 MB, 145 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 23 @0.289s 12%: 0.016+43+0.006 ms clock, 0.12+3.4/46/17+0.055 ms cpu, 181->181->101 MB, 182 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 24 @0.335s 12%: 0.019+63+0.002 ms clock, 0.15+4.6/67/23+0.021 ms cpu, 227->227->126 MB, 227 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 25 @0.405s 12%: 0.018+57+0.26 ms clock, 0.15+0/60/22+2.1 ms cpu, 284->284->157 MB, 284 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 26 @0.472s 12%: 0.028+85+0.003 ms clock, 0.23+7.3/92/34+0.024 ms cpu, 355->355->197 MB, 355 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 27 @0.562s 12%: 0.019+123+0.003 ms clock, 0.15+9.4/132/46+0.026 ms cpu, 444->444->246 MB, 444 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 28 @0.700s 12%: 0.028+129+0.003 ms clock, 0.23+10/140/53+0.028 ms cpu, 555->555->308 MB, 555 MB goal, 0 MB stacks, 0 MB globals, 8 P
GOGC=50 GODEBUG=gctrace=1 ./ex3  1.71s user 0.18s system 186% cpu 1.010 total
$ time GOGC=200 GODEBUG=gctrace=1 ./ex3 
gc 1 @0.002s 4%: 0.008+5.9+0.012 ms clock, 0.068+0/3.4/0.041+0.097 ms cpu, 8->10->4 MB, 8 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.010s 7%: 0.006+7.0+0.003 ms clock, 0.055+0/7.0/0+0.027 ms cpu, 14->14->4 MB, 14 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.018s 9%: 0.007+9.7+0.004 ms clock, 0.056+0/9.7/0.034+0.038 ms cpu, 16->16->6 MB, 16 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 4 @0.030s 9%: 0.030+13+0.004 ms clock, 0.24+0/13/0.090+0.034 ms cpu, 26->26->10 MB, 26 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 5 @0.045s 10%: 0.008+13+0.004 ms clock, 0.066+0/14/3.5+0.032 ms cpu, 41->41->17 MB, 41 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 6 @0.061s 10%: 0.008+21+0.004 ms clock, 0.066+0/21/0+0.039 ms cpu, 64->64->26 MB, 64 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 7 @0.087s 10%: 0.016+21+0.002 ms clock, 0.13+1.5/22/7.3+0.018 ms cpu, 100->100->41 MB, 101 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 8 @0.112s 11%: 0.007+36+0.003 ms clock, 0.058+0/42/5.5+0.026 ms cpu, 157->157->64 MB, 157 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 9 @0.155s 11%: 0.010+50+0.002 ms clock, 0.086+4.1/54/17+0.023 ms cpu, 246->246->101 MB, 246 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 10 @0.216s 12%: 0.009+78+0.002 ms clock, 0.076+4.8/84/30+0.019 ms cpu, 385->385->157 MB, 385 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 11 @0.310s 13%: 0.009+96+5.3 ms clock, 0.076+0/107/56+42 ms cpu, 602->602->246 MB, 602 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 12 @0.439s GOGC=200 GODEBUG=gctrace=1 ./ex3  1.10s user 0.19s system 194% cpu 0.663 total
$ time GOGC=1000 GODEBUG=gctrace=1 ./ex3 
gc 1 @0.011s 5%: 0.036+21+0.043 ms clock, 0.29+0/12/0.91+0.35 ms cpu, 42->53->19 MB, 42 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 2 @0.059s 9%: 0.008+57+5.2 ms clock, 0.068+2.1/32/1.2+41 ms cpu, 223->288->116 MB, 223 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc 3 @0.244s 3%: 0.010+222+0.002 ms clock, 0.087+15/32/269+0.023 ms cpu, 1335->1335->308 MB, 1335 MB goal, 0 MB stacks, 0 MB globals, 8 P
GOGC=1000 GODEBUG=gctrace=1 ./ex3  0.55s user 0.20s system 143% cpu 0.520 total
$ time GOGC=off GODEBUG=gctrace=1 ./ex3 
GOGC=off GODEBUG=gctrace=1 ./ex3  0.14s user 0.23s system 86% cpu 0.424 total
```

You can see the number of GC cycles decrease as GC is assigned a higher value.

Finally, let's recompile the code with the size of the slice pre-allocated to 10,000,000 and see what that does to our garbage collection and our program timing:

```shell
$ time GODEBUG=gctrace=1 ./ex3        
gc 1 @0.004s 18%: 0.012+56+1.8 ms clock, 0.096+0/80/148+14 ms cpu, 381->381->381 MB, 381 MB goal, 0 MB stacks, 0 MB globals, 8 P
GODEBUG=gctrace=1 ./ex3  0.17s user 0.13s system 132% cpu 0.225 total
$ time GOGC=50 GODEBUG=gctrace=1 ./ex3 
gc 1 @0.011s 12%: 0.015+171+0.26 ms clock, 0.12+0/182/106+2.1 ms cpu, 381->381->381 MB, 381 MB goal, 0 MB stacks, 0 MB globals, 8 P
GOGC=50 GODEBUG=gctrace=1 ./ex3  0.34s user 0.10s system 211% cpu 0.208 total
$ time GOGC=200 GODEBUG=gctrace=1 ./ex3 
gc 1 @0.011s 18%: 0.009+64+4.3 ms clock, 0.076+0/88/132+34 ms cpu, 381->381->381 MB, 381 MB goal, 0 MB stacks, 0 MB globals, 8 P
GOGC=200 GODEBUG=gctrace=1 ./ex3  0.18s user 0.12s system 248% cpu 0.121 total
$ time GOGC=1000 GODEBUG=gctrace=1 ./ex3
GOGC=1000 GODEBUG=gctrace=1 ./ex3  0.42s user 0.06s system 166% cpu 0.289 total
$ time GOGC=off GODEBUG=gctrace=1 ./ex3
GOGC=off GODEBUG=gctrace=1 ./ex3  0.05s user 0.04s system 94% cpu 0.099 total
```

The CPU times are much faster. Pre-allocating the slice makes +GOGC=50+ faster than +GOGC=off+ when the slice grows over time. We have, at most, one GC cycle per run. If you know you are going to need a block of memory, it's best to allocate it at once and use it. If you can re-use it, all the better. That's the reason for the slice buffer pattern, also. 