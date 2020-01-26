package main

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/query/lazy"
	"github.com/philipdexter/neolite/lib/storage"
)

func main() {
	storage.Init(1000)
	lazy.InitData(storage.GetData())

	fmt.Println("before")
	storage.Print()
	lazy.Print()

	lazy.Query(
		lazy.Pipeline(
			lazy.ScanAllPipe(),
			lazy.FilterPipe(func(i int64) bool { return i%2 == 0 }),
			lazy.AccumPipe(),
		))

	lazy.Query(
		lazy.Pipeline(
			lazy.ScanAllPipe(),
			lazy.AccumPipe(),
		))

	fmt.Println("=====")
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println(lazy.Step())
	fmt.Println("=====")

	fmt.Println("after")
	storage.Print()
	lazy.Print()
}
