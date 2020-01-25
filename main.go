package main

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/data"
)

func main() {
	fmt.Println("before")
	data.Print()

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
