package eager

import (
	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

type result struct {
	results []storage.Node
}

type pipeResult = []storage.Node

// RunQuery runs a single pipeline
func RunQuery(pipeline *pipeline) result {
	return pipeline.run()
}

type pipeline struct {
	pipes []pipe
	pos   int
}

func (p pipeline) run() result {
	return result{p.pipes[len(p.pipes)-1].run(p, len(p.pipes)-1)}
}

type pipe interface {
	run(pipeline pipeline, pos int) pipeResult
}

type accumPipe struct {
}

func (p *accumPipe) run(pipeline pipeline, pos int) pipeResult {
	res := make([]storage.Node, 0, len(_data.Data))

	for {
		x := pipeline.pipes[pos-1].run(pipeline, pos-1)
		if x == nil {
			break
		}
		res = append(res, x[0])
	}

	return res
}

type scanAllPipe struct {
	pos int
	buf []storage.Node
}

func (p *scanAllPipe) run(pipeline pipeline, pos int) pipeResult {
	if p.buf == nil {
		p.buf = make([]storage.Node, 1, 1)
	}
	if len(_data.Data) > p.pos {
		x := _data.Data[p.pos]
		p.pos++
		p.buf[0] = x
		return p.buf
	}
	return nil
}

type filterPipe struct {
	filterOn func(storage.Node) bool
	buf      []storage.Node
}

func (p *filterPipe) run(pipeline pipeline, pos int) pipeResult {
	if p.buf == nil {
		p.buf = make([]storage.Node, 1, 1)
	}

	x := pipeline.pipes[pos-1].run(pipeline, pos-1)
	for {
		if x == nil || p.filterOn(x[0]) {
			break
		}
		x = pipeline.pipes[pos-1].run(pipeline, pos-1)
	}
	return x
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

// InitData sets the data which the lazy processing engine will use
func InitData(data *storage.Data) {
	_data = data
}
