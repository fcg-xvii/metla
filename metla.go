package metla

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type UpdateState byte

const (
	UpdateNotNeeded UpdateState = iota
	UpdateNeeded
	ResourceNotFound
)

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
	log.Println("METLA :: CONTENT")
	if tpl, err := s.template(path); err == nil {
		//sto := newStorage(vals)
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

func (s *Metla) removeTempalte(path string) {
	s.locker.Lock()
	delete(s.tpls, path)
	s.locker.Unlock()
}

//func (s *Metla) tplArrived(t)
