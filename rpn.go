package metla

import (
	"fmt"
	"reflect"

	"github.com/fcg-xvii/lineman"
	"github.com/golang-collections/collections/stack"
)

type reflectNum struct {
	reflect.Value
}

func (s reflectNum) Int() (res int64) {
	switch s.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return s.Value.Int()
	case reflect.Float64, reflect.Float32:
		return int64(s.Value.Float())
	default:
		return 0
	}
}

func (s reflectNum) Float() (res float64) {
	switch s.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return float64(s.Value.Int())
	case reflect.Float64, reflect.Float32:
		return s.Value.Float()
	default:
		return 0
	}
}

func isArifmeticSymbol(ch byte) bool {
	return isOperatorSymbol(ch) || lineman.CheckLetter(ch) || lineman.CheckNumber(ch) || ch == '(' || ch == ')'
}

func isOperatorSymbol(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '!' || c == '=' || c == '>' || c == '<' || c == '%'
}

func isOperator(val []byte) bool {
	switch len(val) {
	case 0:
		return false
	case 1:
		return val[0] == '+' || val[0] == '-' || val[0] == '*' || val[0] == '/' || val[0] == '^' || val[0] == '!' || val[0] == '>' || val[0] == '<' || val[0] == '%'
	default:
		{
			switch string(val[:2]) {
			case "==", ">=", "<=", "!=", "++", "--":
				return true
			}
		}
	}
	return false
}

func opPriority(op *operator) byte {
	if len(op.data) == 1 {
		switch op.data[0] {
		case '^', '!':
			return 4
		case '*', '/':
			return 3
		case '+', '-':
			return 2
		}
	} else if string(op.data) == "++" || string(op.data) == "--" {
		if op.postfix {
			return 0
		} else {
			return 4
		}
	}
	return 1
}

type operator struct {
	*rawInfoRecord
	data    []byte
	postfix bool
}

func (s *operator) String() string {
	if len(s.data) == 1 {
		switch s.data[0] {
		case 43:
			return "+"
		case 45:
			return "-"
		case 42:
			return "*"
		case 47:
			return "/"
		case 94:
			return "^"
		case 62:
			return ">"
		case 60:
			return "<"
		case 33:
			return "!"
		}
	}
	return string(s.data)
}

func (s *operator) isUnary() bool {
	return (len(s.data) == 1 && s.data[0] == '!') || (string(s.data) == "++" || string(s.data) == "--")
}

func (s *operator) exec(st *stack.Stack) error {
	if s.isUnary() {
		return s.execUnary(st)
	} else {
		return s.execBinary(st)
	}
}

func (s *operator) execUnary(st *stack.Stack) error {
	if st.Len() == 0 {
		return s.fatalError("Value list is empty")
	}
	if val, check := st.Pop().(value); !check {
		return s.fatalError("Expected value object")
	} else if len(s.data) == 1 {
		if bVal, check := val.(valueBoolean); !check {
			return s.fatalError("Expected boolean value")
		} else {
			st.Push(&valBoolean{s.rawInfoRecord, !bVal.Bool()})
		}
	} else {
		if numVal, check := val.(valueNumber); !val.IsNumber() || !check {
			return s.fatalError("Expected number value")
		} else {
			switch string(s.data) {
			case "++":
				numVal.Add(1)
			case "--":
				numVal.Add(-1)
			}
			st.Push(numVal)
		}
	}
	return nil
}

func (s *operator) numberResult(val float64) interface{} {
	if canInt(val) {
		return int64(val)
	} else {
		return val
	}
}

func (s *operator) valsFromStack(st *stack.Stack) (l, r reflectNum, err error) {
	if st.Len() < 2 {
		err = fmt.Errorf("Operands less then 2")
	} else {
		r = reflectNum{reflect.ValueOf(st.Pop())}
		l = reflectNum{reflect.ValueOf(st.Pop())}
	}
	return
}

func (s *operator) execBinary(st *stack.Stack) error {
	l, r, err := s.valsFromStack(st)
	if err != nil {
		return err
	}
	fmt.Println(r, l, st.Len())
	if len(s.data) == 1 {
		switch s.data[0] {
		case '+':
			st.Push(s.numberResult(l.Float() + r.Float()))
		case '-':
			st.Push(s.numberResult(l.Float() - r.Float()))
		case '*':
			st.Push(s.numberResult(l.Float() * r.Float()))
		case '/':
			st.Push(s.numberResult(l.Float() / r.Float()))
		case '%':
			st.Push(int64(l.Int() % r.Int()))
		case '>':
			st.Push(l.Float() > r.Float())
		case '<':
			st.Push(l.Float() < r.Float())
		default:
			return s.fatalError(fmt.Sprintf("Illegal operator '%c'", s.data[0]))
		}
	} else {
		switch string(s.data) {
		case "==":
			st.Push(l.Interface() == r.Interface())
			/*if l.IsNil() || r.IsNil() {
				st.Push(l.Interface() == r.Interface())
			} else {
				st.Push(l.Interface() == r.Interface())
			}*/
		case "!=":
			st.Push(!(l.Interface() == r.Interface()))
		case ">=", "<=":
			{
				switch s.data[0] {
				case '>':
					st.Push(l.Float() >= r.Float())
				case '<':
					st.Push(l.Float() <= r.Float())
				}
			}
		default:
			return s.fatalError(fmt.Sprintf("Illegal operator '%c'", s.data[0]))
		}
	}
	return nil
}

