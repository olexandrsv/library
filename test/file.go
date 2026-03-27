package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"
)

type MockMultipartFile struct {
	Name    string
	Content string
}

func MockMultipartFiles(t testing.TB, fieldName string, files []*MockMultipartFile) []*multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, file := range files {
		part, err := writer.CreateFormFile(fieldName, file.Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = part.Write([]byte(file.Content))
		if err != nil {
			t.Fatal("fatal:", err)
		}
	}

	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		t.Fatal(err)
	}

	return form.File[fieldName]
}

func CheckMultipartFiles(
	t testing.TB,
	fieldName string,
	received []*multipart.FileHeader,
	expected []*MockMultipartFile) {
	if len(received) != len(expected) {
		t.Errorf("received '%d' and expected '%d' files len does not match", len(received), len(expected))
		return
	}
	for i := range received {
		receivedFile := received[i]
		expectedFile := expected[i]
		CheckMultipartFile(t, fieldName, receivedFile, expectedFile)
	}
}

func CheckMultipartFile(t testing.TB, fieldName string, received *multipart.FileHeader,
	expected *MockMultipartFile) {
	if received == nil && expected == nil {
		return
	}
	if received == nil || expected == nil {
		t.Errorf("expected file %v, received file %v", expected, received)
		return
	}
	fileName, content, err := getData(received)
	if err != nil {
		t.Fatal(err)
	}
	if fileName != expected.Name || content != expected.Content {
		t.Errorf("expected file with name '%s' and content '%s', but received file with name '%s' and content '%s'",
			fileName, content, expected.Name, expected.Content)
	}
}

func getData(fileHeader *multipart.FileHeader) (string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", err
	}
	return fileHeader.Filename, string(bytes), nil
}
