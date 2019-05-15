package metla

import (
	_ "fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fcg-xvii/containers"
)

func newTemplate(requester Requester, root *Metla, path string, content []byte, modified time.Time) *template {
	res := &template{
		requester: requester,
		root:      root,
		path:      path,
		modified:  modified,
		locker:    new(sync.RWMutex),
	}
	res.parse(content, res.root)
	return res
}

type template struct {
	requester Requester
	root      *Metla
	path      string
	commands  []executer
	store     *storage
	modified  time.Time
	err       error
	locker    *sync.RWMutex
}

func (s *template) parse(content []byte, root *Metla) error {
	parser := initParser(s.path, content, root)
	if err := parser.parseDocument(); err != nil {
		s.err = err
	} else {
		s.err, s.commands, s.store = nil, parser.execList, parser.store
	}
	return s.err
}

func (s *template) initExec(w io.Writer, parent *tplExec) *tplExec {
	res := &tplExec{
		tplName:  s.path,
		execList: make([]executer, len(s.commands)),
		writer:   w,
		sto:      s.store.execStorage(params),
		stack:    containers.NewStack(5),
		modified: s.modified,
		execStop: time.Now().Add(time.Duration(atomic.LoadInt64(&s.root.maxExecDuration))),
	}
	copy(res.execList, s.commands)
	if parent != nil {
		res.execStop = parent.execStop
		res.layout = parent.layout + 1
	}
	return res
}

func (s *template) content(w io.Writer, params map[string]interface{}, parent *tplExec) (exists bool, modified time.Time, err error) {
	var content []byte
	if content, modified, exists, err = s.requester.RequestUpdate(s.path, s.modified); s.modified.Equal(modified) {
		s.locker.RLock()
		if s.err != nil {
			err = s.err
		} else {

			ex := s.initExec(w, parent)
			s.locker.RUnlock()
			modified = ex.exec()
			return
		}
		s.locker.RUnlock()
	} else if exists {
		s.locker.Lock()
		s.modified, s.err = modified, err
		if s.err == nil {
			if err = s.parse(content, s.root); err == nil {
				ex := s.initExec(w, parent)
				s.locker.Unlock()
				modified = ex.exec()
				return
			}
		}
		s.locker.Unlock()
	}
	return
}

func (s *template) contentWithoutUpdate(w io.Writer, params map[string]interface{}) (modified time.Time, err error) {
	s.locker.RLock()
	if s.err != nil {
		err = s.err
	} else {
		ex := s.initExec(w, nil)
		s.locker.RUnlock()
		modified = ex.exec()
		return
	}
	s.locker.RUnlock()
	return
}

//////////////////////////////////////////////////////////////////////////

type tplExec struct {
	tplName   string
	execList  []executer
	writer    io.Writer
	sto       *execStorage
	stack     *containers.Stack
	modified  time.Time
	breakFlag bool
	layout    byte
	execStop  time.Time
}

func (s *tplExec) Write(data []byte) *execError {
	_, err := s.writer.Write(data)
	if err != nil {
		return &execError{s.tplName, 0, 0, err.Error()}
	}
	return nil
}

func (s *tplExec) exec() (modified time.Time) {
	if s.layout > 200 {
		s.writer.Write([]byte("Fatal error :: include loop - stack owerflow, include layouts > 200\n"))
		return
	}
	//fmt.Println("EXEC.....", s.execList)
	for _, v := range s.execList {
		//fmt.Println(v)
		if execErr := v.exec(s); execErr != nil {
			s.Write([]byte(execErr.Error()))
			return
		}
	}
	return
}
