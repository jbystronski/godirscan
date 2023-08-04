package utils

import "github.com/jbystronski/godirscan/pkg/entry"

func selectEntry(e *entry.Entry, container *map[*entry.Entry]struct{}) {
	if _, ok := (*container)[e]; ok {
		delete(*container, e)
	} else {
		(*container)[e] = struct{}{}
	}
}

func selectAllEntries(entries []*entry.Entry, container *map[*entry.Entry]struct{}) {
	if len(*container) > 0 {
		*container = make(map[*entry.Entry]struct{})
	} else {
		for _, entry := range entries {
			(*container)[entry] = struct{}{}
		}
	}
}
