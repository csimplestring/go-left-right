

# Go Left Right Concurrency
A Go implementation of the left-right concurrency control algorithm in paper *<Left-Right - A Concurrency Control Technique with Wait-Free Population Oblivious Reads>*

This library provides a concurrency primitive for high concurrency reads over a single-writer data structure. The micro benchmark shows the left-right pattern is 2-3x faster than the RWMutex.

# Why

In Go, RWMutex is your best choice in most cases when handling concurrent read/write situation. However, in some special extreme cases: you do not want the high-concurrent reads to be blocked by infrequent writes.

In this library, an example left-right-pattern map was implemented to compare with the classic RWMutex-based map. 

You can apply the same technique on slice, list, stack etc, to make them concurrent-safe if it meets your case: high concurrency reads over a single-writer.

``` Go

// Only writes: RWMutex-Map is 2x faster than LR-Map 
BenchmarkLRMap_Write                  	12231741	        93.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkLRMap_Write-2                	12355836	        93.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkLRMap_Write-4                	12448788	        92.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkLockMap_Write                	20728842	        54.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkLockMap_Write-2              	20953489	        54.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkLockMap_Write-4              	20412688	        54.0 ns/op	       0 B/op	       0 allocs/op

// Only reads: RWMutex-Map is slightly faster than LR-Map
BenchmarkLRMap_Read                   	 4505625	       224 ns/op	      39 B/op	       0 allocs/op
BenchmarkLRMap_Read-2                 	 5114546	       212 ns/op	      34 B/op	       0 allocs/op
BenchmarkLRMap_Read-4                 	 5271961	       212 ns/op	      33 B/op	       0 allocs/op
BenchmarkLockMap_Read                 	 7331290	       158 ns/op	      11 B/op	       0 allocs/op
BenchmarkLockMap_Read-2               	 6270297	       162 ns/op	      14 B/op	       0 allocs/op
BenchmarkLockMap_Read-4               	 6244382	       162 ns/op	      14 B/op	       0 allocs/op

// 5 readers vs 1 writer: on 2.4 CPUs, LR-map is 2-3x faster than RWMutex-map
BenchmarkLRMap_Read_Write_5_1         	    5794	    205062 ns/op	      46 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_5_1-2       	    3088	    405067 ns/op	      72 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_5_1-4       	    3014	    415573 ns/op	      81 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_5_1       	    7491	    160224 ns/op	      27 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_5_1-2     	    1210	   1031627 ns/op	      94 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_5_1-4     	    1106	   1067645 ns/op	     101 B/op	       1 allocs/op


// 10 readers vs 1 writer: on 2.4 CPUs, LR-map is 2-3x faster than RWMutex-map
BenchmarkLRMap_Read_Write_10_1        	    3236	    354466 ns/op	      70 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_10_1-2      	    1609	    772303 ns/op	     123 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_10_1-4      	    1593	    726723 ns/op	     132 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_10_1      	    4263	    271487 ns/op	      36 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_10_1-2    	     676	   1726200 ns/op	     154 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_10_1-4    	     604	   2007962 ns/op	     172 B/op	       1 allocs/op

// 50 readers vs 1 writer: on 2.4 CPUs, LR-map is 2x faster than RWMutex-map
BenchmarkLRMap_Read_Write_50_1        	     586	   1733163 ns/op	     310 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_50_1-2      	     340	   3532511 ns/op	     527 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_50_1-4      	     367	   3192197 ns/op	     493 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_50_1      	    1004	   1177003 ns/op	     101 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_50_1-2    	     223	   5803691 ns/op	     433 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_50_1-4    	     195	   5741745 ns/op	     466 B/op	       1 allocs/op

// 100 readers vs 1 writer: on 2.4 CPUs, LR-map is 2x faster than RWMutex-map
BenchmarkLRMap_Read_Write_100_1       	     384	   3126528 ns/op	     471 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_100_1-2     	     182	   6503080 ns/op	     974 B/op	       1 allocs/op
BenchmarkLRMap_Read_Write_100_1-4     	     182	   6345014 ns/op	     969 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_100_1     	     500	   2351661 ns/op	     188 B/op	       1 allocs/op
BenchmarkLockMap_Read_Write_100_1-2   	     100	  10701998 ns/op	     973 B/op	       2 allocs/op
BenchmarkLockMap_Read_Write_100_1-4   	     124	  10646446 ns/op	     776 B/op	       1 allocs/op
```


# What

- Left-Right is a technique with some similarities with Double Instance Locking because it uses two instances, and has a mutex to serialize mutations. In this sense, it is reminiscent of the “Double Buffering” technique used in computer graphics. 
- Left-Right is a generic, linearizable, technique that can provide mutual exclusivity for any object in memory or data structure
- Left-Right is non-blocking for read-only operations, and fast
- it is Wait-Free Population Oblivious

See all the details: [link](https://github.com/CppCon/CppCon2015/blob/master/Presentations/How%20to%20make%20your%20data%20structures%20wait-free%20for%20reads/How%20to%20make%20your%20data%20structures%20wait-free%20for%20reads%20-%20Pedro%20Ramalhete%20-%20CppCon%202015.pdf)


# How

```bash
go get github.com/csimplestring/go-left-right
```

Then you can use it to wrap any data structures. See the below example.

``` Go

import github.com/csimplestring/go-left-right/primitive

type LRMap struct {
	*primitive.LeftRightPrimitive

    // you have to provides 2 identical instances
	left  map[int]int
	right map[int]int
}

func newIntMap() *LRMap {

	m := &LRMap{
		left:  make(map[int]int),
		right: make(map[int]int),
	}

	m.LeftRightPrimitive = primitive.New()

	return m
}

func (lr *LRMap) Get(k int) (val int, exist bool) {

    // Go does not have generics, so have to use interface{} for lambda's arguments
	lr.ApplyReadFn(lr.left, lr.right, func(instance interface{}) {
		m := instance.(map[int]int)
		val, exist = m[k]
	})

	return
}

func (lr *LRMap) Put(key, val int) {
	lr.ApplyWriteFn(lr.left, lr.right, func(instance interface{}) {
		m := instance.(map[int]int)
		m[key] = val
	})
}
```



