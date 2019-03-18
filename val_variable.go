package metla

import (
	"fmt"
	_ "reflect"

	_ "github.com/fcg-xvii/lineman"
)

type valVariable struct {
	*rawInfoRecord
	name string
}

func (s *valVariable) StorageVal(exec *tplExec) error {
	fmt.Println("STORAGE_VAL")
	if val, check := exec.sto.findVariable(s.name); check {
		exec.st.Push(val)
		return nil
	} else if _, err := exec.w.Write([]byte(s.positionWarning(fmt.Sprintf("variable not found in heap { %v }", s.name)).Error())); err != nil {
		fmt.Println("UNCKECK")
		return err
	}
	return nil
}

/*func (s *valVariable) Val(sto *storage) (res interface{}, err error) {
	if s.v == nil {
		 = sto.findVariable(s.name)
	}
}*/

/*func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValVariable,
		constructor: newValVariable,
	})
}

func checkValVariable(src []byte) bool {
	if lineman.CheckFirsNameChar(src) > 0 {
		for i := 1; i < len(src); i++ {
			if lineman.CheckBodyNameChar(src[i:]) == 0 {
				res := src[i] != '(' && src[i] != '['
				return res
			}
		}
		return true
	}
	return false
}

// Конструктор строки.
func newValVariable(p *parser) (res token, err error) {
	if name, check := p.ReadName(); !check {
		err = p.positionError("Variable parse error :: Unexpected name")
	} else {
		res = &valVariable{p.infoRecordFromMark(), string(name)}
	}
	return
}

type valVariable struct {
	*rawInfoRecord
	name string
}

func (s *valVariable) Val() interface{} {
	return s.name
}

func (s *valVariable) posInfo() *rawInfoRecord { return s.rawInfoRecord }
func (s *valVariable) String() string          { return "[variable :: { " + s.name + " }]" }
func (s *valVariable) IsExecutable() bool      { return false }
func (s *valVariable) IsNil() bool             { return true }
func (s *valVariable) IsNumber() bool          { return false }
func (s *valVariable) IsStatic() bool          { return false }
func (s *valVariable) StaticVal() interface{}  { return nil }
func (s *valVariable) Type() reflect.Kind      { return reflect.Invalid }

func (s *valVariable) execObject(sto *storage) (executor, error) {
	if v, check := sto.findVariable(s.name); !check {
		return nil, s.positionWarning(fmt.Sprintf("Variable not found [%v]", s.name))
	} else {
		return &valVariableExec{s.rawInfoRecord, v}, nil
	}
}

///////////////////////////////////////////////////////////////////

type valVariableExec struct {
	*rawInfoRecord
	v *variable
}

func (s *valVariableExec) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(fmt.Sprint(s.v.value)))
	return
}

func (s *valVariableExec) String() string {
	return "[variable { name: " + s.v.key + ", value: " + fmt.Sprint(s.v.value) + " }]"
}

func (s *valVariableExec) IsNil() bool        { return s.v.IsNil() }
func (s *valVariableExec) Type() reflect.Kind { return s.v.Kind() }
func (s *valVariableExec) ValSingle() bool    { return true }

func (s *valVariableExec) Val() (interface{}, error) {
	return s.v.value, nil
}

func (s *valVariableExec) Vals() ([]interface{}, error) {
	return []interface{}{s.v.value}, nil
}

func (s *valVariableExec) receiveEvent(name string, params []interface{}) bool {
	return false
}*/
