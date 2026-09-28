package envelope

import "testing"

func TestSelfImprovementIterationAcceptRequiresEvidence(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "transition-accept-001",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationProposed,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	accepted, err := iteration.Accept(selfImprovementDigest("c"), selfImprovementDigest("d"))
	if err != nil {
		t.Fatalf("proposal with evidence should be accepted: %v", err)
	}
	if accepted.State != SelfImprovementIterationAccepted {
		t.Fatalf("expected ACCEPTED, got %s", accepted.State)
	}
}

func TestSelfImprovementIterationRejectRequiresEvidence(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "transition-reject-001",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationProposed,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	rejected, err := iteration.Reject(selfImprovementDigest("c"), selfImprovementDigest("d"))
	if err != nil {
		t.Fatalf("proposal with evidence should be rejected: %v", err)
	}
	if rejected.State != SelfImprovementIterationRejected {
		t.Fatalf("expected REJECTED, got %s", rejected.State)
	}
}

func TestSelfImprovementIterationCannotAcceptWithoutEvidence(t *testing.T) {
	iteration := SelfImprovementIteration{
		Version:           SelfImprovementIterationVersion,
		IterationID:       "transition-invalid-001",
		ObservationDigest: selfImprovementDigest("a"),
		CandidateDigest:   selfImprovementDigest("b"),
		State:             SelfImprovementIterationProposed,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if _, err := iteration.Accept("", ""); err == nil {
		t.Fatal("acceptance without verification and rollback evidence must fail")
	}
}

