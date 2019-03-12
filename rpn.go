package metla

import (
	"fmt"
	"reflect"

	"github.com/golang-collections/collections/stack"
)

func isOperator(val []byte) bool {
	if len(val) == 1 {
		return val[0] == '+' || val[0] == '-' || val[0] == '*' || val[0] == '/' || val[0] == '^' || val[0] == '!'
	} else {
		switch string(val) {
		case "==", ">=", "<=", "!=":
			return true
		}
	}
	return false
}

func opPriority(op *operator) byte {
	if len(op.data) == 1 {
		switch op.data[0] {
		case '^', '!':
			return 3
		case '*', '/':
			return 2
		case '+', '-':
			return 1
		}
	}
	return 0
}

type operator struct {
	*rawInfoRecord
	data []byte
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

func (s *operator) exec(st *stack.Stack) error {
	if len(s.data) == 1 && s.data[0] == '!' {
		return s.execUnary(st)
	} else {
		return s.execBinary(st)
	}
}

func (s *operator) execUnary(st *stack.Stack) error {
	if st.Len() == 0 {
		return s.fatalError("Value list is empty")
	}
	v := st.Pop()
	if val, check := v.(value); !check {
		return s.fatalError("Expected value object")
	} else if val.Type() != reflect.Bool {
		return val.fatalError("Expected boolean value")
	} else {
		st.Push(&valBoolean{s.rawInfoRecord, !val.Bool()})
		return nil
	}
}

func (s *operator) execBinary(st *stack.Stack) error {
	if st.Len() < 2 {
		return s.fatalError("Value list is too small")
	}
	r, l, check := st.Pop(), st.Pop(), false
	var lVal, rVal value
	if lVal, check = l.(value); !check {
		return s.fatalError("Expected value object left")
	}
	if rVal, check = r.(value); !check {
		return s.fatalError("Expected value object right")
	}
	fmt.Println(lVal, rVal)
	if lVal.IsNil() || rVal.IsNil() {

	}
	return fmt.Errorf("TEST...")
}

func (s *operator) checkNil(l, r value) (value, error) {

}

func parseRPN(p *parser) (pn []interface{}, err error) {
	sPn := stack.New()
	if p.Char() != '(' && p.Char() != '!' {
		if p.codeStack.Len() == 0 {
			err = p.positionError("Arifmetic left side not found")
			return
		}
		pn = append(pn, p.codeStack.Pop())
	}
	for !p.IsEndLine() && p.Char() != ',' {
		p.PassSpaces()
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
		case '+', '-', '*', '/', '^', '!', '=':
			{
				p.SetupMark()
				op := operator{p.infoRecordFromMark(), []byte{p.Char()}}
				if checkOp := []byte{p.Char(), p.NextChar()}; isOperator(checkOp) {
					op.data = checkOp
					p.IncPos()
				} else if !isOperator(op.data) {
					err = p.positionError(fmt.Sprintf("Unexpected operator '%c'", p.Char()))
					return
				}
				for sPn.Len() > 0 {
					fmt.Println(sPn.Peek())
					if val, check := sPn.Peek().(*operator); check && opPriority(val) >= opPriority(&op) {
						pn = append(pn, sPn.Pop())
					} else {
						break
					}
				}
				sPn.Push(&op)
				p.IncPos()
			}
		default:
			{
				var val token
				if val, err = initVal(p); err != nil {
					return
				}
				pn = append(pn, val)
			}
		}
	}
	for sPn.Len() > 0 {
		pn = append(pn, sPn.Pop())
	}
	return
}

func simpleRPN(pl []interface{}) (res []interface{}, err error) {
	return
}

func execRPN(pl []interface{}) (res interface{}, err error) {
	st := stack.New()
	for _, v := range pl {
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
