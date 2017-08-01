package simple_json

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
)

func TestUnmarshal(t *testing.T) {
	t1 := TestObject{1, "asd"}
	value, _ := Marshal(t1)

	// when
	var t2 TestObject
	readCloser:= &MockReadCloser{value}

	Unmarshal(readCloser, &t2)

	// then
	assert.True(t, t1 == t2)
}

type TestObject struct {
	Field1 int
	Field2 string
}

func (t *TestObject) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	defer (data).Close()
	err := decoder.Decode(t)
	if err != nil {
		return err
	}
	return nil
}

func (t TestObject) ToJson() string {
	val, err := json.Marshal(t)

	if err != nil {
		panic(err)
	}
	return string(val)
}

type MockReadCloser struct {
	data string
}

func (c *MockReadCloser) Read(p []byte) (n int, err error) {
	return copy(p, c.data), nil
}

func (*MockReadCloser) Close() error {
	return nil
}
