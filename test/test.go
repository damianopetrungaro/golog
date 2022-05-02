package test

//go:generate mockgen --build_flags=--mod=mod -destination=./mock.go -package=test github.com/damianopetrungaro/golog Logger,CheckedLogger,Entry,Writer,Flusher
