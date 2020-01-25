package data

import "fmt"

////// Result

const (
	statusDone = iota
	statusNotDone
)

type result struct {
	results []int64
	status  int
}

////// Pipeline

type pipeline struct {
	pipes []pipe
}

func (p pipeline) Run(allowed int64) result {
	return p.pipes[len(p.pipes)-1].Run(allowed, p, len(p.pipes)-1)
}

////// Pipes

type pipe interface {
	Run(allowed int64, pipeline pipeline, pos int) result
}

type AccumPipe struct {
	result result
}

func (p *AccumPipe) Run(allowed int64, pipeline pipeline, pos int) result {
	res := pipeline.pipes[pos-1].Run(allowed, pipeline, pos-1)
	p.result.results = append(p.result.results, res.results...)
	p.result.status = res.status

	return p.result
}

type ScanAllPipe struct {
	pos int64
}

func (p *ScanAllPipe) Run(allowed int64, pipeline pipeline, pos int) result {
	res := make([]int64, 0, allowed)
	end := p.pos + allowed
	for ; p.pos < end && p.pos < int64(len(_data.data)); p.pos++ {
		res = append(res, _data.data[p.pos])
	}

	status := statusNotDone
	if p.pos == int64(len(_data.data)) {
		status = statusDone
	}
	return result{
		results: res,
		status:  status,
	}
}

type FilterPipe struct {
	filterOn func(int64) bool
}

func (p FilterPipe) Run(allowed int64, pipeline pipeline, pos int) result {
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

////// Data and Shadow list

type data struct {
	data []int64
}

type shadow struct {
	items []pipeline
}

var _data data
var _shadow shadow

func init() {
	_data = data{
		make([]int64, 100),
	}
	for i := int64(0); i < 100; i++ {
		_data.data[i] = i
	}

	_shadow = shadow{
		[]pipeline{
			pipeline{[]pipe{
				&ScanAllPipe{},
				&FilterPipe{func(i int64) bool { return i%2 == 0 }},
				&AccumPipe{},
			}},
		},
	}
}

func Step() result {
	return _shadow.items[0].Run(10)
}

func printData() {
	fmt.Println(_data.data)
}

func printShadow() {
	for _, pipeline := range _shadow.items {
		for _, pipe := range pipeline.pipes {
			fmt.Printf("%T ", pipe)
			fmt.Println(pipe)
		}
	}
}

func Print() {
	printData()
	printShadow()
}
