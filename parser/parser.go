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

func (p *Parser) Parse() (prog *ast.Program, err error) {
	prog = p.parse()

	if len(p.Errors) > 0 {
		return nil, p.Errors[0]
	}

	return prog, nil
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

		if !p.expect(token.Semi) {
			return nil
		}
	}

	for !p.peekIs(token.EOF) {
		p.next()

		if p.cur.Type == token.Circuit {
			circuit := p.parseCircuit()
			if circuit == nil {
				return nil
			}

			prog.Circuits = append(prog.Circuits, circuit)
		} else if p.cur.Type == token.Include {
			byName := false

			if p.peekIs(token.Name) {
				p.next()
				byName = true
			}

			path, ok := p.parseInclude()
			if !ok {
				return nil
			}

			prog.Includes = append(prog.Includes, ast.Include{
				ByName: byName,
				Value:  path,
			})
		} else {
			p.curErr("only circuits and include statements can be written in the top-level of a file")
			return nil
		}
	}

	return prog
}

func (p *Parser) parseCircuit() *ast.Circuit {
	if !p.expect(token.Ident) {
		return nil
	}

	circ := &ast.Circuit{
		Name: p.cur.Literal,
	}

	if p.peekIs(token.LeftParen) {
		p.next()
		circ.Inputs = p.parseIdents(token.RightParen)

		if !p.expect(token.Arrow) {
			return nil
		}
		if !p.expect(token.LeftParen) {
			return nil
		}

		circ.Outputs = p.parseIdents(token.RightParen)
	}

	if !p.expect(token.LeftBrace) {
		return nil
	}

	circ.Statements = p.parseStatements()

	return circ
}

func (p *Parser) parseInclude() (path string, ok bool) {
	if !p.expect(token.String) {
		return "", false
	}

	path = p.cur.Literal

	if !p.expect(token.Semi) {
		return "", false
	}

	return path, true
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

	case token.Macro:
		if !p.expect(token.Ident) {
			return nil
		}
		left = &ast.MacroExpr{
			Name: p.cur.Literal,
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.cur.Type {
	case token.Macro:
		stmt := &ast.MacroStmt{}
		if !p.expect(token.Ident) {
			return nil
		}
		stmt.Name = p.cur.Literal

		if !p.expect(token.LeftParen) {
			return nil
		}
		stmt.Registers = p.parseParams(token.RightParen)

		if !p.expect(token.Semi) {
			return nil
		}

		return stmt

	case token.Clock:
		stmt := &ast.Clock{}
		if !p.expect(token.Number) {
			return nil
		}

		delay := p.parseDuration()
		if delay == nil {
			return nil
		}
		stmt.Delay = *delay

		if p.peekIs(token.Macro) {
			p.next()

			if !p.expect(token.Ident) {
				return nil
			}

			stmt.Counter = p.cur.Literal
		}

		if !p.expect(token.LeftBrace) {
			return nil
		}

		stmt.Body = p.parseStatements()

		return stmt

	case token.Ident:
		stmt := &ast.Call{
			Circuit: p.cur.Literal,
		}

		if !p.expect(token.LeftParen) {
			return nil
		}
		stmt.Inputs = p.parseExprs(token.RightParen)

		if !p.peekIs(token.Arrow) {
			if !p.expect(token.Semi) {
				return nil
			}

			return stmt
		}
		p.next()

		if !p.expect(token.LeftParen) {
			return nil
		}
		stmt.Outputs = p.parseParams(token.RightParen)

		if !p.expect(token.Semi) {
			return nil
		}

		return stmt

	case token.Number, token.Prefix:
		stmt := &ast.Pipe{
			Inputs: []ast.Expression{
				p.parseExpression(),
			},
		}

		if !p.expect(token.Arrow) {
			return nil
		}
		p.next()

		param := p.parseParam()
		if param == nil {
			return nil
		}

		stmt.Outputs = ast.Parameters{
			*param,
		}

		if !p.expect(token.Semi) {
			return nil
		}

		return stmt

	case token.LeftParen:
		stmt := &ast.Pipe{
			Inputs: p.parseExprs(token.RightParen),
		}

		if !p.expect(token.Arrow) {
			return nil
		}

		if !p.expect(token.LeftParen) {
			return nil
		}

		stmt.Outputs = p.parseParams(token.RightParen)

		if !p.expect(token.Semi) {
			return nil
		}

		return stmt

	default:
		p.curErr("unexpected token '%s' at the start of a statement", p.cur.Type)
		return nil
	}
}
