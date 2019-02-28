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
			args := make([]execObject, len(s.args))
			for i, v := range s.args {
				if args[i], err = v.execObject(sto, tpl); err != nil {
					return
				}
			}
			res = &valFunctionExec{f, args}
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
	_, err = s.call()
	return
}

func (s *valFunctionExec) IsNil() bool {
	return s.f.IsNil()
}

func (s *valFunctionExec) String() string {
	return "[function { " + s.f.key + " }]"
}

func (s *valFunctionExec) Type() reflect.Kind {
	return reflect.Func
}

func (s *valFunctionExec) Val() (interface{}, error) {
	if vals, err := s.Vals(); err == nil && len(vals) > 0 {
		return vals[0], err
	} else {
		return nil, err
	}
	/*if rRes, err := s.call(); err == nil {
		if len(rRes) > 0 {
			return rRes[0].Interface(), nil
		}
	} else {
		return nil, err
	}*/
}

func (s *valFunctionExec) Vals() (res []interface{}, err error) {
	var rRes []reflect.Value
	if rRes, err = s.call(); err == nil {
		if len(rRes) > 0 {
			res = make([]interface{}, len(rRes))
			for i, v := range rRes {
				res[i] = v.Interface()
			}
		}
	}
	return
}

func (s *valFunctionExec) ValSingle() bool {
	return false
}

func (s *valFunctionExec) call() (res []reflect.Value, err error) {
	fVal := reflect.ValueOf(s.f.value)
	if fVal.Kind() != reflect.Func {
		err = fmt.Errorf("Function exec error :: variable [%v] is not a function", s.f.key)
	} else {
		rArgs := make([]reflect.Value, 0, len(s.args))
		for _, v := range s.args {
			if v.ValSingle() {
				var val interface{}
				if val, err = v.Val(); err == nil {
					rArgs = append(rArgs, reflect.ValueOf(val))
				}
			} else {
				var vals []interface{}
				if vals, err = v.Vals(); err == nil {
					for _, v := range vals {
						rArgs = append(rArgs, reflect.ValueOf(v))
					}
				}
			}
		}
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Function call fatal error (panic) :: %v", r)
			}
		}()
		res = fVal.Call(rArgs)
	}
	return
}
