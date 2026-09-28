package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type replayReport struct {
	Valid                bool                                      `json:"valid"`
	PriorIterationDigest string                                      `json:"prior_iteration_digest"`
	Outcome              envelope.SelfImprovementIterationReplayOutcome `json:"outcome"`
	NonExecuting         bool                                      `json:"non_executing"`
	NonAuthorizing       bool                                      `json:"non_authorizing"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: jev-iteration-replay <iteration.json> <replay.json>")
		os.Exit(2)
	}
	iterationData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	replayData, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	var iteration envelope.SelfImprovementIteration
	if err := json.Unmarshal(iterationData, &iteration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var replay envelope.SelfImprovementIterationReplay
	if err := json.Unmarshal(replayData, &replay); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := replay.ValidateAgainst(iteration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	report := replayReport{
		Valid:                true,
		PriorIterationDigest: replay.PriorIterationDigest,
		Outcome:              replay.Outcome,
		NonExecuting:         replay.NonExecuting,
		NonAuthorizing:       replay.NonAuthorizing,
	}
	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

