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
	query := flag.String("query", "What can this language do?", "capability question to investigate")
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		fail("input", fmt.Errorf("at least one .gooo declaration path is required"))
	}

	results := make([]map[string]any, 0, len(paths))
	for _, path := range paths {
		source, err := os.ReadFile(path)
		if err != nil {
			fail("read declaration "+path, err)
		}
		declaration, err := envelope.BindDeclaration(string(source))
		if err != nil {
			fail("bind declaration "+path, err)
		}
		overview, err := capability.DiscoverOverview(*query, declaration)
		if err != nil {
			fail("discover capabilities "+path, err)
		}
		if err := overview.Validate(); err != nil {
			fail("validate capability evidence "+path, err)
		}
		results = append(results, map[string]any{"path": path, "overview": overview})
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fail("encode capability overviews", err)
	}
	fmt.Println(string(output))
}

func fail(stage string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", stage, err)
	os.Exit(1)
}
