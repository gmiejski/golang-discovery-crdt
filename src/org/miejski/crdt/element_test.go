package crdt

type IntElement struct{
	val int
}

func (i IntElement) Get() interface{} {
	return i.val
}