package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"testing/quick"
)

func TestNode(t *testing.T) {
	assert := assert.New(t)

	n := InsertNode("label")

	f := func(name, val string) bool {
		_nodes[n].firstProp = noId
		SetProperty(n, name, val)
		v, err := FindProp(n, name)
		if err != nil {
			return false
		}
		return val == v
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

	n2 := InsertNode("label2")

	AddRelationship(n, n2, "type 1")
	AddRelationship(n2, n, "type 2")
	typ, err := FindFirstRelTypeTo(n, n2)
	assert.Nil(err)
	assert.Equal(typ, "type 1")
	typ, err = FindFirstRelTypeTo(n2, n)
	assert.Nil(err)
	assert.Equal(typ, "type 2")
}
