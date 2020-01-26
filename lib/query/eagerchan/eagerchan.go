package eagerchan

import (
	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

type result struct {
	results []storage.Node
}

type pipeResult = interface{}

func RunQuery(pipeline *pipeline) result {
	return pipeline.Run()
}

type pipeline struct {
	pipes []pipe
	pos   int
}

func (p pipeline) Run() result {
	chans := make([]chan pipeResult, len(p.pipes))
	for i := 0; i < len(p.pipes); i++ {
		chans[i] = make(chan pipeResult)
	}
	for i := 0; i < len(p.pipes); i++ {
		if i == 0 {
			go p.pipes[i].Run(nil, chans[i], p, len(p.pipes)-1-i)
		} else {
			go p.pipes[i].Run(chans[i-1], chans[i], p, i)
		}
	}
	return result{(<-chans[len(chans)-1]).([]storage.Node)}
}

type pipe interface {
	Run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int)
}

type accumPipe struct {
}

func (p *accumPipe) Run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
	res := make([]storage.Node, 0, len(_data.Data))

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

func (p scanAllPipe) Run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
	if fromChan != nil {
		panic("fromChan != nil")
	}

	dpos := 0
	for {
		if dpos >= len(_data.Data) {
			break
		}
		toChan <- &_data.Data[dpos]
		dpos++
	}
	toChan <- nil
}

type filterPipe struct {
	filterOn func(storage.Node) bool
}

func (p filterPipe) Run(fromChan <-chan pipeResult, toChan chan<- pipeResult, pipeline pipeline, pos int) {
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
