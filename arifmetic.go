package metla

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/golang-collections/collections/stack"
)

type operator byte

func (s operator) String() string {
	switch s {
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
	default:
		return ""
	}
}

func opPriority(op operator) byte {
	switch op {
	case '^':
		return 3
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

var (
	lNumError = errors.New("left side is not number")
	rNumError = errors.New("right side is not number")
)

func arifmeticResult(l, r interface{}, op operator) (res interface{}, err error) {
	if check, _ := checkIfaceNumber(l); !check {
		return nil, lNumError
	}
	if check, _ := checkIfaceNumber(r); !check {
		return nil, rNumError
	}
	lVal, rVal := reflect.ValueOf(l), reflect.ValueOf(r)
	if checkKindInt(lVal.Kind()) && checkKindInt(rVal.Kind()) {
		switch op {
		case '+':
			res = lVal.Int() + rVal.Int()
		case '-':
			res = lVal.Int() - rVal.Int()
		case '*':
			res = lVal.Int() * rVal.Int()
		case '/':
			res = lVal.Int() / rVal.Int()
		case '^':
			res = math.Pow(float64(lVal.Int()), float64(rVal.Int()))
		}
	} else {
		if checkKindFloat(lVal.Kind()) {
			rVal = rVal.Convert(reflect.TypeOf(lVal.Float()))
		} else {
			lVal = lVal.Convert(reflect.TypeOf(rVal.Float()))
		}
		switch op {
		case '+':
			res = lVal.Float() + rVal.Float()
		case '-':
			res = lVal.Float() - rVal.Float()
		case '*':
			res = lVal.Float() * rVal.Float()
		case '/':
			res = lVal.Float() / rVal.Float()
		case '^':
			res = math.Pow(lVal.Float(), rVal.Float())
		}
	}
	return
}

func isOperator(val operator) bool {
	return val == '+' || val == '-' || val == '*' || val == '/' || val == '^'
}

func parseRPL(p *parser) (pn []interface{}, err error) {
	sPn := stack.New()
	if p.Char() != '(' {
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
				fmt.Println("OpenBracket", sPn.Peek())
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
		case '+', '-', '*', '/', '^':
			{
				for sPn.Len() > 0 {
					fmt.Println(sPn.Peek())
					if val, check := sPn.Peek().(operator); check && isOperator(val) && opPriority(val) >= opPriority(operator(p.Char())) {
						pn = append(pn, sPn.Pop())
					} else {
						break
					}
				}
				sPn.Push(operator(p.Char()))
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

func simpleRPL(pl []interface{}) (res []interface{}, err error) {
	return
}

func initArifmetic(p *parser) (res *arifmetic, err error) {
	var pn []interface{}
	if pn, err = parseRPL(p); err != nil {
		return
	}
	fmt.Println(pn)
	/*sPn := stack.New()
	for _, v := range pn {
		if op, check := v.(operator); check {
			r, _ := sPn.Pop().(value).Val()
			l, _ := sPn.Pop().(value).Val()
			if rr, rrr := arifmeticResult(l, r, op); rrr == nil {
				switch rr.(type) {
				case int64:
					sPn.Push(&valInt{p.infoRecordFromMark(), rr.(int64)})
				case float64:
					sPn.Push(&valFloat{p.infoRecordFromMark(), rr.(float64)})
				}
			} else {
				err = rrr
				return
			}
			//return
		} else {
			sPn.Push(v)
		}
	}
	fmt.Println(sPn.Pop())
	*/
	err = fmt.Errorf("Test Error")
	return
}

type arifmetic struct {
	*rawInfoRecord
	pn []interface{}
	//tokens    []token
	//operators []byte
}

func (s *arifmetic) IsExecutable() bool { return false }
func (s *arifmetic) String() string     { return "[arifmetic...]" }

func (s *arifmetic) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
	/*var (
		operands  []execObject
		operators []byte
	)
	for i, v := range s.operators {
		if v == '*' || v == '/' {
			operators = append(operands, v)
			operands = append(operands, s.operands[i], s.operands[i+1])
		}
	}
	for i, v := range s.operators {
		if v == '+' || v == '-' {
			operators = append(operands, v)
			operands = append(operands, s.operands[i], s.operands[i+1])
		}
	}
	return nil, nil*/
	return
}

type arifmeticExec struct {
	*rawInfoRecord
	operands  []execObject
	operators []byte
}

func (s *arifmeticExec) Data(w io.Writer) (err error) {
	var res interface{}
	if res, err = s.Val(); err == nil {
		_, err = w.Write([]byte(fmt.Sprint(res)))
	}
	return
}

func (s *arifmeticExec) Val() (res interface{}, err error) {
	/*result := 0
	for i, v := range s.operands {
		if v == '*' || v == '/' {

		}
	}*/
	return
}

/*func (s *arifmeticExec) Vals() (res []interface{}, err error) {
	if val, err = s.Val(); err == nil {
		res = []interface{val}
	}
	return
}*/

func (s *arifmeticExec) ValSingle() bool { return true }
