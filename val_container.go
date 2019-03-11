package metla

import (
	"fmt"
	"io"
	"reflect"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValArray,
		constructor: newValArray,
	}, &valueCreator{
		checker:     checkValObject,
		constructor: newValObject,
	})
}

func checkValArray(src []byte) bool {
	return src[0] == '['
}

func checkValObject(src []byte) bool {
	return src[0] == '{'
}

func newValArray(p *parser) (res token, err error) {
	var vals []token
	p.IncPos()
	var t token
loop:
	for !p.IsEndDocument() {
		p.PassEndLines()
		if t, err = initVal(p); err != nil {
			return nil, err
		} else {
			vals = append(vals, t)
			p.PassEndLines()
			switch p.Char() {
			case ',':
				p.IncPos()
			case ']':
				p.IncPos()
				break loop
			}
		}
	}
	res = &valArray{p.infoRecordFromMark(), vals}
	fmt.Println(res)
	return
}

/////////////////////////////////////////////////////////////////

type valArray struct {
	*rawInfoRecord
	vals []token
}

func (s *valArray) String() string {
	return fmt.Sprintf("[array :: { %v }]", s.vals)
}

func (s *valArray) IsExecutable() bool { return false }

func (s *valArray) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
	vals := make([]execObject, len(s.vals))
	for i, v := range s.vals {
		if vals[i], err = v.execObject(sto, tpl, nil); err != nil {
			return
		}
	}
	res = &valArrayExec{s.rawInfoRecord, vals}
	return
}

//////////////////////////////////////////////////////////////////

type valArrayExec struct {
	*rawInfoRecord
	vals []execObject
}

func (s *valArrayExec) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(s.String()))
	return
}

func (s *valArrayExec) String() string {
	return fmt.Sprintf("[array { %v }]", s.vals)
}

func (s *valArrayExec) IsNil() bool        { return s.vals == nil }
func (s *valArrayExec) Type() reflect.Kind { return reflect.Slice }

func (s *valArrayExec) Val() (interface{}, error)    { return s.vals, nil }
func (s *valArrayExec) Vals() ([]interface{}, error) { return []interface{}{s.vals}, nil }
func (s *valArrayExec) ValSingle() bool              { return true }

func (s *valArrayExec) receiveEvent(name string, params []interface{}) bool { return false }

//////////////////////////////////////////////////////////////////

func newValObject(p *parser) (res token, err error) {
	m := make(map[string]token)
	var (
		check bool
		key   []byte
		val   token
	)
	p.IncPos()
loop:
	for !p.IsEndDocument() {
		// Парсим
		p.PassEndLines()
		if key, check = p.ReadName(); check {
			fmt.Println("KEY.....................", string(key))
			p.PassSpaces()
			if p.Char() != ':' {
				err = fmt.Errorf("Object parse error :: Unexpected symbol '%c', expected ':'", p.Char())
				return
			}
			p.IncPos()
			if val, err = initVal(p); err != nil {
				return
			}
			m[string(key)] = val
			p.PassEndLines()
			switch p.Char() {
			case ',':
				p.IncPos()
			case '}':
				//fmt.Println("ENDDDDDD .............................")
				p.IncPos()
				break loop
			default:
				err = fmt.Errorf("Object parse error :: Unexpected symbol '%c', ',' or '}' expected", p.Char())
			}
		} else {
			err = fmt.Errorf("Object parse error :: Unexpected variable name...")
			return
		}
	}
	res = &valObject{p.infoRecordFromMark(), m}
	fmt.Println("==================", res)
	return
}

type valObject struct {
	*rawInfoRecord
	vals map[string]token
}

func (s *valObject) Val() interface{} {
	return s.vals
}

func (s *valObject) String() string {
	return fmt.Sprintf("[object :: { %v }]", s.vals)
}

func (s *valObject) IsExecutable() bool { return false }

func (s *valObject) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
	vals := make(map[string]execObject)
	for key, val := range s.vals {
		if vals[key], err = val.execObject(sto, tpl, nil); err == nil {
			return
		}
	}
	res = &valObjectExec{s.rawInfoRecord, vals}
	return
}

///////////////////////////////////////////////////////////////////

type valObjectExec struct {
	*rawInfoRecord
	vals map[string]execObject
}

func (s *valObjectExec) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(s.String()))
	return
}

func (s *valObjectExec) IsNil() bool { return s.vals == nil }

func (s *valObjectExec) String() string {
	return fmt.Sprintf("[object { %v }]", s.vals)
}

func (s *valObjectExec) Type() reflect.Kind           { return reflect.Map }
func (s *valObjectExec) Val() (interface{}, error)    { return s.vals, nil }
func (s *valObjectExec) Vals() ([]interface{}, error) { return []interface{}{s.vals}, nil }
func (s *valObjectExec) ValSingle() bool              { return true }
func (s *valObjectExec) Map() map[string]execObject   { return s.vals }

func (s *valObjectExec) receiveEvent(name string, params []interface{}) bool { return false }
