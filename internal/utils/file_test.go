package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Create a temporary file
	content := "test content"
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading existing file
	loadedContent, err := LoadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadFile() error = %v", err)
	}
	if loadedContent != content {
		t.Errorf("LoadFile() = %v, want %v", loadedContent, content)
	}

	// Test loading non-existent file
	_, err = LoadFile("non_existent_file.txt")
	if err != os.ErrNotExist {
		t.Errorf("LoadFile() error = %v, want %v", err, os.ErrNotExist)
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Test writing to a new file
	content := "test content"
	filename := filepath.Join(tmpdir, "testfile.txt")
	err = WriteFile(filename, content)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	// Verify the file was created with correct content
	loadedContent, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(loadedContent) != content {
		t.Errorf("WriteFile() wrote %v, want %v", string(loadedContent), content)
	}

	// Test writing to a file in a non-existent directory
	newFilename := filepath.Join(tmpdir, "newdir", "testfile.txt")
	err = WriteFile(newFilename, content)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	// Verify the file was created in the new directory
	if _, err := os.Stat(newFilename); os.IsNotExist(err) {
		t.Errorf("WriteFile() did not create file in new directory")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Test existing file
	exists, err := FileExists(tmpfile.Name())
	if err != nil {
		t.Errorf("FileExists() error = %v", err)
	}
	if !exists {
		t.Errorf("FileExists() = %v, want true", exists)
	}

	// Test non-existent file
	exists, err = FileExists("non_existent_file.txt")
	if err != nil {
		t.Errorf("FileExists() error = %v", err)
	}
	if exists {
		t.Errorf("FileExists() = %v, want false", exists)
	}
}
