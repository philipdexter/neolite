package storage

import (
	"fmt"
	"strconv"
)

// Data is an array of nodes
type Data struct {
	Data []Node
}

var _data Data

// GetData returns the singleton data
func GetData() *Data {
	return &_data
}

// Init initializes the singleton data with a size
func Init(size int) {
	_data = Data{
		make([]Node, size),
	}
	for i := 0; i < size; i++ {
		_data.Data[i] = NewNode(strconv.Itoa(i))
	}
}

// Print pretty prints the data
func Print() {
	fmt.Println(_data.Data)
}
