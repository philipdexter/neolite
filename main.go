package main

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/query"
	"github.com/philipdexter/neolite/lib/storage"
)

func main() {
	query.InitData(storage.GetData())

	fmt.Println("before")
	storage.Print()
	query.Print()

	query.Query(
		query.Pipeline(
			query.ScanAllPipe(),
			query.FilterPipe(func(i int64) bool { return i%2 == 0 }),
			query.AccumPipe(),
		))

	query.Query(
		query.Pipeline(
			query.ScanAllPipe(),
			query.AccumPipe(),
		))

	fmt.Println("=====")
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println(query.Step())
	fmt.Println("=====")

	fmt.Println("after")
	storage.Print()
	query.Print()
}
