package metla

import (
	"fmt"
	"io"
	_ "log"
	"sync"
)

type UpdateState byte

const (
	UpdateNotNeeded UpdateState = iota
	UpdateNeeded
	ResourceNotFound
)

func (s UpdateState) String() string {
	switch s {
	case UpdateNotNeeded:
		return "UpdateNotNeeded"
	case UpdateNeeded:
		return "UpdateNeeded"
	case ResourceNotFound:
		return "ResourceNotFound"
	default:
		return "UpdateStateUndefined"
	}
}

type CheckMethod func(string, interface{}) UpdateState
type ContentMethod func(string, interface{}) ([]byte, interface{}, UpdateState)

func New(check CheckMethod, content ContentMethod) *Metla {
	return &Metla{
		check:   check,
		content: content,
		locker:  new(sync.RWMutex),
		tpls:    make(map[string]*template),
	}
}

type Metla struct {
	check   CheckMethod
	content ContentMethod
	locker  *sync.RWMutex
	tpls    map[string]*template
}

func (s *Metla) Content(path string, w io.Writer, vals map[string]interface{}) error {
	if tpl, err := s.template(path); err == nil {
		return tpl.execute(w, vals)
	} else {
		return err
	}
}

func (s *Metla) getTemplate(path string) (res *template, check bool) {
	s.locker.RLock()
	res, check = s.tpls[path]
	s.locker.RUnlock()
	return
}

func (s *Metla) template(path string) (res *template, err error) {
	var check bool
	if res, check = s.getTemplate(path); !check {
		if state := s.check(path, nil); state != ResourceNotFound {
			s.locker.Lock()
			if res, check = s.tpls[path]; !check {
				res = newTemplate(s, path)
				s.tpls[path] = res
			}
			s.locker.Unlock()
		} else {
			err = fmt.Errorf("Document not found :: [%v]", path)
		}
	}
	return
}

/*func (s *Metla) templateResult(path string, sto *storage) (*templateResult, error) {
	if tpl, err := s.template(path); err == nil {
		return tpl.result(sto)
	} else {
		return nil, err
	}
}*/

func (s *Metla) removeTempalte(path string) {
	s.locker.Lock()
	delete(s.tpls, path)
	s.locker.Unlock()
}

//func (s *Metla) tplArrived(t)
