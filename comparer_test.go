package main

import (
	"errors"
	"testing"
)

func TestComparer_Compare(t *testing.T) {
	compare := NewComparer()
	added := false
	reporter := &mockReporter{
		err: nil,
		onAddImage: func(path, category, name string) {
			added = true
			if path != "test/img" || category != "test" || name != "img" {
				t.Fatalf("expected (test/img, test, img) got (%s, %s, %s)", path, category, name)
			}
		},
	}
	walker := &mockWalker{
		err: nil,
		onWalk: func(fn WalkFn) {
			fn("test", "img", "test/img")
		},
	}

	compare.AddReporter(reporter)
	compare.SetWalker(walker)
	err := compare.Compare("test")
	if err != nil {
		t.Fatalf("expected error to be nil: %v", err)
	}
	if !added {
		t.Fatalf("add image was not called")
	}
}

func TestComparer_SetTitle(t *testing.T) {
	compare := NewComparer()
	reporter := &mockReporter{
		onAddImage: func(path, category, name string) {},
	}
	walker := &mockWalker{
		err: nil,
		onWalk: func(fn WalkFn) {
			fn("test", "img", "test/img")
		},
	}
	compare.SetWalker(walker)
	compare.AddReporter(reporter)
	compare.SetTitle("Test")
	err := compare.Compare("test")
	if err != nil {
		t.Fatalf("expected error to be nil: %v", err)
	}
	if reporter.title != "Test" {
		t.Fatalf("expected: Test got: %s", reporter.title)
	}
}

func TestComparer_CompareWithoutWalker(t *testing.T) {
	compare := NewComparer()
	err := compare.Compare("test")
	if err != ErrNoWalkerSet {
		t.Fatalf("expected %v but got %v", ErrNoWalkerSet, err)
	}
}

func TestComparer_WalkError(t *testing.T) {
	testErr := errors.New("test")
	compare := NewComparer()
	walker := &mockWalker{
		err:    testErr,
		onWalk: func(fn WalkFn) {},
	}
	compare.SetWalker(walker)
	err := compare.Compare("test")
	if err != testErr {
		t.Fatalf("expected %v got %v", testErr, err)
	}
}

func TestComparer_FlushError(t *testing.T) {
	testErr := errors.New("test")
	compare := NewComparer()
	walker := &mockWalker{
		onWalk: func(fn WalkFn) {
			fn("test", "img", "test/img")
		},
	}
	reporter := &mockReporter{
		err:        testErr,
		onAddImage: func(path, category, name string) {},
	}
	compare.SetWalker(walker)
	compare.AddReporter(reporter)
	err := compare.Compare("test")
	if err != testErr {
		t.Fatalf("expected %v got %v", testErr, err)
	}
}
