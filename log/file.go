package log

import (
	"bytes"
	"fmt"
	"io"
	"library/errors"
	"os"
	"sync"
)

type LogFile interface {
	read() (string, error)
	print(string)
}

type file struct {
	mx   sync.RWMutex
	file *os.File
}

func newFile(f *os.File) *file {
	return &file{
		file: f,
	}
}

func (f *file) print(msg string) {
	f.mx.Lock()
	defer f.mx.Unlock()

	var b bytes.Buffer
	b.WriteString(msg)
	b.WriteString("\n")
	_, err := f.file.Write(b.Bytes())
	if err != nil {
		fmt.Printf("logger.print error: %+v", err)
	}
}

func (f *file) read() (string, error) {
	f.mx.RLock()
	defer f.mx.RUnlock()
	f.file.Seek(0, 0)
	bytes, err := io.ReadAll(f.file)
	if err != nil {
		fmt.Println(err)
		return "", errors.NewFileErr(err, f.file.Name(), errors.NewInternalError())
	}
	return string(bytes), nil
}

type MockFile struct {
	mockRead  func() (string, error)
	mockPrint func(string)
}

func (mf *MockFile) read() (string, error) {
	return mf.mockRead()
}

func (mf *MockFile) print(msg string) {
	mf.mockPrint(msg)
}
