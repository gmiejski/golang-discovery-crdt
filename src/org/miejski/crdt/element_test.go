package crdt

import "strconv"

type IntElement struct{
	Val int
}

func (i IntElement) Get() string {
	return strconv.Itoa(i.Val)
}