package parser

import (
	"github.com/Zac-Garby/booleang/ast"
	"github.com/Zac-Garby/booleang/lexer"
	"github.com/Zac-Garby/booleang/token"
)

// A Parser takes a sequence of tokens and constructs an abstract
// syntax tree.
type Parser struct {
	Errors []error

	lex       func() token.Token
	text      string
	cur, peek token.Token
}

// New makes a new `Parser` instance.
func New(text, file string) *Parser {
	p := &Parser{
		lex:    lexer.New(text, file),
		text:   text,
		Errors: make([]error, 0),
	}

	p.next()
	p.next()

	return p
}

func (p *Parser) parse() *ast.Program {
	prog := &ast.Program{
		Name: "unnamed",
	}

	if p.curIs(token.Name) {
		if !p.expect(token.Colon) {
			return nil
		}

		if !p.expect(token.String) {
			return nil
		}

		prog.Name = p.cur.Literal
	}

	return prog
}

func (p *Parser) Parse() (prog *ast.Program, err error) {
	prog = p.parse()

	if len(p.Errors) > 0 {
		return nil, p.Errors[0]
	}

	return prog, nil
}
