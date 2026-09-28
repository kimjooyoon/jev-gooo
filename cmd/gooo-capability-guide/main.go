package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
	"github.com/kimjooyoon/jev-gooo/envelope"
)

type input struct {
	Query       string `json:"query"`
	Declaration string `json:"declaration"`
}

type output struct {
	Discovery capability.Discovery `json:"discovery"`
	Plan      capability.Plan      `json:"plan"`
	Guide     capability.Guide      `json:"guide"`
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
	declaration, err := envelope.BindDeclaration(request.Declaration)
	if err != nil {
		fatal(err)
	}
	discovery, err := capability.Discover(request.Query, declaration)
	if err != nil {
		fatal(err)
	}
	plan, err := capability.BuildPlan(discovery)
	if err != nil {
		fatal(err)
	}
	guide, err := capability.BuildGuide(discovery, plan)
	if err != nil {
		fatal(err)
	}
	if err := guide.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(output{Discovery: discovery, Plan: plan, Guide: guide}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}