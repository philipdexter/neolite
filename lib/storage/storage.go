package storage

import (
	"fmt"
)

type Data struct {
	Data []int64
}

var _data Data

func GetData() *Data {
	return &_data
}

func Init(size int) {

	_data = Data{
		make([]int64, size),
	}
	for i := 0; i < size; i++ {
		_data.Data[i] = int64(i)
	}
}

func Print() {
	fmt.Println(_data.Data)
}
