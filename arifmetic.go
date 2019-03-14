package metla

import (
	"fmt"
	"io"
	"reflect"
)

/*var (
	lNumError = errors.New("left side is not number")
	rNumError = errors.New("right side is not number")
)*/

func initArifmetic(p *parser) (res token, err error) {
	fmt.Println("INIT_ART")
	var pn []interface{}
	p.SetupMark()
	info := p.infoRecordFromMark()
	if pn, err = parseRPN(p); err == nil {
		if pn, err = simpleRPN(pn); err == nil {
			res = &arifmetic{info, pn}
		}
	}
	//fmt.Println(execRPN(pn))
	return
}

type arifmetic struct {
	*rawInfoRecord
	pn []interface{}
}

func (s *arifmetic) IsExecutable() bool { return false }
func (s *arifmetic) String() string     { return "[arifmetic...]" }

func (s *arifmetic) execObject(sto *storage, tpl *template, parent executor) (res executor, err error) {
	pn := make([]interface{}, len(s.pn))
	for i, v := range s.pn {
		if t, check := v.(token); check {
			if pn[i], err = t.execObject(sto, tpl, parent); err != nil {
				return
			}
		} else {
			pn[i] = v
		}
	}
	res = &arifmeticExec{s.rawInfoRecord, pn}
	return
}

type arifmeticExec struct {
	*rawInfoRecord
	pn []interface{}
}

func (s *arifmeticExec) result() (res interface{}, err error) {
	if len(s.pn) == 1 {
		res = s.pn[0]
	} else {
		res, err = execRPN(s.pn)
	}
	return
}

func (s *arifmeticExec) ValSingle() bool                                     { return true }
func (s *arifmeticExec) IsNil() bool                                         { return false }
func (s *arifmeticExec) Type() reflect.Kind                                  { return reflect.Invalid }
func (s *arifmeticExec) Val() (interface{}, error)                           { return s.result() }
func (s *arifmeticExec) receiveEvent(name string, params []interface{}) bool { return false }
func (s *arifmeticExec) Vals() (res []interface{}, err error) {
	if val, err := s.result(); err == nil {
		res = []interface{}{val}
	}
	return
}

func (s *arifmeticExec) String() string {
	if res, err := s.result(); err == nil {
		return fmt.Sprint(res)
	} else {
		return s.positionWarning(err.Error()).Error()
	}
}

func (s *arifmeticExec) Data(w io.Writer) (err error) {
	var res interface{}
	if res, err = s.Val(); err == nil {
		_, err = w.Write([]byte(fmt.Sprint(res)))
	}
	return
}
