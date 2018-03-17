package token

import "fmt"

// A Type specifies the type of a token.
type Type string

// The set of valid (and invalid!) token types
const (
	_ Type = ""

	EOF     = "EOF"
	Illegal = "illegal"

	Number = "number"
	Ident  = "identifier"

	Prefix = "prefix"
	Infix  = "infix"

	Semi       = ";"
	LeftParen  = "("
	RightParen = ")"
	LeftBrace  = "{"
	RightBrace = "}"
	Comma      = ","
	Macro      = "%"
	Arrow      = "->"
	Colon      = ":"

	Clock   = "clock"
	Name    = "name"
	Circuit = "circuit"
)

// Keywords maps keyword literals to their types.
var Keywords = map[string]Type{
	"clock":   Clock,
	"name":    Name,
	"circuit": Circuit,
}

// IsKeyword checks whether or not a Type is a keyword.
func (t Type) IsKeyword() bool {
	for _, k := range Keywords {
		if k == t {
			return true
		}
	}

	return false
}

// A Token represents a part of the source code, to make it
// easier to parse a program.
type Token struct {
	Type    Type
	Literal string
	Range   Range
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"[%s %d:%d-%d:%d] %s `%s`",
		t.Range.Start.File,
		t.Range.Start.Line,
		t.Range.Start.Col,
		t.Range.End.Line,
		t.Range.End.Col,
		t.Type,
		t.Literal,
	)
}

// A Position specifies a position in the source code.
type Position struct {
	Line, Col int
	File      string
}

// A Range represents a range of characters between two positions.
type Range struct {
	Start, End Position
}
