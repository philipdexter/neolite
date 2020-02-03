package lazy

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philipdexter/neolite/lib/storage"
)

func TestQueryParse(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		Query("")
	})
	assert.Panics(func() {
		Query("RETURN n")
	})
	assert.Panics(func() {
		Query("MATCH (n)")
	})
	assert.Panics(func() {
		Query("MATCH n\nRETURN n")
	})
	assert.Panics(func() {
		Query("MATCH (n)\nRETURN (n)")
	})
	assert.Panics(func() {
		Query("MATCH (a)\nRETURN n")
	})
	assert.Panics(func() {
		Query("MATCH (a:)\nRETURN a")
	})

	assert.NotNil(Query("MATCH (a)\nRETURN a"))
	assert.NotNil(Query("MATCH (n)\nRETURN n"))
	assert.NotNil(Query("MATCH (n:label)\nRETURN n"))

	pipeline := Query("MATCH (n)\nRETURN n")
	assert.Len(pipeline.pipes, 2)
	assert.IsType(&scanAllPipe{}, pipeline.pipes[0])
	assert.IsType(&accumPipe{}, pipeline.pipes[1])

	pipeline = Query("MATCH (n:alabel)\nRETURN n")
	assert.Len(pipeline.pipes, 2)
	assert.IsType(&scanByLabelPipe{}, pipeline.pipes[0])
	assert.IsType(&accumPipe{}, pipeline.pipes[1])
}

func TestQueryRun(t *testing.T) {
	assert := assert.New(t)

	err := storage.FromString("NODES\n0 0\n1 1\nPROPS\nRELS\n")
	assert.Nil(err)

	Init()
	reqID := SubmitQuery(Query("MATCH (n)\nRETURN n"))
	Run()
	cr := Claim(reqID)
	assert.Len(cr.Rows, 2)
	assert.Equal(cr.Rows[0][0], "0")
	assert.Equal(cr.Rows[1][0], "1")

	storage.Reset()
	err = storage.FromString("NODES\n0 person\n1 car\n2 person\nPROPS\nRELS\n")
	assert.Nil(err)

	Init()
	reqID = SubmitQuery(Query("MATCH (n:person)\nRETURN n"))
	Run()
	cr = Claim(reqID)
	assert.Len(cr.Rows, 2)
	assert.Equal(cr.Rows[0][0], "person")
	assert.Equal(cr.Rows[1][0], "person")

	Init()
	reqID = SubmitQuery(Query("MATCH (n:car)\nRETURN n"))
	Run()
	cr = Claim(reqID)
	assert.Len(cr.Rows, 1)
	assert.Equal(cr.Rows[0][0], "car")
}
