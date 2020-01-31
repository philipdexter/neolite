package storage

import (
	"bufio"
	"fmt"
	"os"
)

var _nodes []Node
var _props []Property
var _rels []Relationship

// Reset removes all nodes, props, and relationships, and sets their storage to nil
func Reset() {
	_nodes = nil
	_props = nil
	_rels = nil
}

// Nodes gets the graph nodes
func Nodes() []Node {
	return _nodes
}

// Props gets the graph properties
func Props() []Property {
	return _props
}

// Rels gets the graph relationships
func Rels() []Relationship {
	return _rels
}

const initCap = 1000

func init() {
	_nodes = make([]Node, 0, initCap)
	_props = make([]Property, 0, initCap)
	_rels = make([]Relationship, 0, initCap)
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
// int int type
// ...
func FromFile(file string) error {
	if len(_nodes) > 0 || len(_props) > 0 || len(_rels) > 0 {
		panic("cannot call FromFile when data exists")
	}

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

	mode := modeNodes
	nodeMap := make(map[int]int)

	s, e := getLine()
	for ; e == nil; s, e = getLine() {
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
			_, err := fmt.Sscanf(s, "%d %s", &i, &label)
			if err != nil {
				panic(err)
			}
			nodeID := InsertNode(label)
			nodeMap[i] = nodeID
		} else if mode == modeProps {
			var i int
			var propName string
			var propVal string
			_, err := fmt.Sscanf(s, "%d %s %s", &i, &propName, &propVal)
			if err != nil {
				panic(err)
			}
			SetProperty(nodeMap[i], propName, propVal)
		} else if mode == modeRels {
			var from int
			var to int
			var typ string
			_, err := fmt.Sscanf(s, "%d %d %s", &from, &to, &typ)
			if err != nil {
				panic(err)
			}
			AddRelationship(nodeMap[from], nodeMap[to], typ)
		} else {
			panic("invalid mode")
		}
	}
	if e != nil {
		return e
	}

	return nil
}

// Print pretty prints the data
func Print() {
	fmt.Println(_nodes)
	fmt.Println(_props)
	fmt.Println(_rels)
}
