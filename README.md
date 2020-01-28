
# neolite

[github.com/philipdexter/neolite](https://github.com/philipdexter/neolite)

[docs](https://godoc.org/github.com/philipdexter/neolite)

A light version of neo4j
to explore data-centric laziness

## benchmark results

```bash
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/philipdexter/neolite
BenchmarkLazy-4        	     697	   1714714 ns/op	 1184421 B/op	      20 allocs/op
BenchmarkLazyFused-4   	    3458	    393271 ns/op	  992648 B/op	       7 allocs/op
BenchmarkEager-4       	     693	   1727678 ns/op	  983233 B/op	       9 allocs/op
BenchmarkEagerChan-4   	     181	   6733494 ns/op	 1000272 B/op	      27 allocs/op
```

## layout

![struct layout](docs/structs.svg)
