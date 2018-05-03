package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalFileWalker_NoWalkFn(t *testing.T) {
	fs := &mockFS{}
	walker := NewLocalFileWalker(fs, DefaultImagePattern)
	err := walker.Walk("test", nil)
	if err != ErrWalkFnNotSet {
		t.Fatalf("expected %v got %v", ErrWalkFnNotSet, err)
	}
}

func TestLocalFileWalker_Walk(t *testing.T) {
	fs := createMockFs()
	walker := NewLocalFileWalker(fs, DefaultImagePattern)

	times := 0
	walker.Walk("test", func(dir, name, path string) {
		if times > 0 {
			t.Fatal("expected only one item")
		}
		times++
		if dir != "asdf" {
			t.Fatalf("expected: asdf but got: %s", dir)
		}
		if name != "image" {
			t.Fatalf("expected: image but got: %s", name)
		}
		expected := filepath.Join("test", "asdf", "image.jpg")
		if expected != path {
			t.Fatalf("expected: %s but got: %s", expected, path)
		}
	})
	if times != 1 {
		t.Fatal("expected exaclty one file")
	}
}

func TestLocalFileWalker_WalkFormats(t *testing.T) {
	testcases := []struct {
		ext string
		ok  bool
	}{
		{ext: ".jpg", ok: true},
		{ext: ".jpeg", ok: true},
		{ext: ".png", ok: true},
		{ext: ".tiff", ok: true},
		{ext: ".bmp", ok: true},
		{ext: ".txt", ok: false},
		{ext: "", ok: false},
	}

	for _, tc := range testcases {
		fs := createFsWithExt(tc.ext)
		walker := NewLocalFileWalker(fs, DefaultImagePattern)
		count := 0
		walker.Walk("test", func(dir, name, path string) {
			count = 1
			if !tc.ok {
				t.Fatalf("expected %s not to be walked", tc.ext)
			}
		})
		if count == 0 && tc.ok {
			t.Fatalf("expected %s to be walked", tc.ext)
		}
	}
}

func TestLocalFileWalker_FileSystemError(t *testing.T) {
	testErr := errors.New("test")
	fs := &mockFS{
		err: testErr,
	}
	walker := NewLocalFileWalker(fs, DefaultImagePattern)
	err := walker.Walk("test", func(dir, name, path string) {})
	if err != testErr {
		t.Fatalf("expected %v but got %v", testErr, err)
	}
}

func createMockFs() *mockFS {
	return &mockFS{
		files: map[string][]os.FileInfo{
			"test": []os.FileInfo{
				&mockFile{
					name:  "asdf",
					isDir: true,
				},
				&mockFile{
					name:  "screenshot.png",
					isDir: false,
				},
			},
			filepath.Join("test", "asdf"): []os.FileInfo{
				&mockFile{
					name:  "image.jpg",
					isDir: false,
				},
				&mockFile{
					name:  "image",
					isDir: true,
				},
			},
		},
	}
}

func createFsWithExt(ext string) *mockFS {
	return &mockFS{
		files: map[string][]os.FileInfo{
			"test": []os.FileInfo{
				&mockFile{
					name:  "asdf",
					isDir: true,
				},
			},
			filepath.Join("test", "asdf"): []os.FileInfo{
				&mockFile{
					name:  "image" + ext,
					isDir: false,
				},
			},
		},
	}
}
