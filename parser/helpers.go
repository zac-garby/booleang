package parser

import (
	"math"
	"strconv"
	"time"

	"github.com/Zac-Garby/booleang/token"
)

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lex()

	if p.peek.Type == token.Illegal {
		p.err(
			"illegal token found: `%s`",
			p.peek.Range, p.peek.Literal,
		)
	}
}

func (p *Parser) curIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.cur.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.peek.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) expect(t token.Type) bool {
	if p.peekIs(t) {
		p.next()
		return true
	}

	p.peekErr(t)
	return false
}

func (p *Parser) parseNumberRaw() (float64, error) {
	return strconv.ParseFloat(p.cur.Literal, 64)
}

func (p *Parser) parseDuration() *time.Duration {
	val, err := p.parseNumberRaw()
	if err != nil {
		p.Errors = append(p.Errors, err)
		return nil
	}

	if math.Floor(val) != val {
		p.curErr("clock duration must be an integer. use a smaller time denomination if necessary")
	}

	if !p.expect(token.Ident) {
		return nil
	}

	var dur time.Duration

	switch p.cur.Literal {
	case "ns":
		dur = time.Nanosecond
	case "ms":
		dur = time.Millisecond
	case "s":
		dur = time.Second
	case "m":
		dur = time.Minute
	case "h":
		dur = time.Hour
	default:
		p.curErr("expexcted ns, ms, s, m, or h to complete the clock duration. got %s", t)
		return nil
	}

	dur *= val

	return &dur
}
