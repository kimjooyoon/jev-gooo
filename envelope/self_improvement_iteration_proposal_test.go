package envelope

import "testing"

func TestSelfImprovementIterationProposeNextFromUnknown(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "proposal-unknown-001",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationUnknown,
		FirstMissingStage: "verification_boundary",
		NextOperation:     "provide_verification_evidence",
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	next, err := iteration.ProposeNext("proposal-next-001", selfImprovementDigest("c"))
	if err != nil {
		t.Fatalf("unknown iteration should accept a bounded next proposal: %v", err)
	}
	if next.State != SelfImprovementIterationProposed {
		t.Fatalf("expected PROPOSED, got %s", next.State)
	}
	if next.IterationID != "proposal-next-001" {
		t.Fatalf("unexpected next iteration id %q", next.IterationID)
	}
}

func TestSelfImprovementIterationProposeNextRequiresCandidateDigest(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "proposal-unknown-002",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationUnknown,
		FirstMissingStage: "verification_boundary",
		NextOperation:     "provide_verification_evidence",
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if _, err := iteration.ProposeNext("proposal-next-002", ""); err == nil {
		t.Fatal("next proposal without a candidate digest must fail")
	}
}

