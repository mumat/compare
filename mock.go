package main

import (
	"os"
	"time"
)

type mockFS struct {
	files map[string][]os.FileInfo
	err   error
}

func (fs *mockFS) ReadDir(path string) ([]os.FileInfo, error) {
	return fs.files[path], fs.err
}

type mockFile struct {
	name  string
	isDir bool
}

func (file *mockFile) Name() string {
	return file.name
}

func (file *mockFile) Size() int64 {
	return 42
}

func (file *mockFile) Mode() os.FileMode {
	return 775
}

func (file *mockFile) ModTime() time.Time {
	return time.Now()
}

func (file *mockFile) IsDir() bool {
	return file.isDir
}

func (file *mockFile) Sys() interface{} {
	return nil
}

type mockAssets string

func (assets mockAssets) String(name string) string {
	return string(assets)
}
