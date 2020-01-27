package main

import (
	"strconv"
	"testing"

	"github.com/philipdexter/neolite/lib/query/eager"
	"github.com/philipdexter/neolite/lib/query/eagerchan"
	"github.com/philipdexter/neolite/lib/query/lazy"
	"github.com/philipdexter/neolite/lib/storage"
)

const numNodes = 10000

var filterFunc = func(n storage.Node) bool {
	i, _ := strconv.ParseInt(n.Label, 10, 32)
	return i%2 == 0
}

func BenchmarkLazy(b *testing.B) {
	b.StopTimer()

	storage.Init(numNodes)
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
		lazy.SubmitQuery(
			lazy.Pipeline(
				lazy.ScanAllPipe(),
				lazy.FilterPipe(filterFunc),
				lazy.AccumPipe(),
			))
		lazy.SubmitQuery(
			lazy.Pipeline(
				lazy.ScanAllPipe(),
				lazy.FilterPipe(filterFunc),
				lazy.AccumPipe(),
			))
		b.StartTimer()

		lazy.Run()
	}
}

func BenchmarkEager(b *testing.B) {
	b.StopTimer()
	storage.Init(numNodes)
	eager.InitData(storage.GetData())

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query1 :=
			eager.Pipeline(
				eager.ScanAllPipe(),
				eager.FilterPipe(filterFunc),
				eager.AccumPipe(),
			)
		query2 :=
			eager.Pipeline(
				eager.ScanAllPipe(),
				eager.FilterPipe(filterFunc),
				eager.AccumPipe(),
			)
		query3 :=
			eager.Pipeline(
				eager.ScanAllPipe(),
				eager.FilterPipe(filterFunc),
				eager.AccumPipe(),
			)
		b.StartTimer()

		eager.RunQuery(&query1)
		eager.RunQuery(&query2)
		eager.RunQuery(&query3)
	}
}

func BenchmarkEagerChan(b *testing.B) {
	b.StopTimer()
	storage.Init(numNodes)
	eagerchan.InitData(storage.GetData())

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		query1 :=
			eagerchan.Pipeline(
				eagerchan.ScanAllPipe(),
				eagerchan.FilterPipe(filterFunc),
				eagerchan.AccumPipe(),
			)
		query2 :=
			eagerchan.Pipeline(
				eagerchan.ScanAllPipe(),
				eagerchan.FilterPipe(filterFunc),
				eagerchan.AccumPipe(),
			)
		query3 :=
			eagerchan.Pipeline(
				eagerchan.ScanAllPipe(),
				eagerchan.FilterPipe(filterFunc),
				eagerchan.AccumPipe(),
			)
		b.StartTimer()

		eagerchan.RunQuery(&query1)
		eagerchan.RunQuery(&query2)
		eagerchan.RunQuery(&query3)
	}
}
