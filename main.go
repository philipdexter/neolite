package main

import (
	"fmt"
	"strconv"

	"github.com/philipdexter/neolite/lib/query/eager"
	"github.com/philipdexter/neolite/lib/query/eagerchan"
	"github.com/philipdexter/neolite/lib/query/lazy"
	"github.com/philipdexter/neolite/lib/storage"
)

func main() {
	storage.Init(1000)

	lazy.Init()
	lazy.InitData(storage.GetData())
	const steps = 100

	eager.InitData(storage.GetData())

	eagerchan.InitData(storage.GetData())

	filterFunc := func(n storage.Node) bool {
		i, _ := strconv.ParseInt(n.Label, 10, 32)
		return i%2 == 0
	}

	lazy.SubmitQuery(
		lazy.Pipeline(
			lazy.ScanAllPipe(),
			lazy.FilterPipe(filterFunc),
			lazy.AccumPipe(),
		))

	eagerQuery :=
		eager.Pipeline(
			eager.ScanAllPipe(),
			eager.FilterPipe(filterFunc),
			eager.AccumPipe(),
		)

	eagerchanQuery :=
		eagerchan.Pipeline(
			eagerchan.ScanAllPipe(),
			eagerchan.FilterPipe(filterFunc),
			eagerchan.AccumPipe(),
		)

	fmt.Println(eager.RunQuery(&eagerQuery))

	fmt.Println(eagerchan.RunQuery(&eagerchanQuery))

	for i := 0; i < steps; i++ {
		if i == steps-1 {
			fmt.Println(lazy.Step())
		} else {
			lazy.Step()
		}
	}
}
