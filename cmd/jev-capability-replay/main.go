package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type report struct {
	Valid                 bool                                  `json:"valid"`
	Outcome               envelope.CapabilityDiscoveryReplayOutcome `json:"outcome"`
	PriorEvidenceDigest   string                                `json:"prior_evidence_digest"`
	CurrentEvidenceDigest string                                `json:"current_evidence_digest"`
	FirstMissingStage     string                                `json:"first_missing_stage,omitempty"`
	NextOperation         string                                `json:"next_operation"`
	Error                 string                                `json:"error,omitempty"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: jev-capability-replay <receipt.json> <replay.json>")
		os.Exit(2)
	}
	receiptData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fail(err)
	}
	replayData, err := os.ReadFile(os.Args[2])
	if err != nil {
		fail(err)
	}
	var receipt envelope.CapabilityDiscoveryReceipt
	if err := json.Unmarshal(receiptData, &receipt); err != nil {
		fail(err)
	}
	var replay envelope.CapabilityDiscoveryReplay
	if err := json.Unmarshal(replayData, &replay); err != nil {
		fail(err)
	}
	result := report{
		Valid:                 true,
		Outcome:               replay.Outcome,
		PriorEvidenceDigest:   replay.PriorEvidenceDigest,
		CurrentEvidenceDigest: replay.CurrentEvidenceDigest,
		FirstMissingStage:     replay.FirstMissingStage,
		NextOperation:         replay.NextOperation,
	}
	if err := replay.ValidateAgainst(receipt); err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fail(err)
	}
	if !result.Valid {
		os.Exit(2)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}