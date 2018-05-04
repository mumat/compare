package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"
)

var assets Assets

var title string
var walker string
var filePattern string
var reporters []string
var htmlFile string
var htmlTemplate string

var version string
var date string
var hash string

var cmd = &cobra.Command{
	Use:   "compare [image directory]",
	Short: "Compare is a static image comparison generator",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run:   printVersion,
}

func init() {
	assets = packr.NewBox("./assets")

	cmd.AddCommand(versionCmd)

	cmd.Flags().StringVarP(&title, "title", "t", "Compare", "Title of the report")

	cmd.Flags().StringArrayVarP(&reporters, "reporter", "r", []string{"html"}, "Reporters to use for generation")
	cmd.Flags().StringVar(&htmlFile, "html-out", "report.html", "HTML file for HTML report")
	cmd.Flags().StringVar(&htmlTemplate, "html-template", "template.html", "HTML template to use for report generation")

	cmd.Flags().StringVarP(&walker, "walker", "w", "local", "File walker to use")
	cmd.Flags().StringVar(&filePattern, "ext-pattern", DefaultImagePattern, "File extension pattern for walker")
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	comparer := NewComparer()
	comparer.SetTitle(title)
	comparer.SetVersionString(fmt.Sprintf("v%s", version))
	filewalker, err := getWalker(walker)
	if err != nil {
		log.Fatal(err)
	}
	comparer.SetWalker(filewalker)
	for _, reporterType := range reporters {
		reporter, err := getReporter(reporterType)
		if err != nil {
			log.Fatal(err)
		}
		comparer.AddReporter(reporter)
	}
	if err := comparer.Compare(args[0]); err != nil {
		log.Fatal(err)
	}
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Compare v%s (built: %s commit: %s)\n", version, date, hash)
}

func getReporter(name string) (Reporter, error) {
	switch name {
	case "html":
		return NewHTMLReporter(assets, getOutput(htmlFile), htmlTemplate), nil
	}
	return nil, fmt.Errorf("unkown reporter %s supported are: [html]", name)
}

func getWalker(name string) (FileWalker, error) {
	switch name {
	case "local":
		return NewLocalFileWalker(NewLocalFileSystem(), filePattern), nil
	}
	return nil, fmt.Errorf("unkown walker %s supported are: [local]", name)
}

func getOutput(path string) io.WriteCloser {
	output, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return output
}
