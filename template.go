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

func (s *template) execute(w io.Writer, vals map[string]interface{}) (err error) {
	switch s.root.check(s.objPath, s.updateMark) {
	case ResourceNotFound:
		{
			s.root.removeTempalte(s.objPath)
			err = fmt.Errorf("Document not found :: [%v]", path)
		}
	case UpdateNeeded:
		{
			s.locker.Lock()
			if content, newMark, state := s.root.content(s.objPath, s.updateMark); state == UpdateNeeded {
				s.updateMark = newMark
				s.parse(content)
			}
			s.locker.Unlock()
		}
	}
	s.locker.RLock()
	if err != nil {
		s.locker.RUnlock()
		return
	}
	err = s.exec(w, vals)
	s.locker.RUnlock()
	return
}

/*func (s *template) execute(w io.Writer, vals map[string]interface{}) (err error) {
	fmt.Println("TEMPLATE_EXECUTE...")
	var (
		check  bool
		source []byte
	)
	s.locker.Lock()
	if source, check, s.updateMark, err = s.root.content(s.objPath, s.updateMark); err == nil {
		if check {
			fmt.Println("TEMPLATE_UPDATE")
			err = s.parse(source)
		}
		if err == nil {
			//fmt.Println("TEMPLATE_EXECUTE")
			sto := newStorage(vals)
			err = s.exec(w, sto)
		}
	}
	s.locker.Unlock()
	fmt.Println(s.tokenList)
	return
}

func (s *template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *tempalte) result() *templateResult {
	res := &templateResult{make([]execObject, len(s.))}
}

func (s *template) exec(w io.Writer, sto *storage) (err error) {
	var execObj execObject
	for _, v := range s.tokenList {
		if execObj, err = v.execObject(sto, s); err != nil {
			fmt.Println("Exec error ::", err)
			return
		}
		if err = execObj.Data(w); err != nil {
			fmt.Println("ExecObj error ::", err)
			return
		}
		if err = v.Data(w, sto); err != nil {
			return
		}
	}
	return

	tplResult =
}*/

////////////////////////////////////////////////////////////////////////

type templateResult struct {
	list []execObject
}
