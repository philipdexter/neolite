package eagerchan

import (
	"github.com/philipdexter/neolite/lib/storage"
)

type result struct {
	results []storage.Node
}

type pipeResult = interface{}

// RunQuery runs a single pipeline
func RunQuery(pipeline *pipeline) result {
	return pipeline.run()
}

type pipeline struct {
	pipes []pipe
	pos   int
}

func (p pipeline) run() result {
	chans := make([]chan pipeResult, len(p.pipes))
	for i := 0; i < len(p.pipes); i++ {
		chans[i] = make(chan pipeResult, 100)
	}
	for i := 0; i < len(p.pipes); i++ {
		if i == 0 {
			go p.pipes[i].run(nil, chans[i], p, len(p.pipes)-1-i)
		} else {
			go p.pipes[i].run(chans[i-1], chans[i], p, i)
		}
	}
	return result{(<-chans[len(chans)-1]).([]storage.Node)}
}

type pipe interface {
	run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int)
}

type accumPipe struct {
}

func (p *accumPipe) run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
	res := make([]storage.Node, 0, len(storage.Nodes()))

	for {
		x := <-fromChan
		if x == nil {
			break
		}
		res = append(res, *x.(*storage.Node))
	}

	toChan <- res
}

type scanAllPipe struct {
}

func (p scanAllPipe) run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
	if fromChan != nil {
		panic("fromChan != nil")
	}

	dpos := 0
	for {
		if dpos >= len(storage.Nodes()) {
			break
		}
		toChan <- &storage.Nodes()[dpos]
		dpos++
	}
	toChan <- nil
}

type filterPipe struct {
	filterOn func(storage.Node) bool
}

func (p filterPipe) run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
	for {
		x := <-fromChan
		if x == nil {
			break
		}
		if p.filterOn(*x.(*storage.Node)) {
			toChan <- x
		}
	}
	toChan <- nil
}

// Pipeline creates a pipeline from pipes
func Pipeline(pipes ...pipe) pipeline {
	return pipeline{pipes, 0}
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
