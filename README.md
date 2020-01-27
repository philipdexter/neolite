
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
BenchmarkLazy-4            21764             54845 ns/op           33408 B/op          3 allocs/op
BenchmarkEager-4           22674             53174 ns/op           32832 B/op          3 allocs/op
BenchmarkEagerChan-4        7064            252803 ns/op           38511 B/op          9 allocs/op
```

## layout

![struct layout](docs/structs.svg)
