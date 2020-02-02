package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromString(t *testing.T) {
	assert := assert.New(t)

	Reset()
	err := FromString("NODES\n0 person\n1 car\n2 person\nPROPS\nRELS\n")
	assert.Nil(err)
	assert.Len(Nodes(), 3)
	assert.Len(Props(), 0)
	assert.Len(Rels(), 0)

	Reset()
	err = FromString("NODES\n0 person\n1 car\n2 person\nPROPS\n0 name bob\n0 age 19\n1 type deluxe\nRELS\n")
	assert.Nil(err)
	assert.Len(Nodes(), 3)
	assert.Len(Props(), 3)
	assert.Len(Rels(), 0)
	assert.Equal(Props()[0].name, "name")
	assert.Equal(Props()[0].val, "bob")
	assert.Equal(Props()[0].next, -1)
	assert.Equal(Props()[1].name, "age")
	assert.Equal(Props()[1].val, "19")
	assert.Equal(Props()[1].next, 0)
	assert.Equal(Props()[2].name, "type")
	assert.Equal(Props()[2].val, "deluxe")
	assert.Equal(Props()[2].next, -1)

	// TODO test rels
}
