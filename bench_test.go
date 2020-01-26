package main

import (
	"strconv"
	"testing"

	"github.com/philipdexter/neolite/lib/query/eager"
	"github.com/philipdexter/neolite/lib/query/lazy"
	"github.com/philipdexter/neolite/lib/storage"
)

var filterFunc = func(n storage.Node) bool {
	i, _ := strconv.ParseInt(n.Label, 10, 32)
	return i%2 == 0
}

func BenchmarkLazy(b *testing.B) {
	b.StopTimer()

	const steps = 100
	storage.Init(1000)
	lazy.InitData(storage.GetData())

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		lazy.Init()
		lazy.SubmitQuery(
			lazy.Pipeline(
				lazy.ScanAllPipe(),
				lazy.FilterPipe(filterFunc),
				lazy.AccumPipe(),
			))
		b.StartTimer()

		for j := 0; j < steps; j++ {
			lazy.Step()
		}
	}
}

func BenchmarkEager(b *testing.B) {
	b.StopTimer()
	storage.Init(1000)
	eager.InitData(storage.GetData())

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query :=
			eager.Pipeline(
				eager.ScanAllPipe(),
				eager.FilterPipe(filterFunc),
				eager.AccumPipe(),
			)
		b.StartTimer()
		eager.RunQuery(&query)
	}
}
