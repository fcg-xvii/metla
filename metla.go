package metla

import (
	"fmt"
	"io"
	_ "log"
	"sync"
	"time"
)

type UpdateState byte

const (
	ResourceNotFound UpdateState = iota
	UpdateNotNeeded
	UpdateNeeded
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

type CheckMethod func(string, *time.Time) UpdateState
type ContentMethod func(string, *time.Time) ([]byte, time.Time, UpdateState)

func New(check CheckMethod, content ContentMethod) *Metla {
	return &Metla{
		check:   check,
		content: content,
		locker:  new(sync.RWMutex),
		tpls:    make(map[string]*Template),
	}
}

type Metla struct {
	check   CheckMethod
	content ContentMethod
	locker  *sync.RWMutex
	tpls    map[string]*Template
}

func (s *Metla) Content(path string, w io.Writer, vals map[string]interface{}) (modified time.Time, err error) {
	if tpl, err := s.Template(path); err == nil {
		return tpl.Execute(w, vals)
	}
	return
}

func (s *Metla) getTemplate(path string) (res *Template, check bool) {
	s.locker.RLock()
	res, check = s.tpls[path]
	s.locker.RUnlock()
	return
}

func (s *Metla) Template(path string) (res *Template, err error) {
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

func (s *Metla) removeTempalte(path string) {
	s.locker.Lock()
	delete(s.tpls, path)
	s.locker.Unlock()
}

//func (s *Metla) tplArrived(t)
