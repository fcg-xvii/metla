package metla

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

func init() {
	keywords["include"] = newKeyInclude
}

func newKeyInclude(p *parser) (t token, err error) {
	p.PassSpaces()
	var pathToken, paramsToken token
	if pathToken, err = initVal(p); err == nil {
		p.PassSpaces()
		if !p.IsEndLine() {
			if !checkValObject(p.EndLineContent()) {
				err = errors.New("Include parse error :: Unexpected params token. [object] token expected")
			} else {
				if paramsToken, err = newValObject(p); err == nil {
					p.PassSpaces()
					if !p.IsEndDocument() || !p.IsEndLine() {
						err = fmt.Errorf("Include parse error :: Unespected symbol [%c], end line expected", p.Char())
					} else {
						t = &keyInclude{rawInfoRecord: p.infoRecordFromMark(), tplPath: pathToken, params: paramsToken}
						//fmt.Println("INCLUDE", t)
					}
				}
			}
		}
	}
	return
}

type keyInclude struct {
	*rawInfoRecord
	tplPath token
	params  token
}

func (s *keyInclude) execObject(sto *storage, tpl *template) (res execObject, err error) {
	var tplPath execObject
	if tplPath, err = s.tplPath.execObject(sto, tpl); err == nil {
		if tplPath.Type() != reflect.String {
			err = fmt.Errorf("Include token exec error :: Unexpected include path token [%s], expected [string] type", tplPath.Type())
			return
		}
	}
	var val interface{}
	if val, err = tplPath.Val(); err == nil {
		sto.newLayout()
		if s.params != nil {
			var params execObject
			if params, err = s.params.execObject(sto, tpl); err != nil {
				return
			}
			for key, val := range params.(*valObjectExec).Map() {
				sto.appendValue(key, val)
			}
		}
		var tplRes *templateResult
		if tplRes, err = tpl.root.templateResult(val.(string), sto); err != nil {
			return
		}
		sto.dropLayout()
		res = &execObjectInclude{s.rawInfoRecord, tplRes}
	}
	return
}

func (s *keyInclude) String() string {
	return "[include :: { " + s.tplPath.String() + " }, { " + s.params.String() + " }]"
}

func (s *keyInclude) IsExecutable() bool { return true }

////////////////////////////////////////////////////////////////////

type execObjectInclude struct {
	*rawInfoRecord
	tpl *templateResult
}

func (s *execObjectInclude) Data(w io.Writer) (err error) {
	s.tpl.exec(w)
	return
}

func (s *execObjectInclude) IsNil() bool { return false }

func (s *execObjectInclude) String() string { return "include..." }

// reflect.Kind тип не очень подходит в этом случае?
func (s *execObjectInclude) Type() reflect.Kind { return reflect.Invalid }

func (s *execObjectInclude) Val() (interface{}, error) {
	return nil, errors.New("Include method is not like value")
}

func (s *execObjectInclude) Vals() ([]interface{}, error) {
	return nil, errors.New("Include method is not like values")
}

func (s *execObjectInclude) ValSingle() bool { return true }
