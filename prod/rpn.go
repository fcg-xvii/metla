package prod

import (
	"fmt"

	"github.com/fcg-xvii/containers"
)

func newRPN(p *parser) *parseError {
	//fmt.Println("NEW_RPN", string(p.Char()))
	obj, valAccepted := rpn{position: position{p.tplName, p.Line(), p.LinePos()}}, false
	st, bracketLayout := new(containers.Stack), 0
	if p.Char() != '!' && p.Char() != '(' {
		if p.stack.Len() == 0 {
			return obj.parseError("RPN parse error :: left side is empty")
		} else {
			obj.pn, valAccepted = append(obj.pn, p.stack.Pop()), true
		}
	}
mainLoop:
	for !p.IsEndDocument() {
		fmt.Println("!!!")
		switch {
		case isOperatorSymbol(p):
			//fmt.Println("OPSYMBOL", st)
			if op, err := initOperator(p); err != nil {
				return err
			} else {
				for st.Len() > 0 {
					if sOp, check := st.Peek().(operator); check {
						if sOp.prefix || op.priority >= sOp.priority {
							obj.pn = append(obj.pn, st.Pop())
						} else {
							break
						}
					} else {
						break
					}
				}
				if valAccepted {
					//fmt.Println("ValAccepted")
					valAccepted = false
				}
				st.Push(op)
				opVal := op.String()
				//fmt.Println("OPVAL", opVal)
				if opVal == "++" || opVal == "--" {
					p.PassSpaces()
					if !isOperatorSymbol(p) {
						break mainLoop
					}
				}
			}
		case p.Char() == '(':
			st.Push('(')
			bracketLayout++
			p.IncPos()
		case p.Char() == ')':
			{
				if bracketLayout == 0 {
					return p.initParseError(p.Line(), p.LinePos(), "Unexpected bracket close ')'")
				}
				bracketLayout--
			popLoop:
				for st.Len() > 0 {
					//fmt.Printf("fff %T\n", st.Peek())
					switch iface := st.Pop(); iface.(type) {
					case int32:
						if iface.(int32) == '(' {
							break popLoop
						}
					default:
						obj.pn = append(obj.pn, iface)
					}
				}
				p.IncPos()
				p.PassSpaces()
				if !isArifmeticSymbol(p) {
					break mainLoop
				}
			}
		default:
			if err := p.initCodeVal(); err != nil {
				return err
			}
			obj.pn = append(obj.pn, p.stack.Pop())
			p.PassSpaces()
			valAccepted = true
			if !isArifmeticSymbol(p) {
				break mainLoop
			}
		}
	}
	//fmt.Println("BRACKET_LAYOUT", bracketLayout)
	if bracketLayout != 0 {
		return obj.parseError("Unclosed bracket")
	}
	obj.pn = append(obj.pn, st.PopAllReverse()...)
	//fmt.Println("PNNNN", obj.pn)
	p.stack.Push(obj)
	return nil
}

type rpn struct {
	position
	pn []interface{}
}

func (s rpn) execRPN(exec *tplExec) (interface{}, *execError) {
	st, execStackLen := containers.NewStack(len(s.pn)), exec.stack.Len()
	for _, v := range s.pn {
		switch v.(type) {
		case getter:
			st.Push(v.(getter).get(exec))
		case executer:
			if err := v.(executer).exec(exec); err != nil {
				return nil, err
			}
			if exec.stack.Len()-1 != execStackLen {
				return nil, v.(coordinator).execError("Not one value returned")
			}
			st.Push(exec.stack.Pop().(getter).get(exec))
		case operator:
			if err := v.(operator).exec(st, exec); err != nil {
				return nil, err
			}
		}
	}
	return st.Pop(), nil
}

func (s rpn) exec(exec *tplExec) *execError {
	if val, err := s.execRPN(exec); err == nil {
		exec.stack.Push(static{s.position, val})
		return nil
	} else {
		return err
	}
}

func (s *rpn) execToBool(exec *tplExec, res *bool) *execError {
	if val, err := s.execRPN(exec); err == nil {
		if b, check := val.(bool); check {
			*res = b
			return nil
		} else {
			return s.execError("Boolean result value expected")
		}
	} else {
		return err
	}
}

func (s rpn) execType() execType {
	return execRPN
}
