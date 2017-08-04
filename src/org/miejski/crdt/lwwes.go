package crdt

import (
	"time"
	"fmt"
	"io"
	"encoding/json"
)

type crdt interface {
	Merge(crdt) crdt
}

type Element interface {
	Get() string
}

type LastWriteWinsElementSet interface {
	Add(element string, time time.Time) bool
	Remove(element string, time time.Time) bool
	Get() []string
	Merge(other LastWriteWinsElementSet) LastWriteWinsElementSet
	Contains(element string) bool
	Size() int
}

type Lwwes struct {
	Add_set    map[string]time.Time
	Remove_set map[string]time.Time
}

func (lwwes *Lwwes) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	defer data.Close()
	err := decoder.Decode(lwwes)
	if err != nil {
		return err
	}
	return nil
}

func CreateLwwes() Lwwes {
	p := Lwwes{map[string]time.Time{}, map[string]time.Time{}}
	return p
}

func (s *Lwwes) Add(element string, t time.Time) bool {
	state, last_operation_time := s.elementInfo(element)
	if state == ADDED {
		return false
	}

	if state == REMOVED && t.Before(last_operation_time) {
		panic(fmt.Sprintf("Adding element: #%v with time earlier than last %s operation!", element, state))
	}

	s.Add_set[element] = t
	return true
}

func (s *Lwwes) Remove(element string, t time.Time) bool {
	state, last_operation_time := s.elementInfo(element)

	if state == REMOVED {
		return false
	}

	if state == ADDED && t.Before(last_operation_time) {
		panic(fmt.Sprintf("Removing element: #%v with time earlier than last %s operation!", element, state))
	}

	s.Remove_set[element] = t
	return true
}

func (s *Lwwes) GetNotWorkingImplementation() []string {
	result := make([]string, 0)
	for added_element, add_time := range s.Add_set {
		removed_time, removed := s.Remove_set[added_element]
		if !removed || add_time.After(removed_time) {
			result = append(result, added_element)
		}
	}
	return result
}

func (s Lwwes) Get() []string {
	result := make([]string, 0)
	for added_element := range s.Add_set {
		state, _ := s.elementInfo(added_element)
		if state == ADDED {
			result = append(result, added_element)
		} else {
			fmt.Sprintf("") // TODO remove
		}
	}
	return result
}

func (s *Lwwes) Merge(other LastWriteWinsElementSet) LastWriteWinsElementSet {
	casted, ok := other.(*Lwwes)
	if !ok {
		return nil
	}
	merged_insert := mergeMap(s.Add_set, casted.Add_set)
	merged_remove := mergeMap(s.Remove_set, casted.Remove_set)
	result := Lwwes{merged_insert, merged_remove}
	return &result
}

func mergeMap(m1 map[string]time.Time, m2 map[string]time.Time) map[string]time.Time {
	result := map[string]time.Time{}
	for m1_element, m1_time := range m1 {
		last_observed, present := result[m1_element]
		if !present || last_observed.Before(m1_time) {
			result[m1_element] = m1_time
		}
	}
	for m2_element, m2_time := range m2 {
		last_observed, present := result[m2_element]
		if !present || last_observed.Before(m2_time) {
			result[m2_element] = m2_time
		}
	}
	return result
}

func (s *Lwwes) Contains(element string) bool {
	state2, _ := (*s).elementInfo(element)
	if state2 == ADDED {
		return true
	}

	state, _ := s.elementInfo(element)
	if state == ADDED {
		return true
	}
	return false
}

func (s *Lwwes) Size() int {
	return len(s.Get())
}

type ElementState string

const (
	ABSENT  ElementState = "ABSENT"
	ADDED   ElementState = "ADDED"
	REMOVED ElementState = "REMOVED"
)

func (s *Lwwes) elementInfo(el string) (ElementState, time.Time) {
	added_time, ok := s.Add_set[el]
	if !ok {
		return ABSENT, time.Time{}
	}
	removed_time, removed_ok := s.Remove_set[el]
	if !removed_ok || added_time.After(removed_time) {
		return ADDED, added_time
	}
	return REMOVED, removed_time

	//var added_time time.Time
	//var added_ok bool
	//for k, t := range s.Add_set {
	//	if k.Get() == el.Get() {
	//		added_ok = true
	//		added_time = t
	//	}
	//}
	//if !added_ok {
	//	return ABSENT, time.Time{}
	//}
	//
	//var removed_time time.Time
	//var removed_ok bool
	//for k, t := range s.Remove_set {
	//	if k.Get() == el.Get() {
	//		removed_ok = true
	//		removed_time = t
	//	}
	//}
	//if !removed_ok || added_time.After(removed_time) {
	//	return ADDED, added_time
	//}
	//return REMOVED, removed_time
}

func (s *Lwwes) elementInfoNotWorking(el string) (ElementState, time.Time) {
	added_time, ok := s.Add_set[el]
	if !ok {
		return ABSENT, time.Time{}
	}
	removed_time, removed_ok := s.Remove_set[el]
	if !removed_ok || added_time.After(removed_time) {
		return ADDED, added_time
	}
	return REMOVED, removed_time
}
