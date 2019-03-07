package metla

import (
	"fmt"
	"io"
	"reflect"
)

func init() {
	keywords["for"] = newKeyFor
}

func newKeyFor(p *parser) (t token, err error) {
	fmt.Println("FOR_ERROR", string(p.EndLineContent()))
	var index token
	p.PassSpaces()
	if index, err = newValVariable(p); err != nil {
		return
	}
	p.PassSpaces()
	switch {
	case p.PosMatchSlice([]byte("in")):
		t, err = newCycleIn(index, p)
	default:
		err = fmt.Errorf("Keyword [for] parse error :: unexpected char after index [%c]", p.Char())
		return
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

type cycle struct {
	index  token
	childs []token
}

func (s *cycle) IsExecutable() bool {
	return true
}

func (s *cycle) String() string {
	res := "[cycle... {\r\n"
	for _, v := range s.childs {
		res += "    " + v.String() + "\r\n"
	}
	res += "}"
	return res
}

func (s *cycle) execObject(*storage, *template) (execObject, error) {
	return nil, fmt.Errorf("Cycle is abstract object")
}

func (s *cycle) appendChild(t token) {
	if t != nil {
		s.childs = append(s.childs, t)
	}
}

type cycleExec struct {
	index  *variable
	childs []execObject
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

func newCycleIn(index token, p *parser) (t token, err error) {

	fmt.Println("ELC", string(p.EndLineContent()))
	p.ForwardPos(2)
	var (
		left token
		name string
	)
	if left, err = initVal(p); err != nil {
		return
	} else {
		v, check := left.(*varVariable); !check {
			err = fmt.Errorf("Keyword [for] parse error :: Index var must be variable")
		}		
	}

		p.PassSpaces()
	if p.Char() != ':' {
		err = fmt.Errorf("Keyword [for] parse error :: unexpected char after left side '%c', ':' expected", p.Char())
	}
	p.IncPos()
	var right token
	if right, err = initVal(p); err != nil {
		return
	}
	fmt.Println("LRRRRR", left, right)
	in := cycleIn{&cycle{index, nil}, left, right}
	if err = p.parseCodeToCloseTag("endfor", &in); err == nil {
		t = &in
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////

type cycleIn struct {
	indexVarName string
	*cycle
	left, right token
}

func (s *cycleIn) execObject(sto *storage, tpl *template) (res execObject, err error) {
	sto.newLayout()
	defer sto.dropLayout()
	var indexObj execObject
	if indexObj, err = s.index.execObject(sto, tpl); err == nil {
		if index, check := indexObj.(*valVariableExec); !check {
			err = fmt.Errorf("Cycle prepare error :: index var must be variable type")
		} else {
			var leftObj, rightObj execObject
			if leftObj, err = s.left.execObject(sto, tpl); err != nil {
				return
			}
			if rightObj, err = s.right.execObject(sto, tpl); err != nil {
				return
			}
			in := &cycleInExec{
				cycleExec: &cycleExec{index: index.v, childs: make([]execObject, len(s.childs))},
				left:      leftObj,
				right:     rightObj,
			}
			for i, v := range s.childs {
				if leftObj, err = v.execObject(sto, tpl); err == nil {
					in.childs[i] = leftObj
				} else {
					return
				}
			}
			res = in
		}
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////

type cycleInExec struct {
	*cycleExec
	left, right execObject
}

func (s *cycleInExec) Data(w io.Writer) (err error) {
	var (
		lVal, rVal interface{}
		lNum, rNum int64
		check      bool
	)
	if lVal, err = s.left.Val(); err != nil {
		_, err = w.Write([]byte(err.Error()))
		return
	}
	if rVal, err = s.right.Val(); err != nil {
		_, err = w.Write([]byte(err.Error()))
		return
	}
	if lNum, check = checkIfaceInt(lVal); !check {
		w.Write([]byte(fmt.Sprintf("Cycle exec error :: left value must be integer, not [%v]", reflect.ValueOf(lVal).Kind())))
	}
	if rNum, check = checkIfaceInt(lVal); !check {
		w.Write([]byte(fmt.Sprintf("Cycle exec error :: right value must be integer, not [%v]", reflect.ValueOf(rVal).Kind())))
	}
	for {
		if lNum > rNum {
			return nil
		}
		s.index.value = lNum
		for _, v := range s.childs {
			if err := v.Data(w); err != nil {
				return err
			}
		}
		lNum++
	}
}

func (s *cycleInExec) IsNil() bool {
	return false
}

func (s *cycleInExec) String() string {
	return "[cycle in...]"
}

func (s *cycleInExec) Type() reflect.Kind           { return reflect.Invalid }
func (s *cycleInExec) Val() (interface{}, error)    { return nil, nil }
func (s *cycleInExec) Vals() ([]interface{}, error) { return nil, nil }
func (s *cycleInExec) ValSingle() bool              { return true }
