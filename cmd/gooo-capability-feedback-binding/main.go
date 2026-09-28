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
	Query         string `json:"query"`
	Declaration   string `json:"declaration"`
	Source        string `json:"source"`
	EvidenceDigest string `json:"evidence_digest"`
	Verified      bool   `json:"verified"`
}

type output struct {
	Options  capability.Options         `json:"options"`
	Coverage capability.Coverage         `json:"coverage"`
	Feedback capability.Feedback        `json:"feedback"`
	Binding  capability.CoverageFeedback `json:"binding"`
}

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1]); if err != nil { fatal(err) }
		defer file.Close(); reader = file
	}
	var request input
	if err := json.NewDecoder(reader).Decode(&request); err != nil { fatal(err) }
	declaration, err := envelope.BindDeclaration(request.Declaration); if err != nil { fatal(err) }
	discovery, err := capability.Discover(request.Query, declaration); if err != nil { fatal(err) }
	options, err := capability.DiscoverOptions(request.Query, declaration); if err != nil { fatal(err) }
	coverage, err := capability.MeasureCoverage(options); if err != nil { fatal(err) }
	feedback := capability.ObserveFeedback(discovery, capability.FeedbackInput{Source: request.Source, EvidenceDigest: request.EvidenceDigest, Verified: request.Verified})
	binding, err := capability.BindCoverageFeedback(coverage, feedback); if err != nil { fatal(err) }
	if err := binding.Validate(); err != nil { fatal(err) }
	if err := json.NewEncoder(os.Stdout).Encode(output{Options: options, Coverage: coverage, Feedback: feedback, Binding: binding}); err != nil { fatal(err) }
}

func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }