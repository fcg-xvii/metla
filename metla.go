package metla

import (
	"io"
	"sync/atomic"
	"time"

	"github.com/fcg-xvii/containers"
)

var (
	DEFAULT_MAX_EXEC_DURATION = time.Second * 30
)

type Requester interface {
	RequestContent(path string) (content []byte, marker time.Time, exists bool, err error)
	RequestUpdate(path string, modified time.Time) (content []byte, marker time.Time, exists bool, err error)
}

func New(requester Requester) *Metla {
	return &Metla{int64(DEFAULT_MAX_EXEC_DURATION), requester, containers.NewCache(time.Hour, time.Hour, nil)}
}

type Metla struct {
	maxExecDuration int64
	requester       Requester
	store           *containers.Cache
}

func (s *Metla) SetMaxExecDuration(execTime time.Duration) {
	atomic.StoreInt64(&s.maxExecDuration, int64(execTime))
}

func (s *Metla) template(path string) (tpl *template, check bool) {
	var iface interface{}
	if iface, check = s.store.GetOrCreate(path, func(key interface{}) (interface{}, bool) {
		if content, modified, exists, err := s.requester.RequestContent(path); exists && err == nil {
			return newTemplate(s.requester, s, path, content, modified), true
		}
		return nil, false
	}); check {
		tpl = iface.(*template)
	}
	return
}

func (s *Metla) Content(path string, w io.Writer, params map[string]interface{}) (modified time.Time, exists bool, err error) {
	if tpl, check := s.template(path); check {
		exists, modified, err = tpl.content(w, params, nil)
		if !exists {
			s.store.Delete(path)
		}
	}
	return
}
