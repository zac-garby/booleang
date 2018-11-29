package lexer_test

import (
	"testing"

	. "github.com/zac-garby/booleang/lexer"
	"github.com/zac-garby/booleang/token"
)

func TestLexer(t *testing.T) {
	input := `
        5 -2 +4.33 0.1290; # some numbers
        ident_1 π___05xyz 5a;
        "hello \" world" 'foo \' bar';

        # the prefixes:
        ! ¬

        # the infixes:
        & | ^ ∧ ∨ ⊻

        ; ( ) { } , % -> :

        clock name circuit include

        $ # this token is illegal
    `

	expected := []token.Type{
		token.Number, token.Number, token.Number, token.Number, token.Semi,
		token.Ident, token.Ident, token.Number, token.Ident, token.Semi,
		token.String, token.String, token.Semi,
		token.Prefix, token.Prefix,
		token.Infix, token.Infix, token.Infix, token.Infix, token.Infix, token.Infix,
		token.Semi, token.LeftParen, token.RightParen, token.LeftBrace, token.RightBrace,
		token.Comma, token.Macro, token.Arrow, token.Colon,
		token.Clock, token.Name, token.Circuit, token.Include,
		token.Illegal,
	}

	next := New(input, "test")

	i := 0
	for tok := next(); tok.Type != token.EOF; tok = next() {
		exp := expected[i]
		i++

		if tok.Type != exp {
			t.Errorf("(%v) expected %s, got %s\n", i, exp, tok.Type)
		}

		if tok.Range.Start.File != "test" || tok.Range.End.File != "test" {
			t.Errorf("(%v) reported wrong file name", i)
		}
	}
}
