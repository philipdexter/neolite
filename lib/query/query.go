package query

import (
	"regexp"
)

type Type uint16

type Item struct {
	Typ Type
	Val string
}

const (
	TypeEOF Type = iota
	TypeMatch
	TypeReturn
	TypeLParen
	TypeRParen
	TypeColon
	TypeName
)

var nameR *regexp.Regexp

func init() {
	nameR = regexp.MustCompile(` |\)|:|\n`)
}

func lex(s string) []Item {
	items := make([]Item, 0)

	pos := 0

	for pos < len(s) {
		if s[pos] == '\n' || s[pos] == ' ' {
			pos += 1
		} else if s[pos] == '(' {
			items = append(items, Item{TypeLParen, ""})
			pos += 1
		} else if s[pos] == ')' {
			items = append(items, Item{TypeRParen, ""})
			pos += 1
		} else if s[pos] == ':' {
			items = append(items, Item{TypeColon, ""})
			pos += 1
		} else if len(s) > pos+len("MATCH ") && s[pos:pos+len("MATCH ")] == "MATCH " {
			items = append(items, Item{TypeMatch, ""})
			pos += len("MATCH ")
		} else if len(s) > pos+len("RETURN ") && s[pos:pos+len("RETURN ")] == "RETURN " {
			items = append(items, Item{TypeReturn, ""})
			pos += len("RETURN ")
		} else {
			nsis := nameR.FindStringIndex(s[pos:])
			var nsi = len(s) - pos
			if len(nsis) > 0 {
				nsi = nsis[0]
			}
			items = append(items, Item{TypeName, s[pos : pos+nsi]})
			pos += nsi
		}
	}

	return items
}

type AST interface {
}

type Match struct {
	nodeName    string
	labelFilter string
}

type Return struct {
	nodeName string
}

type ParseError struct {
	msg string
}

func (e ParseError) Error() string {
	return e.msg
}

func parse(items []Item) ([]AST, error) {
	astNodes := make([]AST, 0)

	if len(items) < 4 {
		return astNodes, ParseError{"missing first, full match statement"}
	}

	if len(items) > 0 && items[0].Typ != TypeMatch {
		return astNodes, ParseError{"first statement must be a match"}
	}

	if len(items) > 1 && items[1].Typ != TypeLParen {
		return astNodes, ParseError{"missing '(' after match statement"}
	}

	if len(items) > 2 && items[2].Typ != TypeName {
		return astNodes, ParseError{"missing node name for match statement"}
	}

	nodeName := items[2].Val
	var labelFilter = ""

	rparenPos := 3
	if len(items) > 3 && items[3].Typ == TypeColon {
		rparenPos += 2
		if len(items) > 4 && items[4].Typ != TypeName {
			return astNodes, ParseError{"missing label after colon in match statement"}
		}
		labelFilter = items[4].Val
	}
	if len(items) > rparenPos && items[rparenPos].Typ != TypeRParen {
		return astNodes, ParseError{"missing ')' in match statement"}
	}

	astNodes = append(astNodes, Match{
		nodeName:    nodeName,
		labelFilter: labelFilter,
	})

	retPos := rparenPos + 1
	if len(items) <= retPos {
		return astNodes, ParseError{"missing return statement after match"}
	}
	if len(items) > retPos && items[retPos].Typ != TypeReturn {
		return astNodes, ParseError{"statement after match must be a return statement"}
	}
	if len(items) > retPos+1 && items[retPos+1].Typ != TypeName {
		return astNodes, ParseError{"missing node for return statement"}
	}
	if len(items) > retPos+1 && items[retPos+1].Val != nodeName {
		return astNodes, ParseError{"node name in return statement does not match node name in match statement"}
	}
	astNodes = append(astNodes, Return{
		nodeName: nodeName,
	})

	return astNodes, nil
}
