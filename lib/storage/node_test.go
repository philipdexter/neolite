package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"testing/quick"
)

func TestSetProperty(t *testing.T) {
	n := NewNode("label")

	f := func(name, val string) bool {
		n.firstProp = nil
		n.SetProperty(name, val)
		v, err := n.FindProp(name)
		if err != nil {
			return false
		}
		return val == v
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

	n2 := NewNode("label2")

	n.AddRelationship(&n2, "type 1")
	n2.AddRelationship(&n, "type 2")
	typ, err := n.FindFirstRelTypeTo(&n2)
	assert.Nil(t, err)
	assert.Equal(t, typ, "type 1")
	typ, err = n2.FindFirstRelTypeTo(&n)
	assert.Nil(t, err)
	assert.Equal(t, typ, "type 2")
}
