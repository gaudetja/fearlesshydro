package main

type Fan interface {
	On() error
	Off() error
	Running() bool
	String() string
}
type MockFan struct {
	running bool
}

func (f *MockFan) On() error {
	f.running = true
	return nil
}

func (f *MockFan) Off() error {
	f.running = false
	return nil
}

func (f *MockFan) Running() bool {
	return f.running
}

func (f *MockFan) String() string {
	if f.running {
		return "ON"
	}
	return "OFF"
}

func MakeFan() (Fan, error) {
	return &MockFan{running: false}, nil
}
