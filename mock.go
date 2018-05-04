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

type mockReporter struct {
	err        error
	title      string
	onAddImage func(path, category, name string)
}

func (reporter *mockReporter) SetTitle(title string) {
	reporter.title = title
}

func (reporter *mockReporter) AddImage(path string, category string, name string) {
	reporter.onAddImage(path, category, name)
}

func (reporter *mockReporter) Flush() error {
	return reporter.err
}

type mockWalker struct {
	onWalk func(fn WalkFn)
	err    error
}

func (walker *mockWalker) Walk(path string, fn WalkFn) error {
	walker.onWalk(fn)
	return walker.err
}
