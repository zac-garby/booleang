package parser

import (
	"fmt"

	"github.com/Zac-Garby/booleang/token"
)

// An Error represents an error encountered while parsing.
type Error struct {
	Message string
	Range   token.Range
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"{* Parse Error @ [%s %d:%d-%d:%d] *} %s",
		e.Range.Start.File,
		e.Range.Start.Line,
		e.Range.Start.Col,
		e.Range.End.Line,
		e.Range.End.Col,
		e.Message,
	)
}

func (p *Parser) err(msg string, r token.Range, format ...interface{}) {
	p.Errors = append(p.Errors, &Error{
		Message: fmt.Sprintf(msg, format...),
		Range:   r,
	})
}

func (p *Parser) curErr(msg string, format ...interface{}) {
	p.err(msg, p.cur.Range, format...)
}

func (p *Parser) peekErr(t token.Type) {
	p.err("expected %s, but got %s", p.peek.Range, t, p.peek.Type)
}

func (p *Parser) unexpectedToken(t token.Type) {
	p.curErr("unexpected token: %s", t)
}
