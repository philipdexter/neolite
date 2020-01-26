package eager

import (
	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

const (
	statusDone = iota
	statusNotDone
)

type result struct {
	results []int64
}

type pipeResult = []int64

func RunQuery(pipeline *pipeline) result {
	return pipeline.Run()
}

type pipeline struct {
	pipes []pipe
	pos   int64
}

func (p pipeline) Run() result {
	return result{p.pipes[len(p.pipes)-1].Run(p, len(p.pipes)-1)}
}

////// Pipes

type pipe interface {
	Run(pipeline pipeline, pos int) pipeResult
}

type accumPipe struct {
}

func (p *accumPipe) Run(pipeline pipeline, pos int) pipeResult {
	res := make([]int64, 0, len(_data.Data))

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
	pos int64
	buf []int64
}

func (p *scanAllPipe) Run(pipeline pipeline, pos int) pipeResult {
	if p.buf == nil {
		p.buf = make([]int64, 1, 1)
	}
	if int64(len(_data.Data)) > p.pos {
		x := _data.Data[p.pos]
		p.pos++
		p.buf[0] = x
		return p.buf
	}
	return nil
}

type filterPipe struct {
	filterOn func(int64) bool
	buf      []int64
}

func (p *filterPipe) Run(pipeline pipeline, pos int) pipeResult {
	if p.buf == nil {
		p.buf = make([]int64, 1, 1)
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

func FilterPipe(f func(int64) bool) pipe {
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
