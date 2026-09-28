package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

type input struct {
	Feedback     capability.FeedbackInput `json:"feedback"`
	Discovery    capability.Discovery    `json:"discovery"`
	CapabilityID string                  `json:"capability_id"`
	Rationale    string                  `json:"rationale"`
}

type output struct {
	Candidate capability.ReviewCandidate `json:"candidate"`
}

func main() {
	var request input
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal(err)
	}
	feedback := capability.ObserveFeedback(request.Discovery, request.Feedback)
	candidate, err := capability.ProposeReview(feedback, request.CapabilityID, request.Rationale)
	if err != nil {
		fatal(err)
	}
	if err := candidate.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(output{Candidate: candidate}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
