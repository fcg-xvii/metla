package metla

import (
	"fmt"
	"io"
	"reflect"

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
		name: string(name),
		args: args,
	}
	return
}

type valFunction struct {
	name string
	args []token
}

func (s *valFunction) Val() interface{} {
	return s.name
}

func (s *valFunction) IsExecutable() bool { return true }

func (s *valFunction) String() string {
	res := fmt.Sprintf("[function :: { %v }, args : { %v }]", string(s.name), s.args)
	return res
}

func (s *valFunction) execObject(sto *storage, tpl *template) (res execObject, err error) {
	if f, check := sto.findVariable(s.name); check {
		if f.Kind() != reflect.Func {
			err = fmt.Errorf("Function exec error :: unexpected variable type [%v], [Func] expected", f.Kind())
		} else {
			res = &valFunctionExec{f, make([]execObject, len(s.args))}
			for i, v := range s.args {
				if res.args[i], err = v.execObject(sto, tpl); err != nil {
					return
				}
			}
		}
	} else {
		err = fmt.Errorf("Function exec error :: [%s] not found", s.name)
	}
	return
}

//////////////////////////////////////////////////////////

type valFunctionExec struct {
	f    *variable
	args []execObject
}

func (s *valFunctionExec) Data(w io.Writer) (err error) {
	fVal := reflect.ValueOf(s.f.value)
	if fVal.Kind() != reflect.Func {
		err = fmt.Errorf("Function exec error :: variable [%v] is not a function", s.f.key)
	} else {
		rArgs := make([]reflect.Value, 0, len(s.args))
		for i, v := range s.args {
			if v.ValSingle() {
				rArgs = append(rArgs, reflect.ValueOf(v.Val()))
			} else {
				for _, v := range v.Vals() {

				}
			}
		}
	}
}
