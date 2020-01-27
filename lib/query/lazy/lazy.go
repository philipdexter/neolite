package lazy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/philipdexter/neolite/lib/storage"
)

var _data *storage.Data

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

func Run() result {
	var res result = result{make([]storage.Node, 0), statusDone, 0}
	for !isDone() {
		res = step()
	}
	return res
}

func step() result {
	if len(_shadow.items) == 0 {
		return result{make([]storage.Node, 0), statusDone, 0}
	}

	randPos := rand.Int31n(int32(len(_shadow.items)))
	pipeline := _shadow.items[randPos]
	allowed := maxAllowed
	if randPos > 0 {
		allowed = _shadow.items[randPos-1].pos - pipeline.pos
	}
	if allowed < 0 {
		panic("allowed < 0")
	}
	res := pipeline.Run(allowed)
	pipeline.pos = res.pos
	pipeline.status = res.status

	return res
}

type pipeline struct {
	pipes  []pipe
	pos    int
	status int
}

func (p pipeline) Run(allowed int) result {
	return *p.pipes[len(p.pipes)-1].Run(allowed, p, len(p.pipes)-1)
}

type pipe interface {
	Run(allowed int, pipeline pipeline, pos int) *result
}

type accumPipe struct {
	result result
}

func (p *accumPipe) Run(allowed int, pipeline pipeline, pos int) *result {
	if p.result.results == nil {
		p.result.results = make([]storage.Node, 0, len(_data.Data))
	}
	res := pipeline.pipes[pos-1].Run(allowed, pipeline, pos-1)
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

func (p *scanAllPipe) Run(allowed int, pipeline pipeline, pos int) *result {
	if p.buf == nil {
		p.buf = make([]storage.Node, 0, maxAllowed)
	} else {
		p.buf = p.buf[:0]
	}
	end := p.pos + allowed
	for ; p.pos < end && p.pos < len(_data.Data); p.pos++ {
		p.buf = append(p.buf, _data.Data[p.pos])
	}

	status := statusNotDone
	if p.pos == len(_data.Data) {
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

func (p *filterPipe) Run(allowed int, pipeline pipeline, pos int) *result {
	res := pipeline.pipes[pos-1].Run(allowed, pipeline, pos-1)
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

func SubmitQuery(pipeline pipeline) {
	_shadow.items = append(_shadow.items, &pipeline)
}

func Pipeline(pipes ...pipe) pipeline {
	return pipeline{pipes, 0, statusNotDone}
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
