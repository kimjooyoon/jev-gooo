package envelope

import (
	"strings"
	"testing"
)

func validCapabilityDiscoveryReplay() CapabilityDiscoveryReplay {
	receipt := validCapabilityDiscoveryReceipt()
	return CapabilityDiscoveryReplay{
		Version:               CapabilityDiscoveryReplayVersion,
		PriorEvidenceDigest:   receipt.EvidenceDigest,
		ObservationDigest:     "sha256:" + strings.Repeat("6", 64),
		CurrentEvidenceDigest: "sha256:" + strings.Repeat("7", 64),
		Outcome:               CapabilityDiscoveryReplayConfirmed,
		NextOperation:         receipt.NextOperation,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
}

func TestCapabilityDiscoveryReplayContinuesPriorEvidence(t *testing.T) {
	receipt := validCapabilityDiscoveryReceipt()
	if err := validCapabilityDiscoveryReplay().ValidateAgainst(receipt); err != nil {
		t.Fatalf("ValidateAgainst() error = %v", err)
	}
}

func TestCapabilityDiscoveryReplayPreservesUnresolvedBoundary(t *testing.T) {
	replay := validCapabilityDiscoveryReplay()
	replay.Outcome = CapabilityDiscoveryReplayUnresolved
	replay.FirstMissingStage = "external_boundary"
	if err := replay.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestCapabilityDiscoveryReplayRejectsPriorEvidenceMismatch(t *testing.T) {
	receipt := validCapabilityDiscoveryReceipt()
	replay := validCapabilityDiscoveryReplay()
	replay.PriorEvidenceDigest = receipt.QueryDigest
	if err := replay.ValidateAgainst(receipt); err == nil {
		t.Fatal("ValidateAgainst() accepted a mismatched prior evidence digest")
	}
}