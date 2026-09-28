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

type output struct {
	Discovery capability.Discovery `json:"discovery"`
	Plan      capability.Plan      `json:"plan"`
}

func main() {
	var request input
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal(err)
	}
	discovery, err := capability.DiscoverFromSource(request.Query, request.Declaration)
	if err != nil {
		fatal(err)
	}
	plan, err := capability.BuildPlan(discovery)
	if err != nil {
		fatal(err)
	}
	if err := plan.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(output{Discovery: discovery, Plan: plan}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
