package golog

type Writer interface {
	Write(Entry)
}
