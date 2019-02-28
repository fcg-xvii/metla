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
						t = &keyInclude{tplPath: pathToken, params: paramsToken}
						//fmt.Println("INCLUDE", t)
					}
				}
			}
		}
	}
	return
}

type keyInclude struct {
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
	r := execObjectInclude{sto: sto}
	var val interface{}
	if val, err = tplPath.Val(); err == nil {
		if r.tpl, err = tpl.root.template(val.(string)); err != nil {
			return
		}
		r.params, err = s.params.execObject(sto, tpl)
	}
	return
}

func (s *keyInclude) String() string {
	return "[include :: {" + s.tplPath.String() + "}, { " + s.params.String() + " }]"
}

func (s *keyInclude) IsExecutable() bool { return true }

////////////////////////////////////////////////////////////////////

type execObjectInclude struct {
	tpl    *template
	params execObject
	sto    *storage
}

func (s *execObjectInclude) Data(w io.Writer) (err error) {
	// Создать слой параметров в sto
	err = s.tpl.exec(w, s.sto)
	// Прибить слой параметрв в sto
	return
}
