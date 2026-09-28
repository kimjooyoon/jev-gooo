package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

type input struct {
	Coverage capability.Coverage         `json:"coverage"`
	Binding  capability.CoverageFeedback `json:"binding"`
	Plan     capability.FocusPlan        `json:"plan"`
}

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fatal(err)
		}
		defer file.Close()
		reader = file
	}
	var request input
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		fatal(err)
	}
	ledger, err := capability.BuildImprovementLedger(request.Coverage, request.Binding, request.Plan)
	if err != nil {
		fatal(err)
	}
	if err := ledger.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(ledger); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
