package log

type MockLog struct {
	MockInfo       func(string)
	MockInfof      func(string, ...any)
	MockTrace      func(string)
	MockError      func(error)
	MockReadLogs   func() (string, error)
	MockReadTraces func() (string, error)
}

func (l *MockLog) Info(msg string) {
	l.MockInfo(msg)
}

func (l *MockLog) Infof(pattern string, values ...any) {
	l.MockInfof(pattern, values...)
}

func (l *MockLog) Trace(msg string) {
	l.MockTrace(msg)
}

func (l *MockLog) Error(err error) {
	l.MockError(err)
}

func (l *MockLog) ReadLogs() (string, error) {
	return l.MockReadLogs()
}

func (l *MockLog) ReadTraces() (string, error) {
	return l.MockReadTraces()
}
