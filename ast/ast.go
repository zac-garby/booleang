package ast

import "time"

// A Node is the interface from which both expression types
// and statement types extend.
type Node interface {
	String() string
}

// A Statement is a piece of code which doesn't evaluate to
// a value, such as a function call.
type Statement interface {
	Node
	Statement()
}

// An Expression is a piece of code which _does_ evaluate to
// a value, such as an infix expression.
type Expression interface {
	Node
	Expression()
}

// A Parameter is either a macro or an identifier.
type Parameter struct {
	Macro bool
	Name  string
}

// Parameters are a list of strings.
type Parameters []Parameter

// A Circuit is similar to a function in other languages - it
// is a bit of code you can call upon later. A Circuit is neither
// a statement or an expression.
type Circuit struct {
	Name            string
	Inputs, Outputs []string
	Statements      []Statement
}

// A Program is an optionally named sequence of circuit
// definitions. The entrance point is the circuit called
// 'main' - if there isn't one, it won't be able to run.
//
// A Program also contains a list of includes, in order
// of their lexical position.
type Program struct {
	Name     string
	Includes []string
	Circuits []*Circuit
}

type stmt struct{}

func (s *stmt) Statement() {}

type (
	// A MacroStmt statement gives a name to a set of registers.
	// e.g. %num (a0, a1, a2, a3);
	MacroStmt struct {
		*stmt
		Name      string
		Registers Parameters
	}

	// A Call statement calls a circuit.
	// e.g. add (a, b, 0) -> (d, e);
	Call struct {
		*stmt
		Circuit string
		Inputs  []Expression
		Outputs Parameters
	}

	// A Pipe statement pipes expressions into registers.
	// e.g. (0, x) -> (a, b);
	Pipe struct {
		*stmt
		Inputs  []Expression
		Outputs Parameters
	}

	// A Clock executes some statements with a set interval.
	// Doesn't need a semi colon.
	// e.g. clock 1.5s { !a -> a; }
	Clock struct {
		*stmt
		Delay   time.Duration
		Counter string
		Body    []Statement
	}
)

type expr struct{}

func (e *expr) Expression() {}

type (
	// A Bit is a bit literal - either 1 or 0.
	Bit struct {
		*expr
		Value bool
	}

	// An Identifier usually denotes the name of a register.
	// e.g. foobar
	Identifier struct {
		*expr
		Value string
	}

	// An Infix is an infix expression.
	// e.g. a ⊻ b
	Infix struct {
		*expr
		Left, Right Expression
		Operator    string
	}

	// A Prefix is a prefix expression.
	// e.g. !foo
	Prefix struct {
		*expr
		Right    Expression
		Operator string
	}

	// A MacroExpr expands to all the registers inside a macro.
	// e.g. %a
	MacroExpr struct {
		*expr
		Name string
	}
)
