package lazy

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/philipdexter/neolite/lib/query"
	"github.com/philipdexter/neolite/lib/storage"
)

const maxAllowed = 100

const (
	statusDone = iota
	statusNotDone
)

type pipeResult struct {
	results []storage.Node
	status  int
	pos     int
	reqID   RequestID
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
func Run() {
	for !isDone() {
		res := step()
		if res.status == statusDone {
			cr := &ClientResult{
				Columns: []string{"label"},
				Rows:    make([][]string, 0, len(res.results)),
			}
			for i := 0; i < len(res.results); i++ {
				cr.Rows = append(cr.Rows, []string{res.results[i].Label})
			}
			finishRequest(res.reqID, cr)
		}
	}
}

type ClientResult struct {
	Columns []string
	Rows    [][]string
}

type RequestID = int

var clientCache map[RequestID]*ClientResult

func Claim(reqID RequestID) *ClientResult {
	cr, ok := clientCache[reqID]
	if !ok {
		return nil
	}
	return cr
}

func finishRequest(reqID RequestID, cr *ClientResult) {
	if _, ok := clientCache[reqID]; ok {
		panic("duplicate RequestID in clientCache")
	}
	clientCache[reqID] = cr
}

func step() pipeResult {
	if len(_shadow.items) == 0 {
		return pipeResult{make([]storage.Node, 0), statusDone, 0, 0}
	}

	randPos := random.Int31n(int32(len(_shadow.items)))
	pipeline := _shadow.items[randPos]
	for pipeline.status == statusDone {
		randPos = (randPos + 1) % int32(len(_shadow.items))
		pipeline = _shadow.items[randPos]
	}
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

var reqIDCounter = 0

type pipeline struct {
	pipes  []pipe
	pos    int
	status int
	reqID  RequestID
}

func (p pipeline) run(allowed int) pipeResult {
	return *p.pipes[len(p.pipes)-1].run(allowed, p, len(p.pipes)-1)
}

type pipe interface {
	run(allowed int, pipeline pipeline, pos int) *pipeResult
}

type limitPipe struct {
	limit         int
	result        pipeResult
}

func (p *limitPipe) run(allowed int, pipeline pipeline, pos int) *pipeResult {
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
	p.result.reqID = pipeline.reqID

	return &p.result
}

type accumPipe struct {
	result        pipeResult
}

func (p *accumPipe) run(allowed int, pipeline pipeline, pos int) *pipeResult {
	if p.result.results == nil {
		p.result.results = make([]storage.Node, 0, len(storage.Nodes()))
	}
	res := pipeline.pipes[pos-1].run(allowed, pipeline, pos-1)
	p.result.results = append(p.result.results, res.results...)
	p.result.status = res.status
	p.result.pos = res.pos
	p.result.reqID = pipeline.reqID

	return &p.result
}

type scanByLabelPipe struct {
	label  string
	pos    int
	buf    []storage.Node
	result pipeResult
}

func (p *scanByLabelPipe) run(allowed int, pipeline pipeline, pos int) *pipeResult {
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
	result pipeResult
}

func (p *scanAllPipe) run(allowed int, pipeline pipeline, pos int) *pipeResult {
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

func (p *filterPipe) run(allowed int, pipeline pipeline, pos int) *pipeResult {
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
func SubmitQuery(pipeline pipeline) RequestID {
	_shadow.items = append(_shadow.items, &pipeline)

	return pipeline.reqID
}

// Pipeline creates a pipeline from pipes
func Pipeline(pipes ...pipe) pipeline {
	i := reqIDCounter
	reqIDCounter++
	return pipeline{pipes, 0, statusNotDone, i}
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
	_shadow = shadow{make([]*pipeline, 0)}
	clientCache = make(map[RequestID]*ClientResult, 0)
}

var matchRegexp *regexp.Regexp
var returnRegexp *regexp.Regexp

func init() {
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	matchRegexp = regexp.MustCompile(`MATCH \(([a-z])(?::([A-Za-z]+))?\)`)
	returnRegexp = regexp.MustCompile(`RETURN ([a-z])`)

	Init()
}

// Query takes a string and builds a pipeline
// from it
func Query(queryString string) pipeline {
	astNodes, err := query.Parse(query.Lex(queryString))
	if err != nil {
		panic(err)
	}

	if len(astNodes) != 2 {
		panic("len(astNodes) != 2")
	}

	pipes := make([]pipe, 0, 2)

	switch astNode := astNodes[0].(type) {
	case query.Match:
		if astNode.LabelFilter != "" {
			pipes = append(pipes, ScanByLabelPipe(astNode.LabelFilter))
		} else {
			pipes = append(pipes, ScanAllPipe())
		}
	default:
		panic("first astNode is not a match statement")
	}

	switch astNodes[1].(type) {
	case query.Return:
		pipes = append(pipes, AccumPipe())
	default:
		panic("second astNode is not a return statement")
	}


	return Pipeline(pipes...)
}
