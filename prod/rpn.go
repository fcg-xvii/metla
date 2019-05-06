package prod

import (
	"fmt"

	"github.com/fcg-xvii/containers"
)

func newRPN(p *parser) *parseError {
	fmt.Println("NEW_RPN", string(p.Char()))
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
		switch {
		case isOperatorSymbol(p):
			fmt.Println("OPSYMBOL", st)
			if op, err := initOperator(p); err != nil {
				return err
			} else {
				for st.Len() > 0 {
					if sOp, check := st.Peek().(operator); check {
						if sOp.prefix || op.priority <= sOp.priority {
							obj.pn = append(obj.pn, st.Pop())
						} else {
							break
						}
					} else {
						break
					}
				}
				if valAccepted {
					fmt.Println("ValAccepted")
					valAccepted = false
				}
				st.Push(op)
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
				fmt.Println("%%%%%%%%%%%%%%%%%%", st)
			popLoop:
				for st.Len() > 0 {
					fmt.Printf("fff %T\n", st.Peek())
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
			fmt.Println("888", string(p.Char()))
			if !isArifmeticSymbol(p) {
				break mainLoop
			}
			fmt.Println("PASS")
		}
	}
	fmt.Println("BRACKET_LAYOUT", bracketLayout)
	if bracketLayout != 0 {
		return obj.parseError("Unclosed bracket")
	}
	obj.pn = append(obj.pn, st.PopAllReverse()...)
	p.stack.Push(obj)
	return nil
}

type rpn struct {
	position
	pn []interface{}
}

func (s rpn) exec(exec *tplExec) *execError {
	st := new(containers.Stack)
	for _, v := range s.pn {
		if op, check := v.(operator); check {
			o
		}
	}
}
