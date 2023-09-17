package common

type Store struct {
	DataAccessor
	name string
}

func (s *Store) Name() string {
	return s.name
}

func (s *Store) SetName(name string) {
	s.name = name
}

func (s *Store) Data() DataAccessor {
	return s.DataAccessor
}

func (s *Store) SetData(acc DataAccessor) {
	s.DataAccessor = acc
}

func (s *Store) Items() DataAccessor {
	return s.DataAccessor
}

func (s *Store) SetItems(acc DataAccessor) {
	s.DataAccessor = acc
}
