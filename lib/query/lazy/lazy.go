package lazy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

const (
	statusDone = iota
	statusNotDone
)

type result struct {
	results []int64
	status  int
	pos     int64
}

func Step() result {
	if len(_shadow.items) == 0 {
		return result{make([]int64, 0), 0, 0}
	}
	randPos := rand.Int31n(int32(len(_shadow.items)))
	pipeline := _shadow.items[randPos]
	allowed := int64(10)
	if int32(len(_shadow.items)) > randPos+1 {
		allowed = _shadow.items[randPos+1].pos - pipeline.pos
	}
	if allowed < 0 {
		panic("allowed < 0")
	}
	res := pipeline.Run(allowed)
	pipeline.pos = res.pos

	return res
}

type pipeline struct {
	pipes []pipe
	pos   int64
}

func (p pipeline) Run(allowed int64) result {
	return p.pipes[len(p.pipes)-1].Run(allowed, p, len(p.pipes)-1)
}

////// Pipes

type pipe interface {
	Run(allowed int64, pipeline pipeline, pos int) result
}

type accumPipe struct {
	result result
}

func (p *accumPipe) Run(allowed int64, pipeline pipeline, pos int) result {
	res := pipeline.pipes[pos-1].Run(allowed, pipeline, pos-1)
	p.result.results = append(p.result.results, res.results...)
	p.result.status = res.status
	p.result.pos = res.pos

	return p.result
}

type scanAllPipe struct {
	pos int64
}

func (p *scanAllPipe) Run(allowed int64, pipeline pipeline, pos int) result {
	res := make([]int64, 0, allowed)
	end := p.pos + allowed
	for ; p.pos < end && p.pos < int64(len(_data.Data)); p.pos++ {
		res = append(res, _data.Data[p.pos])
	}

	status := statusNotDone
	if p.pos == int64(len(_data.Data)) {
		status = statusDone
	}
	return result{
		results: res,
		status:  status,
		pos:     p.pos,
	}
}

type filterPipe struct {
	filterOn func(int64) bool
}

func (p filterPipe) Run(allowed int64, pipeline pipeline, pos int) result {
	res := pipeline.pipes[pos-1].Run(allowed, pipeline, pos-1)
	filtered := make([]int64, 0, len(res.results)/2)

	for _, i := range res.results {
		if p.filterOn(i) {
			filtered = append(filtered, i)
		}
	}

	res.results = filtered
	return res
}

var _shadow shadow

type shadow struct {
	items []*pipeline
}

func SubmitQuery(pipeline pipeline) {
	_shadow.items = append(_shadow.items, &pipeline)
}

func Pipeline(pipes ...pipe) pipeline {
	return pipeline{pipes, 0}
}

func ScanAllPipe() pipe {
	return &scanAllPipe{}
}

func FilterPipe(f func(int64) bool) pipe {
	return &filterPipe{f}
}

func AccumPipe() pipe {
	return &accumPipe{}
}

func Print() {
	for _, pipeline := range _shadow.items {
		fmt.Printf(" : %v\n", pipeline.pos)
		for _, pipe := range pipeline.pipes {
			fmt.Printf("%T ", pipe)
			fmt.Println(pipe)
		}
	}
}

func Init() {
	rand.Seed(time.Now().UTC().UnixNano())
	_shadow = shadow{make([]*pipeline, 0)}
}

func InitData(data *storage.Data) {
	_data = data
}
