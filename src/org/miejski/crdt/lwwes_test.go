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

func TestMergeOperation(t *testing.T) {
	lwwes := CreateLwwes()
	now := time.Now()
	lwwes.Add(el(1), now.Add(time.Minute))
	lwwes.Add(el(2), now.Add(2 * time.Minute))
	lwwes.Add(el(3), now.Add(3 * time.Minute))
	lwwes.Add(el(6), now.Add(7 * time.Minute))
	lwwes.Remove(el(3), now.Add(8 * time.Minute))

	lwwes2 := CreateLwwes()
	lwwes2.Add(el(4), now.Add(2 * time.Minute))
	lwwes2.Add(el(5), now.Add(3 * time.Minute))
	lwwes2.Add(el(6), now.Add(5 * time.Minute))
	lwwes2.Remove(el(6), now.Add(6 * time.Minute))
	lwwes2.Remove(el(2), now.Add(7 * time.Minute))

	merged := lwwes.Merge(&lwwes2)

	assert.True(t, merged.Contains(el(1)))
	assert.False(t, merged.Contains(el(2)))
	assert.False(t, merged.Contains(el(3)))
	assert.True(t, merged.Contains(el(4)))
	assert.True(t, merged.Contains(el(5)))
	assert.True(t, merged.Contains(el(6)))
}



func el(i int) IntElement {
	return IntElement{i}
}
