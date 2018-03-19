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

func (p *Parser) parseExpression() ast.Expression {
	var left ast.Expression

	switch p.cur.Type {
	case token.Ident:
		left = &ast.Identifier{
			Value: p.cur.Literal,
		}

	case token.Number:
		if !(p.cur.Literal == "0" || p.cur.Literal == "1") {
			p.curErr("a bit literal must be 0 or 1")
			return nil
		}

		left = &ast.Bit{
			Value: p.cur.Literal == "1",
		}

	case token.LeftParen:
		p.next()
		left = p.parseExpression()
		if !p.expect(token.RightParen) {
			return nil
		}

	case token.Prefix:
		op := p.cur.Literal
		p.next()
		left = &ast.Prefix{
			Operator: op,
			Right:    p.parseExpression(),
		}
	}

	if p.peekIs(token.Infix) {
		op := p.peek.Literal
		p.next()
		p.next()
		right := p.parseExpression()
		left = &ast.Infix{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}
