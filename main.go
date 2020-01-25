package main

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/data"
)

func main() {
	fmt.Println("before")
	data.Print()

	data.Query(
		data.Pipeline(
			data.ScanAllPipe(),
			data.FilterPipe(func(i int64) bool { return i%2 == 0 }),
			data.AccumPipe(),
		))

	data.Query(
		data.Pipeline(
			data.ScanAllPipe(),
			data.AccumPipe(),
		))

	fmt.Println("=====")
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println(data.Step())
	fmt.Println("=====")

	fmt.Println("after")
	data.Print()
}
