package main

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// File ending regex patterns
const (
	DefaultImagePattern = `.*\.(png|jpg|jpeg|tiff|bmp)`
)

// FileWalker errors
var (
	ErrWalkFnNotSet = errors.New("walker function not set")
)

// FileWalker represents a walkable filesystem API
type FileWalker interface {
	Walk(path string, fn WalkFn) error
}

// WalkFn represents a filesystem walk function
// each step WalkFn reveives the current directory name
// the current file name as well as the full path
type WalkFn func(dir, name, path string)

// LocalFileWalker offers a walkable filesystem API for local directories.
// Only files in the first subdirectory level are walked over.
type LocalFileWalker struct {
	filesystem FileSystem
	extensions *regexp.Regexp
}

// NewLocalFileWalker creates a new LocalFileSystem
func NewLocalFileWalker(
	filesystem FileSystem,
	extensionRegex string,
) *LocalFileWalker {
	fw := &LocalFileWalker{}
	fw.filesystem = filesystem
	fw.extensions = regexp.MustCompile(extensionRegex)
	return fw
}

// Walk starts the filesystem walk, calling fn for every file.
func (fw *LocalFileWalker) Walk(base string, fn WalkFn) error {
	if fn == nil {
		return ErrWalkFnNotSet
	}
	dirs, err := fw.listDirs(base)
	if err != nil {
		return err
	}
	return fw.walkDirs(base, dirs, fn)
}

func (fw *LocalFileWalker) walkDirs(base string, dirs []string, fn WalkFn) error {
	for _, dir := range dirs {
		files, err := fw.listFiles(filepath.Join(base, dir))
		if err != nil {
			return err
		}
		if err = fw.walkFiles(base, dir, files, fn); err != nil {
			return err
		}
	}
	return nil
}

func (fw *LocalFileWalker) walkFiles(base, dir string, files []string, fn WalkFn) error {
	for _, file := range files {
		name := strings.TrimSuffix(file, filepath.Ext(file))
		fn(dir, name, filepath.Join(base, dir, file))
	}
	return nil
}

func (fw *LocalFileWalker) listDirs(path string) ([]string, error) {
	dirs := make([]string, 0)
	content, err := fw.filesystem.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, item := range content {
		if item.IsDir() {
			dirs = append(dirs, item.Name())
		}
	}
	return dirs, nil
}

func (fw *LocalFileWalker) listFiles(path string) ([]string, error) {
	dirs := make([]string, 0)
	content, err := fw.filesystem.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, item := range content {
		if !item.IsDir() && fw.extensions.MatchString(item.Name()) {
			dirs = append(dirs, item.Name())
		}
	}
	return dirs, nil
}
