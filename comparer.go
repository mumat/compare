package main

import "errors"

// Comparer errors
var (
	ErrNoWalkerSet = errors.New("no filewalker set")
)

// Comparer uses a filewalker to create reports with registered reporters
type Comparer struct {
	filewalker FileWalker
	reporters  []Reporter

	title   string
	version string
}

// NewComparer creates a new Comparer
func NewComparer() *Comparer {
	comparer := &Comparer{}
	comparer.reporters = make([]Reporter, 0)
	return comparer
}

// SetTitle sets the title on add reporters
func (comparer *Comparer) SetTitle(title string) {
	comparer.title = title
}

// SetVersionString sets build version information
func (comparer *Comparer) SetVersionString(version string) {
	comparer.version = version
}

// AddReporter adds a reporter that should be used for report generation
func (comparer *Comparer) AddReporter(reporter Reporter) {
	comparer.reporters = append(comparer.reporters, reporter)
}

// SetWalker sets the walker to use for image collection
func (comparer *Comparer) SetWalker(walker FileWalker) {
	comparer.filewalker = walker
}

// Compare walks over given path and creates the reports
func (comparer *Comparer) Compare(path string) error {
	if comparer.filewalker == nil {
		return ErrNoWalkerSet
	}
	comparer.setTitle(comparer.title)
	err := comparer.filewalker.Walk(path, comparer.onAddImage)
	if err != nil {
		return err
	}
	return comparer.flush()
}

func (comparer *Comparer) setTitle(title string) {
	for _, reporter := range comparer.reporters {
		reporter.SetTitle(title)
		reporter.SetVersionString(comparer.version)
	}
}

func (comparer *Comparer) onAddImage(dir, name, path string) {
	for _, reporter := range comparer.reporters {
		reporter.AddImage(path, dir, name)
	}
}

func (comparer *Comparer) flush() error {
	for _, reporter := range comparer.reporters {
		if err := reporter.Flush(); err != nil {
			return err
		}
	}
	return nil
}
