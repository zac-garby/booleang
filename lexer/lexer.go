package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/zac-garby/booleang/token"
)

// New takes a string and returns a stream of tokens in the form of
// a generator closure function.
func New(str, file string) func() token.Token {
	// Add a newline at the end of the file to prevent errors
	str += "\n"

	var (
		index = 0
		col   = 1
		line  = 1
		ch    = make(chan token.Token)
	)

	go func() {
		for {
			if index < len(str) {
				foundSpace := false

				for index < len(str) && (unicode.IsSpace(rune(str[index])) || str[index] == '#') {
					if unicode.IsSpace(rune(str[index])) {
						index++
						col++

						if str[index-1] == '\n' {
							col = 1
							line++
						}

						foundSpace = true
					} else {
						for index < len(str) && str[index] != '\n' {
							index++
						}

						col = 1
					}
				}

				if foundSpace {
					continue
				}

				found := false

				remainingSubstring := str[index:]

				for _, pair := range lexemes {
					var (
						regex   = pair.regex
						handler = pair.handler
						pattern = regexp.MustCompile(regex)
						match   = pattern.FindStringSubmatch(remainingSubstring)
					)

					if len(match) > 0 {
						found = true
						t, literal, whole := handler(match)
						l := len(whole)

						ch <- token.Token{
							Type:    t,
							Literal: literal,
							Range: token.Range{
								Start: token.Position{Line: line, Col: col, File: file},
								End:   token.Position{Line: line, Col: col + l - 1, File: file},
							},
						}

						index += l
						col += l

						for index < len(str) && unicode.IsSpace(rune(str[index])) && str[index] != '\n' {
							index++
							col++
						}

						if index < len(str) && str[index] == '#' {
							for index < len(str) && str[index] != '\n' {
								index++
							}
						}

						break
					}
				}

				if !found {
					ch <- token.Token{
						Type:    token.Illegal,
						Literal: string(str[index]),
						Range: token.Range{
							Start: token.Position{Line: line, Col: col, File: file},
							End:   token.Position{Line: line, Col: col, File: file},
						},
					}

					index++
					col++
				}
			} else {
				index++
				col++

				ch <- token.Token{
					Type:    token.EOF,
					Literal: "",
					Range: token.Range{
						Start: token.Position{Line: line, Col: col, File: file},
						End:   token.Position{Line: line, Col: col, File: file},
					},
				}
			}
		}
	}()

	return func() token.Token {
		return <-ch
	}
}

type transformer func(token.Type, string, string) (token.Type, string, string)
type handler func([]string) (token.Type, string, string)

func h(t token.Type, group int, transformer transformer) handler {
	return func(m []string) (token.Type, string, string) {
		return transformer(t, m[group], m[0])
	}
}

func none(t token.Type, literal, whole string) (token.Type, string, string) {
	return t, literal, whole
}

func stringTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	escapes := map[string]string{
		`\n`: "\n",
		`\"`: "\"",
		`\'`: "'",
		`\a`: "\a",
		`\b`: "\b",
		`\f`: "\f",
		`\r`: "\r",
		`\t`: "\t",
		`\v`: "\v",
	}

	for k, v := range escapes {
		literal = strings.Replace(literal, k, v, -1)
	}

	return t, literal, whole
}

func idTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	if kwType, ok := token.Keywords[literal]; ok {
		return kwType, literal, whole
	}

	return t, literal, whole
}

type lexicalPair struct {
	regex   string
	handler handler
}

var lexemes = []lexicalPair{
	// literals
	{`^[-+]?\d+(?:\.\d+)?`, h(token.Number, 0, none)},
	{`^"((\\"|[^"])*)"`, h(token.String, 1, stringTransformer)},
	{`^'((\\'|[^'])*)'`, h(token.String, 1, stringTransformer)},
	{`^[\p{L}\p{M}_][\p{L}\p{M}\d_!?]*`, h(token.Ident, 0, idTransformer)},

	// punctuation
	{`^;`, h(token.Semi, 0, none)},
	{`^\(`, h(token.LeftParen, 0, none)},
	{`^\)`, h(token.RightParen, 0, none)},
	{`^\{`, h(token.LeftBrace, 0, none)},
	{`^\}`, h(token.RightBrace, 0, none)},
	{`^\,`, h(token.Comma, 0, none)},
	{`^%`, h(token.Macro, 0, none)},
	{`^->`, h(token.Arrow, 0, none)},
	{`^:`, h(token.Colon, 0, none)},

	// prefix operators
	{`^!`, h(token.Prefix, 0, none)},
	{`^¬`, h(token.Prefix, 0, none)},

	// infix operators
	{`^&`, h(token.Infix, 0, none)},
	{`^\|`, h(token.Infix, 0, none)},
	{`^\^`, h(token.Infix, 0, none)},
	{`^∧`, h(token.Infix, 0, none)},
	{`^∨`, h(token.Infix, 0, none)},
	{`^⊻`, h(token.Infix, 0, none)},
}
