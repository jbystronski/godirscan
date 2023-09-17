package common

type Entry struct {
	name string
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) SetName(n string) {
	e.name = n
}
