package crdt

import (
	"time"
	"fmt"
)

type crdt interface {
	Merge(crdt) crdt
}

type Element interface {
	Get() interface{}
}

type LastWriteWinsElementSet interface {
	Add(element Element, time time.Time) bool
	Remove(element Element, time time.Time) bool
	Get() []Element
	Merge(LastWriteWinsElementSet) LastWriteWinsElementSet
	Contains(element Element) bool
	Size() int
}

type lwwes struct {
	add_set    map[Element]time.Time
	remove_set map[Element]time.Time
}

func CreateLwwes() LastWriteWinsElementSet {
	p := lwwes{map[Element]time.Time{}, map[Element]time.Time{}}
	return &p
}

func (s *lwwes) Add(element Element, t time.Time) bool {
	state, last_operation_time := s.elementInfo(element)
	if state == ADDED {
		return false
	}

	if state == REMOVED && t.Before(last_operation_time) {
		panic(fmt.Sprintf("Adding element: #%v with time earlier than last %s operation!", element, state))
	}

	s.add_set[element] = t
	return true
}

func (s *lwwes) Remove(element Element, t time.Time) bool {
	state, last_operation_time := s.elementInfo(element)

	if state == REMOVED {
		return false
	}

	if state == ADDED && t.Before(last_operation_time) {
		panic(fmt.Sprintf("Removing element: #%v with time earlier than last %s operation!", element, state))
	}

	s.remove_set[element] = t
	return true
}

func (s *lwwes) Get() []Element {
	result := make([]Element, 0)
	for added_element, add_time := range s.add_set {
		removed_time, removed := s.remove_set[added_element]
		if !removed || add_time.After(removed_time) {
			result = append(result, added_element)
		}
	}
	return result
}

func (s *lwwes) Merge(LastWriteWinsElementSet) LastWriteWinsElementSet {
	panic("Panicked on merge")
}

func (s *lwwes) Contains(element Element) bool {
	state, _ := s.elementInfo(element)
	if state == ADDED {
		return true
	}
	return false
}

func (s *lwwes) Size() int {
	return len(s.Get())
}

type ElementState string

const (
	ABSENT  ElementState = "ABSENT"
	ADDED   ElementState = "ADDED"
	REMOVED ElementState = "REMOVED"
)

func (s *lwwes) elementInfo(el Element) (ElementState, time.Time) {
	added_time, ok := s.add_set[el]
	if !ok {
		return ABSENT, time.Time{}
	}
	removed_time, removed_ok := s.remove_set[el]
	if !removed_ok || added_time.After(removed_time) {
		return ADDED, added_time
	}
	return REMOVED, removed_time
}
