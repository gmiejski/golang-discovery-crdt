package crdt

import "testing"

func TestSimpleSet_Add(t *testing.T) {
	set := SimpleSet{map[Element]int{}}

	//when
	set.Add(IntElement{1})
	set.Add(IntElement{1})
	set.Add(IntElement{2})


	if set.Size() != 2 {
		t.Fail()
	}
	if !set.Contains(IntElement{2}) {
		t.Fail()
	}
	if !set.Contains(IntElement{1}) {
		t.Fail()
	}

	removed := set.Remove(IntElement{1})
	if !removed {
		t.Fail()
	}
	if set.Contains(IntElement{1}) {
		t.Fail()
	}
	if set.Size() != 1 {
		t.Fail()
	}
}

