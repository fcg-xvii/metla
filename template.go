package metla

import (
	"fmt"
	"io"
	_ "reflect"
	"sync"

	"github.com/golang-collections/collections/stack"
)

var (
	MaxStorageLayouts = 150
)

func newTemplate(root *Metla, objPath string) *Template {
	return &Template{
		root:    root,
		objPath: objPath,
		locker:  new(sync.RWMutex),
	}
}

type Template struct {
	root       *Metla
	objPath    string
	locker     *sync.RWMutex
	tokenList  []interface{}
	updateMark interface{}
	err        error
}

func (s *Template) checkUpdate() error {
	switch s.root.check(s.objPath, s.updateMark) {
	case ResourceNotFound:
		{
			s.root.removeTempalte(s.objPath)
			s.err = fmt.Errorf("Document not found :: [%v]", s.objPath)
			return s.err
		}
	case UpdateNeeded:
		{
			s.locker.Lock()
			if content, newMark, state := s.root.content(s.objPath, s.updateMark); state == UpdateNeeded {
				s.updateMark = newMark
				s.err = s.parse(content)
			}
			s.locker.Unlock()
		}
	}
	return nil
}

func (s *Template) execute(w io.Writer, vals map[string]interface{}) error {
	s.checkUpdate()
	if s.err != nil {
		return s.err
	}
	sto := newStorage(vals)
	return s.result(sto, w)
}

func (s *Template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *Template) pushToken(t interface{}) {
	s.tokenList = append(s.tokenList, t)
}

func (s *Template) result(sto *storage, w io.Writer) (err error) {
	s.locker.RLock()
	if s.err != nil {
		s.locker.RUnlock()
		return s.err
	}
	list := make([]interface{}, len(s.tokenList))
	copy(list, s.tokenList)
	s.locker.RUnlock()
	sto.newLayout()
	if len(sto.layouts) >= MaxStorageLayouts {
		return fmt.Errorf("Fatal error :: Include loop arrived - max storage layouts")
	}
	tplExec := &tplExec{list, stack.New(), sto, 0, w, 0, false, s.root}
	err = tplExec.exec()
	sto.dropLayout()
	return
}

type tplExec struct {
	list        []interface{}
	st          *stack.Stack
	sto         *storage
	index       int
	w           io.Writer
	fieldLayout int
	breakFlag   bool
	root        *Metla
}

func (s *tplExec) exec() (err error) {
	s.index = -1
	for s.index < len(s.list) {
		if err = s.execNext(); err != nil {
			return
		}
	}
	return
}

func (s *tplExec) execNext() (err error) {
	s.index++
	if s.index >= len(s.list) {
		return
	}
	//fmt.Println("EXEC_NEXT", s.list[s.index])
	switch s.list[s.index].(type) {
	case *execCommand:
		//fmt.Println("EXEC_COMMAND", s)
		exec := s.list[s.index].(*execCommand)
		//fmt.Println(exec.name)
		if s.fieldLayout > 0 && (exec.name != "field-end" && exec.name != "field-start") {
			s.st.Push(exec)
			break
		}
		if err = exec.method(s, exec.rawInfoRecord); err != nil {
			return
		}
	case *valVariable:
		//fmt.Println("VAL_VARIABLE")
		s.list[s.index] = s.list[s.index].(*valVariable).StorageVal(s)

	case *operator:
		//fmt.Println("OPERATOR")
		if err = s.list[s.index].(*operator).exec(s.st); err != nil {
			return
		}
	default:
		//fmt.Println("DEFAULT...", s.list[s.index])
		s.st.Push(s.list[s.index])
	}
	return
}

////////////////////////////////////////////////////////////////////////
