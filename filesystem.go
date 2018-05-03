package main

import (
	"io/ioutil"
	"os"
)

// FileSystem offers a simple filesystem api
type FileSystem interface {
	ReadDir(path string) ([]os.FileInfo, error)
}

// LocalFileSystem offers access to the local disk
type LocalFileSystem struct {
}

// NewLocalFileSystem creates a new LocalFileSystem
func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

// ReadDir list the contents of the given path
func (fs *LocalFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
