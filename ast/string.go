package ast

import (
	"fmt"
	"strings"
)

func (p *Program) String() string {
	var circs []string

	for _, circ := range p.Circuits {
		circs = append(circs, circ.String())
	}

	return fmt.Sprintf(
		"{%s, includes %s\n%s}",
		p.Name,
		p.Includes,
		strings.Join(circs, "\n"),
	)
}

func (c *Circuit) String() string {
	return fmt.Sprintf(
		`%s (%s) -> (%s) {%s}`,
		c.Name,
		c.Inputs,
		c.Outputs,
		stmts(c.Statements),
	)
}

func (p Parameters) String() string {
	var strs []string

	for _, param := range p {
		if param.Macro {
			strs = append(strs, "%"+param.Name)
		} else {
			strs = append(strs, param.Name)
		}
	}

	return strings.Join(strs, ", ")
}

func exprs(in []Expression) string {
	var strs []string

	for _, x := range in {
		strs = append(strs, x.String())
	}

	return strings.Join(strs, ", ")
}

func stmts(in []Statement) string {
	var strs []string

	for _, x := range in {
		strs = append(strs, x.String())
	}

	return strings.Join(strs, "; ")
}

func (m *MacroStmt) String() string {
	return fmt.Sprintf("<macro %s (%s)>", m.Name, m.Registers.String())
}

func (c *Call) String() string {
	return fmt.Sprintf(
		"<call %s (%s) -> (%s)",
		c.Circuit,
		exprs(c.Inputs),
		c.Outputs.String(),
	)
}

func (p *Pipe) String() string {
	return fmt.Sprintf(
		"<pipe (%s) -> (%s)>",
		exprs(p.Inputs),
		p.Outputs.String(),
	)
}

func (c *Clock) String() string {
	return fmt.Sprintf(
		"<clock %s (%s) [%s]>",
		c.Delay.String(),
		c.Counter,
		stmts(c.Body),
	)
}

func (b *Bit) String() string {
	if b.Value {
		return "<bit 1>"
	}
	return "<bit 0>"
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Infix) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

func (p *Prefix) String() string {
	return fmt.Sprintf("%s%s", p.Operator, p.Right.String())
}

func (m *MacroExpr) String() string {
	return fmt.Sprintf("%%%s", m.Name)
}
