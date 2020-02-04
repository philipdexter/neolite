package main

import (
	"fmt"
	"strconv"

	"github.com/philipdexter/neolite/lib/query/engine/eager"
	"github.com/philipdexter/neolite/lib/query/engine/eagerchan"
	"github.com/philipdexter/neolite/lib/query/engine/lazy"
	"github.com/philipdexter/neolite/lib/query/engine/lazyfused"
	"github.com/philipdexter/neolite/lib/storage"
)

func main() {
	storage.FromFile("data.dat")

	lazy.Init()

	lazyfused.Init()

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

	lazyfused.SubmitQuery(
		lazyfused.Pipeline(
			lazyfused.FilterPipe(filterFunc),
			lazyfused.AccumPipe(),
		))
	lazyfused.SubmitQuery(
		lazyfused.Pipeline(
			lazyfused.FilterPipe(filterFunc),
			lazyfused.AccumPipe(),
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

	fmt.Println("eager")
	fmt.Println(eager.RunQuery(&eagerQuery))

	fmt.Println("eagerchan")
	fmt.Println(eagerchan.RunQuery(&eagerchanQuery))

	fmt.Println("lazy")
	lazy.Run()
	fmt.Println(lazy.Claim(0))

	fmt.Println("lazyfuzed")
	fmt.Println(lazyfused.Run())
}