func (s *operator) checkNil(l, r value) (res *valBoolean) {
	res = &valBoolean{l.posInfo(), false}
	if l.IsNil() == r.IsNil() {
		res.val = true
	}
	return
}

func parseRPN(p *parser) (pn []interface{}, err error) {
	fmt.Println("PARSE_RPN")
	prevVal := false
	sPn := stack.New()
	if p.stack.Len() > 0 {
		pn = append(pn, p.stack.Pop())
		prevVal = true
	}
	p.PassSpaces()
	for !p.IsEndLine() && isArifmeticSymbol(p.Char()) {
		switch p.Char() {
		case '(':
			{
				sPn.Push(byte('('))
				p.IncPos()
			}
		case ')':
			{
				accepted := false
				for sPn.Len() > 0 {
					if c, check := sPn.Peek().(byte); check && c == '(' {
						sPn.Pop()
						accepted = true
						break
					} else {
						pn = append(pn, sPn.Pop())
					}
				}
				if !accepted {
					err = p.positionError("Not closed bracked in arifmetic expression.")
					return
				}
				p.IncPos()
			}
		case '+', '-', '*', '/', '^', '!', '=', '>', '<', '%':
			{
				p.SetupMark()
				op := operator{p.infoRecordFromMark(), []byte{p.Char()}, true}
				if checkOp := []byte{p.Char(), p.NextChar()}; isOperator(checkOp) {
					op.data = checkOp
					p.IncPos()
					if !prevVal {
						op.postfix = false
					}
					fmt.Println("POSTFIX", op.postfix)
				} else if !isOperator(op.data) {
					err = p.positionError(fmt.Sprintf("Unexpected operator '%c'", p.Char()))
					return
				}
				for sPn.Len() > 0 {
					if val, check := sPn.Peek().(*operator); check && opPriority(val) >= opPriority(&op) {
						pn = append(pn, sPn.Pop())
					} else {
						break
					}
				}
				sPn.Push(&op)
				p.IncPos()
				prevVal = false
			}
		default:
			{

				if _, err = initCodeVal(p); err != nil {
					return
				} else {
					pn, prevVal = append(pn, p.stack.Pop()), true
				}
			}
		}
		p.PassSpaces()
	}
	for sPn.Len() > 0 {
		pn = append(pn, sPn.Pop())
	}
	fmt.Println(pn, len(pn), cap(pn))
	return
}

func checkSimple(op *operator, st *stack.Stack) error {
	fmt.Println("SIMPLE_RPN")
	if st.Len() == 0 {
		return op.fatalError("Empty args list")
	}
	if _, rCheck := st.Peek().(*valVariable); !rCheck {
		if op.isUnary() {
			return op.exec(st)
		} else {
			fmt.Println("PPPPP")
			r := st.Pop()
			if _, lCheck := st.Peek().(*valVariable); !lCheck || !rCheck {
				fmt.Println("LCH", lCheck, rCheck)
				st.Push(r)
				st.Push(op)
				return nil
			} else {
				st.Push(r)
				return op.exec(st)
			}
		}
	}
	st.Push(op)
	return nil
}

func simpleRPN(pl []interface{}) (res []interface{}, err error) {
	st := stack.New()
	for _, v := range pl {
		if op, check := v.(*operator); check {
			if err = checkSimple(op, st); err != nil {
				return
			}
		} else {
			st.Push(v)
		}
	}
	res = make([]interface{}, st.Len())
	for i := len(res) - 1; i >= 0; i-- {
		res[i] = st.Pop()
	}
	return
}

func execRPN(pn []interface{}) (res interface{}, err error) {
	fmt.Println("PNNNNN", pn)
	for _, v := range pn {
		fmt.Printf("%T ::\n", v)
	}
	st := stack.New()
	for _, v := range pn {
		if op, check := v.(*operator); check {
			if err = op.exec(st); err != nil {
				return
			}
		} else {
			st.Push(v)
		}
	}
	res = st.Pop()
	return
}
