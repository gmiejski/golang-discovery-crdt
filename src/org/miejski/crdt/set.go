package crdt

type Set interface {
	Add(element Element) bool
	Remove(element Element) bool
	Contains(element Element) bool
	Size() int
}

type SimpleSet struct {
	data map[Element]int
}

func (s *SimpleSet) Add(element Element) bool {
	if s.Contains(element) {
		return false
	}

	ints := (*s).data
	ints[element] = 1
	return true
}

func (s *SimpleSet) Remove(element Element) bool {
	if !s.Contains(element) {
		return false
	}
	delete(s.data, element)
	return true
}

func (s *SimpleSet) Contains(element Element) bool {
	for k, _ := range s.data {
		if k == element {
			return true
		}
	}
	return false
}

func (s *SimpleSet) Size() int {
	return len(s.data)
}
