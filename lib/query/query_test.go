package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLex(t *testing.T) {
	assert := assert.New(t)

	items := Lex("MATCH (n)\nRETURN n")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "n"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "n"},
	})

	items = Lex("MATCH (a)RETURN a")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "a"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "a"},
	})

	items = Lex("MATCH ( x ) RETURN x")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "x"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "x"},
	})

	items = Lex("MATCH (n:alabel) RETURN n")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "n"},
		Item{TypeColon, ""},
		Item{TypeName, "alabel"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "n"},
	})
}

func TestParse(t *testing.T) {
	assert := assert.New(t)

	astNodes, err := Parse(Lex("MATCH (n:alabel) RETURN n"))
	assert.Nil(err)
	assert.Len(astNodes, 2)

	astNodes, err = Parse(Lex("MATCH (n) RETURN n"))
	assert.Nil(err)
	assert.Len(astNodes, 2)

	astNodes, err = Parse(Lex("MATCH (x) RETURN n"))
	assert.NotNil(err)

	astNodes, err = Parse(Lex("(x) RETURN n"))
	assert.NotNil(err)

	astNodes, err = Parse(Lex("MATCH (x)"))
	assert.NotNil(err)

	astNodes, err = Parse(Lex("MATCH (x) MATCH (x)"))
	assert.NotNil(err)

	astNodes, err = Parse(Lex(""))
	assert.NotNil(err)
}
