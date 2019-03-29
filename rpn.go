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

func (s reflectNum) IsNumber() bool {
	switch s.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
		return true
	}
	return false
}

func (s reflectNum) Add(val int64) interface{} {
	switch s.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return s.Value.Int() + val
	case reflect.Float64, reflect.Float32:
		return s.Value.Float() + float64(val)
	default:
		return 0
	}
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

func (s reflectNum) IsNil() bool {
	switch s.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Func, reflect.Interface, reflect.Map, reflect.UnsafePointer, reflect.Slice:
		return s.Value.IsNil()
	default:
		return false
	}
}

func isArifmeticSymbol(ch byte) bool {
	return isOperatorSymbol(ch) || lineman.CheckLetter(ch) || lineman.CheckNumber(ch) || ch == '(' || ch == ')'
}

func isOperatorSymbol(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '!' || c == '=' || c == '>' || c == '<' || c == '%' || c == '&' || c == '|'
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
			case "==", ">=", "<=", "!=", "++", "--", "&&", "||":
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
	iface := st.Pop()
	if vlr, check := iface.(valuer); check {
		iface = vlr.Value()
	}
	val := reflectNum{reflect.ValueOf(iface)}
	if len(s.data) == 1 {
		if val.Kind() != reflect.Bool {
			return fmt.Errorf("Boolean value expected, [%v] given", val.Kind())
		}
		st.Push(!val.Bool())
	} else {
		if !val.IsNumber() {
			return fmt.Errorf("Number value expected, [%v] given", val.Kind())
		}
		switch string(s.data) {
		case "++":
			st.Push(val.Add(1))
		case "--":
			st.Push(val.Add(-1))
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

func (s *operator) valFromStack(st *stack.Stack) reflectNum {
	val := st.Pop()
	if vVal, check := val.(*variable); check {
		return reflectNum{reflect.ValueOf(vVal.value)}
	}
	return reflectNum{reflect.ValueOf(val)}
}

func (s *operator) valsFromStack(st *stack.Stack) (l, r reflectNum, err error) {
	if st.Len() < 2 {
		err = fmt.Errorf("Operands less then 2")
	} else {

		r = s.valFromStack(st)
		l = s.valFromStack(st)
		if !l.IsNil() && !r.IsNil() && (l.Kind() != r.Kind()) {
			lt := l.Type()
			if !r.Type().ConvertibleTo(lt) {
				err = fmt.Errorf("Coudn't convert type [%v] to [%v]", lt, r.Type())
				return
			}
			r = reflectNum{r.Convert(lt)}
		}
	}
	return
}

func (s *operator) execBinary(st *stack.Stack) error {
	l, r, err := s.valsFromStack(st)
	if err != nil {
		return err
	}
	if len(s.data) == 1 {
		switch s.data[0] {
		case '+':
			st.Push(s.numberResult(l.Float() + r.Float()))
		case '-':
			st.Push(s.numberResult(l.Float() - r.Float()))
		case '*':
			st.Push(s.numberResult(l.Float() * r.Float()))
		case '/':
			if r.Float() == 0 {
				return fmt.Errorf("Division by zero")
			}
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
			if l.IsNil() || r.IsNil() {
				st.Push(l.IsNil() && r.IsNil())
			} else {
				st.Push(l.Interface() == r.Interface())
			}
		case "!=":
			if l.IsNil() || r.IsNil() {
				st.Push(!(l.IsNil() && r.IsNil()))
			} else {
				st.Push(!(l.Interface() == r.Interface()))
			}
		case "&&", "||":
			if l.Kind() != reflect.Bool || r.Kind() != reflect.Bool {
				return fmt.Errorf("Expected boolean values")
			}
			if string(s.data) == "&&" {
				st.Push(l.Bool() && r.Bool())
			} else {
				st.Push(l.Bool() || r.Bool())
			}
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
	//fmt.Println("PARSE_RPN", p.stack.Len(), p.stack.Peek())
	prevVal := false
	sPn := stack.New()
	/*if p.stack.Len() > 0 {
		for p.stack.Len() > 0 {
			if _, check := p.stack.Peek().(splitter); !check {
				pn = append(pn, p.stack.Pop())
				prevVal = true
			} else {
				break
			}
		}
	}*/
	if p.Char() != '(' && p.Char() != '!' {
		//fmt.Println("RRRRRRRR", p.readStackVal())
		pn = append(pn, p.readStackVal()...)
	}

	bracketOpened := 0
	p.PassSpaces()
loop:
	for !p.IsEndLine() && isArifmeticSymbol(p.Char()) {
		switch p.Char() {
		case '(':
			{
				bracketOpened++
				sPn.Push(byte('('))
				p.IncPos()
			}
		case ')':
			{
				if bracketOpened == 0 {
					break loop
				}
				bracketOpened--
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
		case '+', '-', '*', '/', '^', '!', '=', '>', '<', '%', '&', '|':
			{
				p.SetupMark()
				op := operator{p.infoRecordFromMark(), []byte{p.Char()}, true}
				if checkOp := []byte{p.Char(), p.NextChar()}; isOperator(checkOp) {
					op.data = checkOp
					p.IncPos()
					if !prevVal {
						op.postfix = false
					}
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
				fmt.Println("VAL.....")
				if _, errCV := initCodeVal(p); errCV != nil {
					return
				} else {
					pn, prevVal = append(pn, p.readStackVal()...), true
				}
			}
		}
		p.PassSpaces()
	}
	for sPn.Len() > 0 {
		pn = append(pn, sPn.Pop())
	}
	//fmt.Println(pn, len(pn), cap(pn))
	return
}

func checkSimple(op *operator, st *stack.Stack) error {
	if st.Len() == 0 {
		return op.fatalError("Empty args list")
	}
	if isStaticDataObject(st.Peek()) {
		if op.isUnary() {
			return op.exec(st)
		} else {
			r := st.Pop()
			if !isStaticDataObject(st.Peek()) {
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
	/*fmt.Println("SIMPLE_RPN")
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
	fmt.Println("SIMPLE", res)*/
	return pl, nil
}
