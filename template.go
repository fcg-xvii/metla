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
	//fmt.Println("Result...", s.tokenList)
	if s.err != nil {
		s.locker.RUnlock()
		return s.err
	}
	/*res := &templateResult{make([]interface{}, 0, len(s.tokenList))}
	for _, v := range s.tokenList {
		if eObj, err := v.execObject(sto); err == nil {
			res.list = append(res.list, eObj)
		} else {
			s.locker.RUnlock()
			return nil, err
		}
	}*/
	list, st := s.tokenList, stack.New()
	s.locker.RUnlock()
	fmt.Println("LIST", list)
	for len(list) > 0 {
		//fmt.Println(len(list))
		if obj, check := list[0].(*execCommand); check {
			if list, err = obj.method(list, st, sto, w); err != nil {
				return err
			}
		} else {
			st.Push(list[0])
			list = list[1:]
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////
