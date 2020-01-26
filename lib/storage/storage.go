package storage

import (
	"fmt"
)

////// Result

type Data struct {
	Data []int64
}

var _data Data

func GetData() Data {
	return _data
}

func init() {

	_data = Data{
		make([]int64, 100),
	}
	for i := int64(0); i < 100; i++ {
		_data.Data[i] = i
	}
}

func Print() {
	fmt.Println(_data.Data)
}
