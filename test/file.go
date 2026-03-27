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

func MockMultipartFiles(t *testing.T, fieldName string, files []*MockMultipartFile) []*multipart.FileHeader {
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
	t *testing.T,
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
		if receivedFile == nil && expectedFile == nil {
			return
		}
		if receivedFile == nil || expectedFile == nil {
			t.Errorf("expected file %v, received file %v", expectedFile, receivedFile)
			return
		}
		file, err := receivedFile.Open()
		if err != nil {
			t.Fatal(err)
		}
		bytes, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		content := string(bytes)
		if receivedFile.Filename != expectedFile.Name || content != expectedFile.Content {
			t.Errorf("expected file with name '%s' and content '%s', but received file with name '%s' and content '%s'",
				receivedFile.Filename, content, expectedFile.Name, expectedFile.Content)
		}
	}
}
