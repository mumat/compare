package main

import (
	"bytes"
	"testing"
)

func TestHTMLReporter_SetTitle(t *testing.T) {
	var buffer bytes.Buffer
	assets := mockAssets("{{.Title}}")
	reporter := NewHTMLReporter(assets, &buffer, "")

	reporter.SetTitle("test")
	err := reporter.Flush()
	if err != nil {
		t.Fatalf("expected error to be nil %v", err)
	}
	result := buffer.String()
	if result != "test" {
		t.Fatalf("expected result to be 'test' but got '%s'", result)
	}
}

func TestHTMLReporter_AddImage(t *testing.T) {
	var buffer bytes.Buffer
	assets := mockAssets("{{range .Categories}}{{.Title}}{{range .Images}}|{{.Name}}{{end}}#{{end}}")
	reporter := NewHTMLReporter(assets, &buffer, "")

	reporter.AddImage("test/fdsa/image.jpg", "fdsa", "image")
	reporter.AddImage("test/asdf/pic.jpg", "asdf", "pic")
	reporter.AddImage("test/asdf/image.jpg", "asdf", "image")
	err := reporter.Flush()
	if err != nil {
		t.Fatalf("expected error to be nil %v", err)
	}
	result := buffer.String()
	if result != "asdf|image|pic#fdsa|image#" {
		t.Fatalf("expected result to be 'asdfimage' but got '%s'", result)
	}
}
