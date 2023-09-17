package common

type Entries []*StoreItem

func (e *Entries) All() []*StoreItem {
	return *e
}

func (e *Entries) Len() int {
	return len(*e)
}

func (e *Entries) Find(index int) *StoreItem {
	if index > len(*e)-1 {
		return nil
	}

	return (*e)[index]
}
