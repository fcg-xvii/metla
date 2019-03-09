package metla

import (
	"fmt"
	"io"
	"sync"
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
	tokenList  []token
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
	fmt.Println("Update proceed...", s.err)
	if s.err != nil {
		return s.err
	}
	sto := newStorage(vals)
	if result, err := s.result(sto); err != nil {
		return err
	} else {
		result.exec(w)
		return nil
	}
}

func (s *template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *template) result(sto *storage) (*templateResult, error) {
	s.locker.RLock()
	//fmt.Println("Result...", s.tokenList)
	if s.err != nil {
		s.locker.RUnlock()
		return nil, s.err
	}
	res := &templateResult{make([]execObject, 0, len(s.tokenList))}
	for _, v := range s.tokenList {
		if eObj, err := v.execObject(sto, s); err == nil {
			res.list = append(res.list, eObj)
		} else {
			s.locker.RUnlock()
			return nil, err
		}
	}
	s.locker.RUnlock()
	return res, nil
}

////////////////////////////////////////////////////////////////////////

type templateResult struct {
	list []execObject
}

func (s *templateResult) exec(w io.Writer) (err error) {
	for _, v := range s.list {
		if err = v.Data(w); err != nil {
			return
		}
	}
	return
}
