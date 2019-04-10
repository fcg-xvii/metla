package prod

import (
	"fmt"
	"time"
)

type value struct {
	layout int
	key    string
}

type store struct {
	list   []*value
	layout int
}

func (s *storage) findObject(key string) int {
	for i := len(s.list) - 1; i >= 0; i-- {
		if key = s.list[i].key {
			return i
		}
	}
	val := &value{s.layout, false, key, nil}
	s.list = append(s.list, val)
	return len(s.list) - 1
}

/////////////////////////////////////////////

type execStorage struct {
	keys []string
	values []interface{}
}
