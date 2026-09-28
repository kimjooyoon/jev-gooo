package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

type input struct {
	Query       string `json:"query"`
	Declaration string `json:"declaration"`
}

func main() {
	var request input
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal(err)
	}
	analysis, err := capability.AnalyzeSource(request.Query, request.Declaration)
	if err != nil {
		fatal(err)
	}
	if err := analysis.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(analysis); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
