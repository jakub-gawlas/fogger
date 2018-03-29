package notifier

type store struct {
	state map[string]string
}

func newStore() *store {
	return &store{
		state: map[string]string{},
	}
}

func (s *store) Get(key string) (string, error) {
	return s.state[key], nil
}

func (s *store) Set(key, value string) error {
	s.state[key] = value
	return nil
}
