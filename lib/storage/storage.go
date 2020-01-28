package storage

import (
	"bufio"
	"fmt"
	"os"
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

const (
	modeNodes = iota
	modeProps
	modeRels
)

// FromFile loads data from a file with the syntax
// NODES
// int label
// ...
// PROPS
// int propName propVal
// ...
// RELS
// int int
// ...
func FromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("could not read file %v\n", file)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	getLine := func() (string, error) {
		var pref = true
		var err error
		var l []byte
		var line []byte
		for pref && err == nil {
			line, pref, err = r.ReadLine()
			l = append(l, line...)
		}
		return string(l), err
	}

	_data = Data{
		make([]Node, 0, 100),
	}

	mode := modeNodes
	newNodeIndex := 0
	nodeMap := make(map[int]int)

	s, e := getLine()
	for ; e == nil; s, e = getLine() {
		fmt.Println(s)
		if s == "NODES" {
			mode = modeNodes
			continue
		} else if s == "PROPS" {
			mode = modeProps
			continue
		} else if s == "RELS" {
			mode = modeRels
			continue
		}

		if mode == modeNodes {
			var i int
			var label string
			fmt.Println(s)
			_, err := fmt.Sscanf(s, "%d %s", &i, &label)
			if err != nil {
				panic(err)
			}
			fmt.Printf("got %v and %v\n", i, label)
			node := NewNode(label)
			_data.Data = append(_data.Data, node)
			nodeMap[i] = newNodeIndex
			newNodeIndex++
		} else {
			panic("not implemented")
		}
	}
	if e != nil {
		return e
	}

	return nil
}

// Print pretty prints the data
func Print() {
	fmt.Println(_data.Data)
}
