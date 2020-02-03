package lazy

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

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

	randPos := random.Int31n(int32(len(_shadow.items)))
	pipeline := _shadow.items[randPos]
	allowed := maxAllowed
	if randPos > 0 {
		allowed = _shadow.items[randPos-1].pos - pipeline.pos
	}
	if allowed < 0 {
		panic("allowed < 0")
	}
	res := pipeline.run(allowed)
	pipeline.pos = res.pos
	pipeline.status = res.status

	return res
}

type pipeline struct {
	pipes  []pipe
	pos    int
	status int
}

func (p pipeline) run(allowed int) result {
	return *p.pipes[len(p.pipes)-1].run(allowed, p, len(p.pipes)-1)
}

type pipe interface {
	run(allowed int, pipeline pipeline, pos int) *result
}

type limitPipe struct {
	limit         int
	result        result
}

func (p *limitPipe) run(allowed int, pipeline pipeline, pos int) *result {
	if p.result.results == nil {
		p.result.results = make([]storage.Node, 0, p.limit)
	}
	if p.limit < allowed {
		allowed = p.limit
	}
	res := pipeline.pipes[pos-1].run(allowed, pipeline, pos-1)

	remaining := p.limit - len(p.result.results)
	if remaining >= len(p.result.results) {
		p.result.results = append(p.result.results, res.results...)
	} else {
		p.result.results = append(p.result.results, res.results[:remaining]...)
	}

	if len(p.result.results) < p.limit {
		p.result.status = res.status
	} else if len(p.result.results) == p.limit {
		p.result.status = statusDone
	} else {
		panic("len(p.result.results) > p.limit")
	}

	p.result.pos = res.pos

	return &p.result
}

type accumPipe struct {
	result result
}

func (p *accumPipe) run(allowed int, pipeline pipeline, pos int) *result {
	if p.result.results == nil {
		p.result.results = make([]storage.Node, 0, len(storage.Nodes()))
	}
	res := pipeline.pipes[pos-1].run(allowed, pipeline, pos-1)
	p.result.results = append(p.result.results, res.results...)
	p.result.status = res.status
	p.result.pos = res.pos

	return &p.result
}

type scanByLabelPipe struct {
	label  string
	pos    int
	buf    []storage.Node
	result result
}

func (p *scanByLabelPipe) run(allowed int, pipeline pipeline, pos int) *result {
	if p.buf == nil {
		p.buf = make([]storage.Node, 0, maxAllowed)
	} else {
		p.buf = p.buf[:0]
	}
	end := p.pos + allowed
	for ; p.pos < end && p.pos < len(storage.Nodes()); p.pos++ {
		if storage.Nodes()[p.pos].Label == p.label {
			p.buf = append(p.buf, storage.Nodes()[p.pos])
		} else {
			end++
		}
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

type scanAllPipe struct {
	pos    int
	buf    []storage.Node
	result result
}

func (p *scanAllPipe) run(allowed int, pipeline pipeline, pos int) *result {
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

func (p *filterPipe) run(allowed int, pipeline pipeline, pos int) *result {
	res := pipeline.pipes[pos-1].run(allowed, pipeline, pos-1)
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

// Pipeline creates a pipeline from pipes
func Pipeline(pipes ...pipe) pipeline {
	return pipeline{pipes, 0, statusNotDone}
}

// ScanAllPipe creates a pipe
// which scans all nodes of the graph sequentially
func ScanAllPipe() pipe {
	return &scanAllPipe{}
}

// ScanByLabelPipe creates a pipe
// which scans all nodes of the graph sequentially
// and only keeping those with the set label
func ScanByLabelPipe(label string) pipe {
	return &scanByLabelPipe{
		label: label,
	}
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

// LimitPipe creates a pipe
// which accumulates at most n items
// of its input into an array
func LimitPipe(n int) pipe {
	return &limitPipe{
		limit: n,
	}
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

var random *rand.Rand

// Init initializes the lazy processing engine
func Init() {
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	_shadow = shadow{make([]*pipeline, 0)}
}

var matchRegexp *regexp.Regexp
var returnRegexp *regexp.Regexp

func init() {
	matchRegexp = regexp.MustCompile(`MATCH \(([a-z])(?::([A-Za-z]+))?\)`)
	returnRegexp = regexp.MustCompile(`RETURN ([a-z])`)
}

// Query takes a string and builds a pipeline
// from it
func Query(query string) pipeline {
	// TODO use lexer/parser

	lines := strings.Split(query, "\n")

	if len(lines) != 2 {
		panic("query must have one match statement and one return statement")
	}

	if !strings.HasPrefix(lines[0], "MATCH") {
		panic("first statement of query must start with MATCH")
	}

	if !strings.HasPrefix(lines[len(lines)-1], "RETURN") {
		panic("last statement of query must start with RETURN")
	}

	matchM := matchRegexp.FindStringSubmatch(lines[0])
	if len(matchM) != 3 {
		panic("invalid match statement")
	}
	nodeVar := matchM[1]
	labelSelect := matchM[2]

	matchR := returnRegexp.FindStringSubmatch(lines[len(lines)-1])
	if len(matchR) != 2 {
		panic("invalid return statement")
	}
	if matchR[1] != nodeVar {
		panic(fmt.Sprintf("invalid return statement %s %s", nodeVar, lines[1]))
	}

	pipes := make([]pipe, 0, 2)

	if labelSelect != "" {
		pipes = append(pipes, ScanByLabelPipe(labelSelect))
	} else {
		pipes = append(pipes, ScanAllPipe())
	}

	pipes = append(pipes, AccumPipe())

	return pipeline{pipes, 0, statusNotDone}
}
