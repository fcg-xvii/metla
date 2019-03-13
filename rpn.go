package metla

import (
	"fmt"
	"reflect"

	"github.com/golang-collections/collections/stack"
)

func isOperator(val []byte) bool {
	if len(val) == 1 {
		return val[0] == '+' || val[0] == '-' || val[0] == '*' || val[0] == '/' || val[0] == '^' || val[0] == '!' || val[0] == '>' || val[0] == '<'
	} else {
		switch string(val) {
		case "==", ">=", "<=", "!=", "++", "--":
			return true
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

func (s *operator) exec(st *stack.Stack) error {
	if (len(s.data) == 1 && s.data[0] == '!') || (string(s.data) == "++" || string(s.data) == "--") {
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
		if val.Type() != reflect.Bool {
			return s.fatalError("Expected boolean value")
		} else {
			st.Push(&valBoolean{s.rawInfoRecord, !val.Bool()})
		}
	} else {
		if !val.IsNumber() {
			return s.fatalError("Expected number value")
		} else {
			switch string(s.data) {
			case "++":
				st.Push(s.numberResult(val.Float() + 1))
			case "--":
				st.Push(s.numberResult(val.Float() - 1))
			}
		}
	}
	return nil
}

func (s *operator) numberResult(val float64) value {
	if canInt(val) {
		return &valInt{s.rawInfoRecord, int64(val)}
	} else {
		return &valFloat{s.rawInfoRecord, val}
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
	if len(s.data) == 1 {
		if !lVal.IsNumber() {
			return s.fatalError("Left operand must be a number")
		}
		if !rVal.IsNumber() {
			return s.fatalError("Right operand must be a number")
		}
		switch s.data[0] {
		case '+':
			st.Push(s.numberResult(lVal.Float() + rVal.Float()))
		case '-':
			st.Push(s.numberResult(lVal.Float() - rVal.Float()))
		case '*':
			st.Push(s.numberResult(lVal.Float() * rVal.Float()))
		case '/':
			st.Push(s.numberResult(lVal.Float() / rVal.Float()))
		case '%':
			st.Push(&valInt{s.rawInfoRecord, lVal.Int() % rVal.Int()})
		case '>':
			st.Push(&valBoolean{s.rawInfoRecord, lVal.Float() > rVal.Float()})
		case '<':
			st.Push(&valBoolean{s.rawInfoRecord, lVal.Float() < rVal.Float()})
		default:
			return s.fatalError(fmt.Sprintf("Illegal operator '%c'", s.data[0]))
		}
	} else {
		switch string(s.data) {
		case "==":
			if lVal.IsNil() || rVal.IsNil() {
				st.Push(s.checkNil(lVal, rVal))
			} else {
				st.Push(&valBoolean{s.rawInfoRecord, lVal.StaticVal() == rVal.StaticVal()})
			}
		case "!=":
			st.Push(&valBoolean{s.rawInfoRecord, !(lVal.StaticVal() == rVal.StaticVal())})
		case ">=", "<=":
			{
				if !lVal.IsNumber() {
					return s.fatalError("Left operand must be a number")
				}
				if !rVal.IsNumber() {
					return s.fatalError("Right operand must be a number")
				}
				switch s.data[0] {
				case '>':
					st.Push(&valBoolean{s.rawInfoRecord, lVal.Float() >= rVal.Float()})
				case '<':
					st.Push(&valBoolean{s.rawInfoRecord, lVal.Float() <= rVal.Float()})
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
		case '+', '-', '*', '/', '^', '!', '=', '>', '<':
			{
				p.SetupMark()
				op := operator{p.infoRecordFromMark(), []byte{p.Char()}, true}
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
