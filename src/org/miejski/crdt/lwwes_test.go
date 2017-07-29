package crdt

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestAddingRemovingElements(t *testing.T) {
	lwwes := CreateLwwes()
	lwwes.Add(el(1), time.Now())

	if !lwwes.Contains(el(1)) || lwwes.Size() != 1 {
		panic(t)
	}

	lwwes.Remove(el(1), time.Now())

	if lwwes.Contains(el(1)) || lwwes.Size() != 0 {
		panic(t)
	}
}

func TestOperationsHappeningInPast(t *testing.T) {
	assert.Panics(t, func() {
		lwwes := CreateLwwes()
		now := time.Now()
		lwwes.Add(el(1), now)
		lwwes.Remove(el(1), now.Add(-1*time.Minute))
	})

	assert.Panics(t, func() {
		lwwes := CreateLwwes()
		now := time.Now()
		lwwes.Add(el(1), now)
		lwwes.Remove(el(1), now.Add(time.Minute))
		lwwes.Add(el(1), now.Add(-time.Millisecond))
	})
}

func el(i int) IntElement {
	return IntElement{i}
}
