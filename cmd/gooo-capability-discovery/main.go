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
	declarationPath := flag.String("declaration", "examples/capability-discovery.gooo", "path to the .gooo declaration")
	flag.Parse()

	source, err := os.ReadFile(*declarationPath)
	if err != nil {
		fail("read declaration", err)
	}
	declaration, err := envelope.BindDeclaration(string(source))
	if err != nil {
		fail("bind declaration", err)
	}
	overview, err := capability.DiscoverOverview(*query, declaration)
	if err != nil {
		fail("discover capabilities", err)
	}
	if err := overview.Validate(); err != nil {
		fail("validate capability evidence", err)
	}
	output, err := json.MarshalIndent(overview, "", "  ")
	if err != nil {
		fail("encode capability overview", err)
	}
	fmt.Println(string(output))
}

func fail(stage string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", stage, err)
	os.Exit(1)
}
