package metla

import (
	"fmt"
	"io"
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

func (s *template) execute(w io.Writer, vals map[string]interface{}) error {
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
	fmt.Println("Update proceed...")
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
	list := s.tokenList
	s.locker.RUnlock()
	fmt.Println("LIST", list)
	tplExec := &tplExec{list, stack.New(), sto, 0, w}
	return tplExec.exec()
}

type tplExec struct {
	list  []interface{}
	st    *stack.Stack
	sto   *storage
	index int
	w     io.Writer
}

func (s *tplExec) exec() (err error) {
	fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	for s.index < len(s.list) {
		switch s.list[s.index].(type) {
		case *execCommand:
			exec := s.list[s.index].(*execCommand)
			if err = exec.method(s, exec.rawInfoRecord); err != nil {
				return
			}
		case *valVariable:
			if err = s.list[s.index].(*valVariable).StorageVal(s); err != nil {
				return
			}
		default:
			s.st.Push(s.list[s.index])
		}
		s.index++
	}
	return
}

////////////////////////////////////////////////////////////////////////
