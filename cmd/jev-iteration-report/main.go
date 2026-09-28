package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type iterationReport struct {
	Valid          bool                                      `json:"valid"`
	IterationID    string                                    `json:"iteration_id"`
	State          envelope.SelfImprovementIterationState   `json:"state"`
	CanonicalDigest string                                   `json:"canonical_digest"`
	FirstMissingStage string                                `json:"first_missing_stage,omitempty"`
	NextOperation  string                                    `json:"next_operation,omitempty"`
	NonExecuting   bool                                      `json:"non_executing"`
	NonAuthorizing bool                                      `json:"non_authorizing"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: jev-iteration-report <iteration.json>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var iteration envelope.SelfImprovementIteration
	if err := json.Unmarshal(data, &iteration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := iteration.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	digest, err := iteration.CanonicalDigest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	report := iterationReport{
		Valid:             true,
		IterationID:       iteration.IterationID,
		State:             iteration.State,
		CanonicalDigest:   digest,
		FirstMissingStage: iteration.FirstMissingStage,
		NextOperation:     iteration.NextOperation,
		NonExecuting:      iteration.NonExecuting,
		NonAuthorizing:    iteration.NonAuthorizing,
	}
	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

