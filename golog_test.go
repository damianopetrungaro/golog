package golog_test

// FakeErrorHandler used for internal testing purposes
type FakeErrorHandler struct {
	Err error
}

func (fe *FakeErrorHandler) Handle(err error) {
	fe.Err = err
}
