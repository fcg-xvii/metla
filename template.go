package metla

import (
	"fmt"
	"io"
	_ "reflect"
	"sync"

	"github.com/golang-collections/collections/stack"
)

func newTemplate(root *Metla, objPath string) *template {
	return &template{
		root:    root,
		objPath: objPath,
		locker:  new(sync.RWMutex),
	}
}

type template struct {
	root       *Metla
	objPath    string
	locker     *sync.RWMutex
	tokenList  []interface{}
	updateMark interface{}
	err        error
	//lastRequest time.Time
}

func (s *template) checkUpdate() error {
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

func (s *template) execute(w io.Writer, vals map[string]interface{}) error {
	s.checkUpdate()
	if s.err != nil {
		return s.err
	}
	sto := newStorage(vals)
	return s.result(sto, w)
}

func (s *template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *template) pushToken(t interface{}) {
	s.tokenList = append(s.tokenList, t)
}

func (s *template) result(sto *storage, w io.Writer) (err error) {
	s.locker.RLock()
	if s.err != nil {
		s.locker.RUnlock()
		return s.err
	}
	list := make([]interface{}, len(s.tokenList))
	copy(list, s.tokenList)
	s.locker.RUnlock()
	//fmt.Println("LIST", list)
	sto.newLayout()
	tplExec := &tplExec{list, stack.New(), sto, 0, w, 0, s.root}
	sto.dropLayout()
	return tplExec.exec()
}

type tplExec struct {
	list        []interface{}
	st          *stack.Stack
	sto         *storage
	index       int
	w           io.Writer
	fieldLayout int
	root        *Metla
}

func (s *tplExec) exec() (err error) {
	//fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
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
