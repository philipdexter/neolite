package main

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/query/eager"
	"github.com/philipdexter/neolite/lib/query/lazy"
	"github.com/philipdexter/neolite/lib/storage"
)

func main() {
	storage.Init(1000)

	lazy.Init()
	lazy.InitData(storage.GetData())
	const steps = 100

	eager.InitData(storage.GetData())

	lazy.SubmitQuery(
		lazy.Pipeline(
			lazy.ScanAllPipe(),
			lazy.FilterPipe(func(i int64) bool { return i%2 == 0 }),
			lazy.AccumPipe(),
		))

	eagerQuery :=
		eager.Pipeline(
			eager.ScanAllPipe(),
			eager.FilterPipe(func(i int64) bool { return i%2 == 0 }),
			eager.AccumPipe(),
		)

	fmt.Println(eager.RunQuery(&eagerQuery))

	for i := 0; i < steps; i++ {
		if i == steps-1 {
			fmt.Println(lazy.Step())
		} else {
			lazy.Step()
		}
	}
}
