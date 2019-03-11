package metla

import (
	"errors"
	"fmt"
	"reflect"
	"io"
)

var (
	lNumError = errors.New("left side is not number")
	rNumError = errors.New("right side is not number")
)

func arifmeticResultInt(l, r interface{}, op byte) (res interface{}, err error) {
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
		}
	}
	return
}

func initNestedArifmetic(p *parser) (res *arifmetic, err error) {
	fmt.Println("NESTED_ARIFMETIC")
	var (
		operands         []token
		operators        []byte
		operand          token
		operatorExpected bool
	)
	p.SetupMark()
	for p.Char() != ')' {
		p.PassSpaces()
		if operatorExpected {
			//fmt.Println("nested operator expected")
			if !isOpArifmetic(p.Char()) {
				err = p.positionError(fmt.Sprintf("Unexpected operator character '%c'", p.Char()))
				return
			}
			operators = append(operators, p.Char())
			operatorExpected = false
			p.IncPos()
		} else {
			//fmt.Println("nested operand expected", p.MarkLine(), p.MarkLinePos())
			if p.Char() == '(' {
				p.IncPos()
				if operand, err = initNestedArifmetic(p); err == nil {
					operands = append(operands, operand)
				} else {
					return
				}
			} else {
				if operand, err = initVal(p); err == nil {
					operands = append(operands, operand)
				} else {
					return
				}
			}
			operatorExpected = true
		}
	}
	p.IncPos()
	res = &arifmetic{
		p.infoRecordFromMark(),
		operands,
		operators,
	}
	fmt.Println(res)
	return
}

func initArifmetic(p *parser) (res *arifmetic, err error) {
	p.SetupMark()
	var (
		bracketOpened, operatorExpected bool
		operators                       []byte
		operands                        []token
		operand                         token
	)
	if bracketOpened = p.Char() == '('; bracketOpened {
		p.IncPos()
		if operand, err = initNestedArifmetic(p); err != nil {
			return
		} else {
			operands = append(operands, operand)
		}
	} else {
		if p.codeStack.Len() == 0 {
			err = p.positionError("Arifmetic parse error, operator left side is empty")
			return
		}
		operands = append(operands, p.codeStack.Pop().(token))
	}
	operatorExpected = true
	for !p.IsEndLine() && p.Char() != ',' {
		fmt.Println(operatorExpected, p.MarkLine(), p.MarkLinePos())
		p.PassSpaces()
		if operatorExpected {
			if !isOpArifmetic(p.Char()) {
				err = fmt.Errorf("Unexpected character '%c', expected arifmetic operator", p.Char())
				return
			} else {
				operators = append(operators, p.Char())
				p.IncPos()
				operatorExpected = false
			}
		} else {
			if p.Char() == '(' {
				p.IncPos()
				if operand, err = initNestedArifmetic(p); err != nil {
					return
				} else {
					operands = append(operands, operand)
				}
			} else {
				if operand, err = initVal(p); err != nil {
					return
				} else {
					operands = append(operands, operand)
					p.PassSpaces()
					if bracketOpened {
						if p.Char() == ')' {
							p.IncPos()
							break
						}
					} else if p.IsEndLine() || p.Char() == ',' {
						break
					}
				}
			}
			operatorExpected = true
		}
	}
	fmt.Println(operands)
	fmt.Println(operators)
	res = &arifmetic{
		p.infoRecordFromMark(),
		operands,
		operators,
	}
	return
}

type arifmetic struct {
	*rawInfoRecord
	tokens    []token
	operators []byte
}

func (s *arifmetic) IsExecutable() bool { return false }
func (s *arifmetic) String() string     { return "[arifmetic...]" }

func (s *arifmetic) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
	var (
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
	return nil, nil
}

type arifmeticExec struct {
	*rawInfoRecord
	operands  []execObject
	operators []byte
}

func (s *arifmeticExec) Data(w io.Writer) (err error) {
	var res interface{}
	res, err = s.Val(); err == nil {
		err = w.Write([]byte(fmt.Sprint(res)))
	}
	return
}

func (s *arifmeticExec) Val() (res interface{}, err error) {
	for i, v := range s.operands {
		if 
	}
}

func (s *arifmeticExec) Vals() (res []interface{}, err error) {
	if val, err = s.Val(); err == nil {
		res = []interface{val}
	}
	return
}

func (s *arifmeticExec) ValSingle() bool { return true }
