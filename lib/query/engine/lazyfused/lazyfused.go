package lazyfused

import (
	"fmt"

	"github.com/philipdexter/neolite/lib/storage"
)

const maxAllowed = 100

const (
	statusDone = iota
	statusNotDone
)

type result struct {
	results []storage.Node
	status  int
	pos     int
}

func isDone() bool {
	for _, pipeline := range _shadow.items {
		if pipeline.status == statusNotDone {
			return false
		}
	}
	return true
}

// Run runs all submitted queries
func Run() result {
	var res = result{make([]storage.Node, 0), statusDone, 0}
	for !isDone() {
		res = step()
	}
	return res
}

// TODO auto detect when two items are at the same spot and can be fused

func step() result {
	if len(_shadow.items) == 0 {
		return result{make([]storage.Node, 0), statusDone, 0}
	}

	var fuseRef *result

	var res result
	for i := 0; i < len(_shadow.items); i++ {
		pipeline := _shadow.items[i]
		allowed := maxAllowed
		res = pipeline.run(allowed, &fuseRef)
		pipeline.pos = res.pos
		pipeline.status = res.status
	}

	return res
}

type pipeline struct {
	pipes  []pipe
	pos    int
	status int
}

func (p pipeline) run(allowed int, fuseRef **result) result {
	return *p.pipes[len(p.pipes)-1].run(allowed, p, len(p.pipes)-1, fuseRef)
}

type pipe interface {
	run(allowed int, pipeline pipeline, pos int, fuseRef **result) *result
}

type accumPipe struct {
	result result
}

func (p *accumPipe) run(allowed int, pipeline pipeline, pos int, fuseRef **result) *result {
	if p.result.results == nil {
		p.result.results = make([]storage.Node, 0, len(storage.Nodes()))
	}
	res := pipeline.pipes[pos-1].run(allowed, pipeline, pos-1, fuseRef)
	p.result.results = append(p.result.results, res.results...)
	p.result.status = res.status
	p.result.pos = res.pos

	return &p.result
}

type scanAllPipe struct {
	pos    int
	buf    []storage.Node
	result result
}

func (p *scanAllPipe) run(allowed int, pipeline pipeline, pos int, fuseRef **result) *result {
	if p.buf == nil {
		p.buf = make([]storage.Node, 0, maxAllowed)
	} else {
		p.buf = p.buf[:0]
	}
	end := p.pos + allowed
	for ; p.pos < end && p.pos < len(storage.Nodes()); p.pos++ {
		p.buf = append(p.buf, storage.Nodes()[p.pos])
	}

	status := statusNotDone
	if p.pos == len(storage.Nodes()) {
		status = statusDone
	}

	p.result.results = p.buf
	p.result.status = status
	p.result.pos = p.pos
	return &p.result
}

type filterPipe struct {
	filterOn func(storage.Node) bool
	buf      []storage.Node
}

func (p *filterPipe) run(allowed int, pipeline pipeline, pos int, fuseRef **result) *result {
	var res *result
	if *fuseRef == nil {
		res = pipeline.pipes[pos-1].run(allowed, pipeline, pos-1, fuseRef)
		*fuseRef = res
	} else {
		res = *fuseRef
	}
	if p.buf == nil {
		p.buf = make([]storage.Node, 0, maxAllowed)
	} else {
		p.buf = p.buf[:0]
	}

	for _, i := range res.results {
		if p.filterOn(i) {
			p.buf = append(p.buf, i)
		}
	}

	res.results = p.buf
	return res
}

var _shadow shadow

type shadow struct {
	items []*pipeline
}

// SubmitQuery queues a query to be processed by Run
func SubmitQuery(pipeline pipeline) {
	_shadow.items = append(_shadow.items, &pipeline)
}

var fuzedScanAll = ScanAllPipe()

// Pipeline creates a pipeline from pipes
func Pipeline(pipes ...pipe) pipeline {
	pipes = append([]pipe{fuzedScanAll}, pipes...)
	return pipeline{pipes, 0, statusNotDone}
}

// ScanAllPipe creates a pipe
// which scans all nodes of the graph sequentially
func ScanAllPipe() pipe {
	return &scanAllPipe{}
}

// FilterPipe creates a pipe
// which filters its input by the provided function
func FilterPipe(f func(storage.Node) bool) pipe {
	return &filterPipe{
		filterOn: f,
	}
}

// AccumPipe creates a pipe
// which accumulates its input into an array
func AccumPipe() pipe {
	return &accumPipe{}
}

// Print prints out the state of the lazy processing engine
func Print() {
	for _, pipeline := range _shadow.items {
		fmt.Printf(" : %v\n", pipeline.pos)
		for _, pipe := range pipeline.pipes {
			fmt.Printf("%T ", pipe)
			fmt.Println(pipe)
		}
	}
}

// Init initializes the lazy processing engine
func Init() {
	_shadow = shadow{make([]*pipeline, 0)}
}
