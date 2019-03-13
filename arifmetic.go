package metla

import (
	"fmt"
	"io"
)

/*var (
	lNumError = errors.New("left side is not number")
	rNumError = errors.New("right side is not number")
)*/

func initArifmetic(p *parser) (res token, err error) {
	fmt.Println("INIT_ART")
	var pn []interface{}
	if pn, err = parseRPN(p); err != nil {
		return
	}
	fmt.Println(pn)
	fmt.Println("EXEC_RPN")
	fmt.Println(execRPN(pn))
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
