package main

import (
	"html/template"
	"io"
	"sort"
	"strings"
)

// Reporter represents a report builder
type Reporter interface {
	SetTitle(title string)
	SetVersionString(version string)
	AddImage(path string, category string, name string)
	Flush() error
}

// Report represents a comparison report
type Report struct {
	Title      string
	Version    string
	Categories []*Category
}

// Category represents a comparison category
type Category struct {
	Title  string
	Images []*Image
}

type byTitle []*Category
type byName []*Image

// Image represents a single comparison image
type Image struct {
	Path string
	Name string
}

// HTMLReporter generates a html based report
type HTMLReporter struct {
	assets     Assets
	template   string
	categories map[string]*Category
	report     *Report
	writer     io.Writer
}

// NewHTMLReporter creates a new HTML reporter.
func NewHTMLReporter(assets Assets, writer io.Writer, template string) *HTMLReporter {
	reporter := &HTMLReporter{}
	reporter.assets = assets
	reporter.report = &Report{}
	reporter.categories = make(map[string]*Category)
	reporter.writer = writer
	reporter.template = template
	return reporter
}

// SetTitle sets the main title of the Report
func (reporter *HTMLReporter) SetTitle(title string) {
	reporter.report.Title = title
}

// SetVersionString set build version information
func (reporter *HTMLReporter) SetVersionString(version string) {
	reporter.report.Version = version
}

// AddImage adds a new image with the given name to the report
// creating a new category for each title.
func (reporter *HTMLReporter) AddImage(path string, category string, name string) {
	image := &Image{
		Path: path,
		Name: name,
	}
	cat := reporter.getCategory(category)
	cat.Images = append(cat.Images, image)
}

// Flush generates the entire report writing it to the given writer
func (reporter *HTMLReporter) Flush() error {
	html := reporter.assets.String(reporter.template)
	tmpl := template.Must(template.New("base").Parse(html))
	reporter.report.Categories = reporter.getSortedCategories()
	return tmpl.Execute(reporter.writer, reporter.report)
}

func (reporter *HTMLReporter) getSortedCategories() []*Category {
	categories := make([]*Category, 0)
	for _, category := range reporter.categories {
		sort.Sort(byName(category.Images))
		categories = append(categories, category)
	}
	sort.Sort(byTitle(categories))
	return categories
}

func (reporter *HTMLReporter) getCategory(title string) *Category {
	var cat *Category
	var exists bool
	if cat, exists = reporter.categories[title]; !exists {
		cat = &Category{}
		cat.Title = title
		cat.Images = make([]*Image, 0)
		reporter.categories[title] = cat
	}
	return cat
}

func (by byTitle) Len() int           { return len(by) }
func (by byTitle) Less(i, j int) bool { return strings.Compare(by[i].Title, by[j].Title) < 0 }
func (by byTitle) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

func (by byName) Len() int           { return len(by) }
func (by byName) Less(i, j int) bool { return strings.Compare(by[i].Name, by[j].Name) < 0 }
func (by byName) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }
