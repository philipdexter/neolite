package eager

import (
	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

type result struct {
	results []storage.Node
}

type pipeResult = []storage.Node

func RunQuery(pipeline *pipeline) result {
	return pipeline.Run()
}

type pipeline struct {
	pipes []pipe
	pos   int
}

func (p pipeline) Run() result {
	return result{p.pipes[len(p.pipes)-1].Run(p, len(p.pipes)-1)}
}

type pipe interface {
	Run(pipeline pipeline, pos int) pipeResult
}

type accumPipe struct {
}

func (p *accumPipe) Run(pipeline pipeline, pos int) pipeResult {
	res := make([]storage.Node, 0, len(_data.Data))

	for {
		x := pipeline.pipes[pos-1].Run(pipeline, pos-1)
		if x == nil {
			break
		}
		res = append(res, x...)
	}

	return res
}

type scanAllPipe struct {
	pos int
	buf []storage.Node
}

func (p *scanAllPipe) Run(pipeline pipeline, pos int) pipeResult {
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

func (p *filterPipe) Run(pipeline pipeline, pos int) pipeResult {
	if p.buf == nil {
		p.buf = make([]storage.Node, 1, 1)
	}

	x := pipeline.pipes[pos-1].Run(pipeline, pos-1)
	for {
		if x == nil || p.filterOn(x[0]) {
			break
		}
		x = pipeline.pipes[pos-1].Run(pipeline, pos-1)
	}
	return x
}

func Pipeline(pipes ...pipe) pipeline {
	return pipeline{pipes, 0}
}

func ScanAllPipe() pipe {
	return &scanAllPipe{}
}

func FilterPipe(f func(storage.Node) bool) pipe {
	return &filterPipe{
		filterOn: f,
	}
}

func AccumPipe() pipe {
	return &accumPipe{}
}

func InitData(data *storage.Data) {
	_data = data
}
