package prod

import (
	"fmt"

	"github.com/fcg-xvii/containers"
)

func isOperatorSymbol(p *parser) bool {
	switch p.Char() {
	case '+', '-', '*', '/', '<', '>', '&', '|', '!', '=':
		return true
	}
	return false
}

func isArifmeticSymbol(p *parser) bool {
	if !isOperatorSymbol(p) {
		return p.Char() == '(' || p.Char() == ')'
	}
	return true
}

func initOperator(p *parser) (op operator, err *parseError) {
	op.position, op.source = position{p.tplName, p.Line(), p.LinePos()}, []byte{p.Char()}
	p.IncPos()
	if isOperatorSymbol(p) {
		op.source = append(op.source, p.Char())
		p.IncPos()
		switch op.String() {
		case ">=", "<=":
			op.priority = 3
		case "==":
			op.priority = 4
		case "&&", "||":
			op.priority = 5
		default:
			err = op.parseError(fmt.Sprintf("Unexpected operator '%v'", op))
		}
	} else {
		switch op.source[0] {
		case '!':
			op.priority, op.prefix = 0, true
		case '*', '/':
			op.priority = 1
		case '+', '-':
			op.priority = 2
		case '<', '>':
			op.priority = 3
		default:
			err = op.parseError(fmt.Sprintf("Unexpected operator '%v'", op))
		}
	}
	return
}

type operator struct {
	position
	source   []byte
	priority uint8
	prefix   bool
}

func (s operator) String() string {
	return string(s.source)
}

func (s *opeartor) lrNumber() (left, right float64)

func (s operator) exec(stack *containers.Stack, exec *tplExec) *execError {
	switch string(s.source) {
	case "+", "-", "*", "/":

	}
}

func (s operator) execArifmetic(stack *containers.Stack) *execError {

}
