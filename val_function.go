package metla

import (
	"fmt"
	"io"

	"github.com/fcg-xvii/lineman"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkFunction,
		constructor: newValFunction,
	})
}

func checkFunction(src []byte) bool {
	if lineman.CheckFirsNameChar(src) == 0 {
		return false
	}
	code := lineman.NewCodeLine(src)
	if _, check := code.ReadName(); check {
		return code.Char() == '('
	} else {
		return false
	}
	return true
}

func newValFunction(p *parser) (res token, err error) {
	name, _ := p.ReadName()
	var args []token
	p.IncPos()
loop:
	for {
		if res, err = initVal(p); err != nil {
			return
		} else {
			args = append(args, res)
			if p.IsEndDocument() {
				err = fmt.Errorf("Function parse error :: unexpected end of document")
				return
			} else {
				switch ch := p.Char(); ch {
				case ',':
					p.IncPos()
				case ')':
					p.IncPos()
					break loop
				default:
					err = fmt.Errorf("Function parse error :: unexpected symbol '%c', expected ',' or ')'", ch)
					return
				}
			}
		}
	}
	res = &valFunction{
		name: name,
		args: args,
	}
	return
}

type valFunction struct {
	name []byte
	args []token
}

func (s *valFunction) Val() interface{} {
	return s.name
}

func (s *valFunction) Data(w io.Writer, sto *storage) (err error) {
	return
}

func (s *valFunction) IsExecutable() bool { return true }

func (s *valFunction) String() string {
	res := fmt.Sprintf("[function :: { %v }, args : { %v }]", string(s.name), s.args)
	return res
}
