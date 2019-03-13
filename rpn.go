package metla

import (
	"fmt"

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
	//fmt.Println(lVal, rVal)
	if len(s.data) == 1 {
		var lNum, rNum valueNumber
		if lNum, check = lVal.(valueNumber); !lVal.IsNumber() || !check {
			return s.fatalError("Left operand must be a number")
		}
		if rNum, check = rVal.(valueNumber); !lVal.IsNumber() || !check {
			return s.fatalError("Right operand must be a number")
		}
		switch s.data[0] {
		case '+':
			st.Push(s.numberResult(lNum.Float() + rNum.Float()))
		case '-':
			st.Push(s.numberResult(lNum.Float() - rNum.Float()))
		case '*':
			st.Push(s.numberResult(lNum.Float() * rNum.Float()))
		case '/':
			st.Push(s.numberResult(lNum.Float() / rNum.Float()))
		case '%':
			st.Push(&valInt{s.rawInfoRecord, lNum.Int() % rNum.Int()})
		case '>':
			st.Push(&valBoolean{s.rawInfoRecord, lNum.Float() > rNum.Float()})
		case '<':
			st.Push(&valBoolean{s.rawInfoRecord, lNum.Float() < rNum.Float()})
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
				var lNum, rNum valueNumber
				if lNum, check = lVal.(valueNumber); !check || !lVal.IsNumber() {
					return s.fatalError("Left operand must be a number")
				}
				if rNum, check = rVal.(valueNumber); !check || !rVal.IsNumber() {
					return s.fatalError("Right operand must be a number")
				}
				switch s.data[0] {
				case '>':
					st.Push(&valBoolean{s.rawInfoRecord, lNum.Float() >= rNum.Float()})
				case '<':
					st.Push(&valBoolean{s.rawInfoRecord, lNum.Float() <= rNum.Float()})
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
	if p.codeStack.Len() > 0 {
		pn = append(pn, p.codeStack.Pop())
		prevVal = true
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
				var val token
				if val, err = initVal(p); err != nil {
					return
				}
				pn, prevVal = append(pn, val), true
			}
		}
	}
	for sPn.Len() > 0 {
		pn = append(pn, sPn.Pop())
	}
	return
}

func checkSimple(op *operator, st *stack.Stack) error {
	if st.Len() == 0 {
		return op.fatalError("Empty args list")
	}
	if right, check := st.Peek().(value); check && right.IsStatic() {
		if op.isUnary() {
			return op.exec(st)
		} else {
			st.Pop()
			if left, check := st.Peek().(value); !check || !left.IsStatic() {
				st.Push(right)
				st.Push(op)
				return nil
			} else {
				st.Push(right)
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
