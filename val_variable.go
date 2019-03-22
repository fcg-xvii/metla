package metla

import (
	_ "fmt"
)

type valVariable struct {
	*rawInfoRecord
	name string
}

func (s *valVariable) StorageVal(exec *tplExec) *variable {
	val, check := exec.sto.findVariable(s.name)
	if !check {
		val = &variable{key: s.name, value: nil}
	}
	exec.st.Push(val)
	return val
}

func (s *valVariable) String() string {
	return "{ variable :: " + s.name + " }"
}
