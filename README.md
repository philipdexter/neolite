
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
BenchmarkLazy-4        	     716	   1710507 ns/op	 1180998 B/op	      20 allocs/op
BenchmarkEager-4       	     718	   1756687 ns/op	  983232 B/op	       9 allocs/op
BenchmarkEagerChan-4   	     188	   6517803 ns/op	 1000291 B/op	      27 allocs/op
```

## layout

![struct layout](docs/structs.svg)
