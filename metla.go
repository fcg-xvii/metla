package metla

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type CheckMethod func(string) bool
type ContentMethod func(string, interface{}) ([]byte, bool, interface{}, error)

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

func (s *Metla) getTemplate(path string) (res *template, check bool) {
	s.locker.RLock()
	res, check = s.tpls[path]
	s.locker.RUnlock()
	return
}

func (s *Metla) Content(path string, w io.Writer, vals map[string]interface{}) (err error) {
	log.Println("METLA :: CONTENT")
	if tpl, check := s.getTemplate(path); check {
		tpl.execute(w, vals)
	} else {
		s.locker.Lock()
		if tpl, check := s.tpls[path]; !check {
			if s.check(path) {
				tpl = newTemplate(s, path)
				s.tpls[path] = tpl
				s.locker.Unlock()
				err = tpl.execute(w, vals)
				return
			} else {
				err = fmt.Errorf("Document not found [%v]", path)
			}
		} else {
			s.locker.Unlock()
			err = tpl.execute(w, vals)
			return
		}
		s.locker.Unlock()
	}
	return
}

//func (s *Metla) tplArrived(t)
