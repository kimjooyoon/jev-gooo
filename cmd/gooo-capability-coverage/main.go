package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
	"github.com/kimjooyoon/jev-gooo/envelope"
)

type input struct { Query string `json:"query"`; Declaration string `json:"declaration"` }
type output struct { Options capability.Options `json:"options"`; Coverage capability.Coverage `json:"coverage"` }

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1]); if err != nil { fatal(err) }
		defer file.Close(); reader = file
	}
	var request input
	if err := json.NewDecoder(reader).Decode(&request); err != nil { fatal(err) }
	declaration, err := envelope.BindDeclaration(request.Declaration); if err != nil { fatal(err) }
	options, err := capability.DiscoverOptions(request.Query, declaration); if err != nil { fatal(err) }
	coverage, err := capability.MeasureCoverage(options); if err != nil { fatal(err) }
	if err := coverage.Validate(); err != nil { fatal(err) }
	if err := json.NewEncoder(os.Stdout).Encode(output{Options: options, Coverage: coverage}); err != nil { fatal(err) }
}

func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }