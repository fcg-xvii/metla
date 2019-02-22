package metla

type storage struct {
	names []string
}

func (s *storage) appendName(name string) {
	s.names = append(s.names, name)
}
