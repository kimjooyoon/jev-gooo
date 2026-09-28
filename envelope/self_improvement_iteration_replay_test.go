package envelope

import (
	"strings"
	"testing"
)

func TestSelfImprovementIterationReplayPreservesContinuity(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "iteration-replay-001",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationProposed,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	priorDigest, err := iteration.CanonicalDigest()
	if err != nil {
		t.Fatalf("iteration should have a canonical digest: %v", err)
	}
	replay := SelfImprovementIterationReplay{
		Version:              SelfImprovementIterationReplayVersion,
		PriorIterationDigest: priorDigest,
		ObservationDigest:    "sha256:" + strings.Repeat("c", 64),
		Outcome:              SelfImprovementIterationReplayUnresolved,
		FirstMissingStage:    "verification_boundary",
		NextOperation:        "provide_verification_evidence",
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := replay.ValidateAgainst(iteration); err != nil {
		t.Fatalf("replay should preserve continuity: %v", err)
	}
}

func TestSelfImprovementIterationReplayRejectsMismatchedPrior(t *testing.T) {
	replay := SelfImprovementIterationReplay{
		Version:              SelfImprovementIterationReplayVersion,
		PriorIterationDigest: "sha256:" + strings.Repeat("a", 64),
		ObservationDigest:    "sha256:" + strings.Repeat("b", 64),
		Outcome:              SelfImprovementIterationReplayUnresolved,
		FirstMissingStage:    "verification_boundary",
		NextOperation:        "provide_verification_evidence",
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "iteration-replay-002",
		ObservationDigest: selfImprovementDigest("c"),
		CandidateDigest:   selfImprovementDigest("d"),
		State:             SelfImprovementIterationProposed,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if err := replay.ValidateAgainst(iteration); err == nil {
		t.Fatal("replay with a mismatched prior digest must be rejected")
	}
}

