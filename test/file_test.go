package test

import (
	"mime/multipart"
	"testing"
)

func TestCheckMultipartFile(t *testing.T) {
	files := []*MockMultipartFile{{Name: "a.txt", Content: "abc"}}
	receivedFile := MockMultipartFiles(t, "file", files)[0]
	table := []struct {
		name      string
		received  *multipart.FileHeader
		expected  *MockMultipartFile
		wantErr   bool
		wantFatal bool
	}{
		{"both nil", nil, nil, false, false},
		{"received nil, expected not nil", nil, &MockMultipartFile{Name: "a.txt", Content: "abc"}, true, false},
		{"expected nil, received not nil", &multipart.FileHeader{Filename: "a.txt"}, nil, true, false},
		{"name and content match", receivedFile, files[0], false, false},
		{"name mismatch", receivedFile, &MockMultipartFile{Name: "b.txt", Content: "abc"}, true, false},
		{"content mismatch", receivedFile, &MockMultipartFile{Name: "a.txt", Content: "def"}, true, false},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				CheckMultipartFile(m, "file", tc.received, tc.expected)
			}()
			if m.fataled != tc.wantFatal {
				t.Errorf("wantFatal=%v got=%v", tc.wantFatal, m.fataled)
			}
			if !m.fataled && m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestCheckMultipartFiles_MockTB(t *testing.T) {
	files := []*MockMultipartFile{{Name: "a.txt", Content: "abc"}, {Name: "b.txt", Content: "def"}}
	received := MockMultipartFiles(t, "file", files)
	table := []struct {
		name      string
		received  []*multipart.FileHeader
		expected  []*MockMultipartFile
		wantErr   bool
		wantFatal bool
	}{
		{"all match", received, files, false, false},
		{"length mismatch", received, files[:1], true, false},
		{"content mismatch", received, []*MockMultipartFile{{Name: "a.txt", Content: "abc"}, {Name: "b.txt", Content: "xyz"}}, true, false},
		{"all nil", nil, nil, false, false},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				CheckMultipartFiles(m, "file", tc.received, tc.expected)
			}()
			if m.fataled != tc.wantFatal {
				t.Errorf("wantFatal=%v got=%v", tc.wantFatal, m.fataled)
			}
			if !m.fataled && m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestMockMultipartFiles(t *testing.T) {
	table := []struct {
		name          string
		files         []*MockMultipartFile
		expectedFiles []*MockMultipartFile
	}{
		{
			name:          "single file",
			files:         []*MockMultipartFile{{Name: "a.txt", Content: "abc"}},
			expectedFiles: []*MockMultipartFile{{Name: "a.txt", Content: "abc"}},
		},
		{
			name:          "multiple files",
			files:         []*MockMultipartFile{{Name: "a.txt", Content: "abc"}, {Name: "b.txt", Content: "def"}},
			expectedFiles: []*MockMultipartFile{{Name: "a.txt", Content: "abc"}, {Name: "b.txt", Content: "def"}},
		},
		{
			name:          "empty files",
			files:         []*MockMultipartFile{},
			expectedFiles: []*MockMultipartFile{},
		},
		{
			name:          "nil files",
			files:         nil,
			expectedFiles: nil,
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			result := MockMultipartFiles(t, "file", tc.files)
			if len(result) != len(tc.expectedFiles) {
				t.Errorf("expected %d file(s), got %d", len(tc.expectedFiles), len(result))
			}
			for i := range result {
				receivedFile := result[i]
				expectedFile := tc.expectedFiles[i]
				if receivedFile.Filename != tc.expectedFiles[i].Name {
					t.Errorf("expected filename '%s', got '%s'", expectedFile.Name, receivedFile.Filename)
				}
				name, content, err := getData(receivedFile)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if name != expectedFile.Name || content != expectedFile.Content {
					t.Errorf("expected name '%s' and content '%s', got name '%s' and content '%s'", tc.expectedFiles[i].Name, tc.expectedFiles[i].Content, name, content)
				}
			}
		})
	}
}
