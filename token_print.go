package metla

import (
	"io"
)

type tokenPrint struct {
	val token
}

func (s *tokenPrint) execObject(sto *storage, tpl *template) (obj execObject, err error) {
	if obj, err = s.val.execObject(sto, tpl); err == nil {
		obj = &execObjectPrint{obj}
	}
	return
}

func (s *tokenPrint) String() string {
	return "[tokenPrint :: { " + s.String() + " }]"
}

func (s *tokenPrint) IsExecutable() bool { return false }

//////////////////////////////////////////////////////////////////////////////

type execObjectPrint struct {
	val execObject
}

func (s *execObjectPrint) Data(w io.Writer) (err error) {
	err = s.val.Data(w)
	return
}

func (s *execObjectPrint) IsNil() bool {
	return false
}

func (s *execObjectPrint) Type() execObjectType {
	return execObjectToken
}

func (s *execObjectPrint) Val() interface{} {
	return s.val.Val()
}

func (s *execObjectPrint) Vals() []interface{} {
	return []interface{}{s.val.Val()}
}

func (s *execObjectPrint) ValSingle() bool {
	return true
}

func (s *execObjectPrint) String() string {
	return "[print { " + s.val.String() + " }"
}
