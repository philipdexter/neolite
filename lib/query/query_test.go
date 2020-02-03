package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLex(t *testing.T) {
	assert := assert.New(t)

	items := lex("MATCH (n)\nRETURN n")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "n"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "n"},
	})

	items = lex("MATCH (a)RETURN a")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "a"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "a"},
	})

	items = lex("MATCH ( x ) RETURN x")
	assert.Equal(items, []Item{
		Item{TypeMatch, ""},
		Item{TypeLParen, ""},
		Item{TypeName, "x"},
		Item{TypeRParen, ""},
		Item{TypeReturn, ""},
		Item{TypeName, "x"},
	})

	items = lex("MATCH (n:alabel) RETURN n")
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

	astNodes, err := parse(lex("MATCH (n:alabel) RETURN n"))
	assert.Nil(err)
	assert.Len(astNodes, 2)

	astNodes, err = parse(lex("MATCH (n) RETURN n"))
	assert.Nil(err)
	assert.Len(astNodes, 2)

	astNodes, err = parse(lex("MATCH (x) RETURN n"))
	assert.NotNil(err)

	astNodes, err = parse(lex("(x) RETURN n"))
	assert.NotNil(err)

	astNodes, err = parse(lex("MATCH (x)"))
	assert.NotNil(err)

	astNodes, err = parse(lex("MATCH (x) MATCH (x)"))
	assert.NotNil(err)

	astNodes, err = parse(lex(""))
	assert.NotNil(err)
}
