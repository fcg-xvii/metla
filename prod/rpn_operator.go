package prod

import (
	"fmt"
	"reflect"

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
		case "==", "!=":
			op.priority = 4
		case "&&", "||":
			op.priority = 5
		case "++", "--":
			op.priority = 1
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

func (s operator) ifaceNumber(iface interface{}) (res float64, err *execError) {
	lVal := reflect.ValueOf(iface)
	switch lVal.Kind() {
	case reflect.Float32, reflect.Float64:
		res = lVal.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = float64(lVal.Int())
	default:
		err = s.execError(fmt.Sprintf("Number type expected, not '%v'", lVal.Kind()))
	}
	return
}

func (s operator) lrNumber(r, l interface{}) (left, right float64, err *execError) {
	if left, err = s.ifaceNumber(l); err != nil {
		return
	}
	right, err = s.ifaceNumber(r)
	return
}

func (s operator) lrBool(r, l interface{}) (left, right bool, err *execError) {
	var check bool
	if left, check = l.(bool); !check {
		err = s.execError("Left bool value expected")
		return
	}
	if right, check = r.(bool); !check {
		err = s.execError("Right bool value expected")
		return
	}
	return
}

func (s operator) lrEqual(r, l interface{}) bool {
	lVal, rVal := reflect.ValueOf(l), reflect.ValueOf(r)
	if lVal.Kind() == reflect.Invalid {
		return rVal.Kind() == reflect.Invalid || rVal.IsNil()
	} else if rVal.Kind() == reflect.Invalid {
		return lVal.Kind() == reflect.Invalid || lVal.IsNil()
	} else if lVal.Kind() == rVal.Kind() {
		return l == r
	} else {
		lType, rType := lVal.Type(), rVal.Type()
		if !rType.ConvertibleTo(lType) {
			return false
		}
		rVal = rVal.Convert(lType)
		return lVal.Interface() == rVal.Interface()
	}
}

func (s operator) exec(stack *containers.Stack, exec *tplExec) *execError {
	switch opVal := string(s.source); opVal {
	case "+", "-", "*", "/", ">", "<", ">=", "<=":
		left, right, err := s.lrNumber(stack.Pop(), stack.Pop())
		if err != nil {
			return err
		}
		if len(s.source) == 1 {
			switch s.source[0] {
			case '+':
				stack.Push(left + right)
			case '-':
				stack.Push(left - right)
			case '*':
				stack.Push(left * right)
			case '/':
				stack.Push(left / right)
			case '>':
				stack.Push(left > right)
			case '<':
				stack.Push(left < right)
			}
		} else {
			switch opVal {
			case ">=":
				stack.Push(left >= right)
			case "<=":
				stack.Push(left <= right)
			}
		}
	case "&&", "||":
		left, right, err := s.lrBool(stack.Pop(), stack.Pop())
		if err != nil {
			return err
		}
		switch opVal {
		case "&&":
			stack.Push(left && right)
		case "||":
			stack.Push(left || right)
		}
	case "==":
		stack.Push(s.lrEqual(stack.Pop(), stack.Pop()))
	case "!=":
		stack.Push(!s.lrEqual(stack.Pop(), stack.Pop()))
	case "!":
		if b, check := stack.Pop().(bool); !check {
			return s.execError("Expected boolean value")
		} else {
			stack.Push(!b)
		}
	case "++", "--":
		if num, err := s.ifaceNumber(stack.Pop()); err != nil {
			return err
		} else {
			if opVal == "++" {
				stack.Push(num + 1)
			} else {
				stack.Push(num - 1)
			}
		}
	default:
		return s.execError("Unexpected opeartor " + string(s.source))
	}
	return nil
}

/*func (s *operator) execArifmetic(stack *containers.Stack) *execError {

}*/
