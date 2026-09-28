package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
	"github.com/kimjooyoon/jev-gooo/envelope"
)

func main() {
	query := flag.String("query", "What can this language do?", "natural-language capability question")
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		fatal(fmt.Errorf("at least one .gooo declaration path is required"))
	}

	items := make([]capability.AssistantObservationItem, 0, len(paths))
	for _, path := range paths {
		source, err := os.ReadFile(path)
		if err != nil {
			fatal(fmt.Errorf("read declaration %s: %w", path, err))
		}
		declaration, err := envelope.BindDeclaration(string(source))
		if err != nil {
			fatal(fmt.Errorf("bind declaration %s: %w", path, err))
		}
		assistant, err := capability.DiscoverAssistant(*query, declaration)
		if err != nil {
			fatal(fmt.Errorf("discover capabilities %s: %w", path, err))
		}
		items = append(items, capability.AssistantObservationItem{Path: path, Assistant: assistant})
	}

	report, err := capability.BuildAssistantObservationReport(*query, items)
	if err != nil {
		fatal(err)
	}
	if err := report.Validate(); err != nil {
		fatal(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
