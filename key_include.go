package metla

import (
	"errors"
	"fmt"
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
						t = &keyInclude{tplPath: pathToken, paramsToken: paramsToken}
						//fmt.Println("INCLUDE", t)
					}
				}
			}
		}
	}
	return
}

type keyInclude struct {
	tplPath     token
	paramsToken token
}

func (s *keyInclude) Val() interface{} {
	return s.tplPath
}

func (s *keyInclude) Data() (res []byte, err error) {
	return nil, nil
}

func (s *keyInclude) String() string {
	return "[include :: {" + s.tplPath.String() + "}, { " + s.paramsToken.String() + " }]"
}
func (s *keyInclude) IsExecutable() bool { return true }
