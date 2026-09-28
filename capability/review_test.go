package capability

import (
	"strings"
	"testing"
)

func reviewFeedback(t *testing.T) Feedback {
	t.Helper()
	discovery, err := DiscoverFromSource(
		"What can this language do?",
		"package jev\nactivity discover_capability\nproperty evidence_digest string",
	)
	if err != nil {
		t.Fatal(err)
	}
	feedback := ObserveFeedback(discovery, FeedbackInput{
		Source:         "external evaluator receipt",
		EvidenceDigest: strings.Repeat("b", 64),
		Verified:       true,
	})
	if err := feedback.Validate(); err != nil {
		t.Fatal(err)
	}
	return feedback
}

func TestProposeReviewRequiresExplicitProposal(t *testing.T) {
	candidate, err := ProposeReview(reviewFeedback(t), "new-capability", "")
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Status != ReviewPendingProposal {
		t.Fatalf("status = %s, want %s", candidate.Status, ReviewPendingProposal)
	}
	if err := candidate.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProposeReviewDoesNotPromoteDeferred(t *testing.T) {
	discovery, err := DiscoverFromSource("Can this language run it?", "package jev\nactivity inspect")
	if err != nil {
		t.Fatal(err)
	}
	feedback := ObserveFeedback(discovery, FeedbackInput{Source: "note"})
	candidate, err := ProposeReview(feedback, "execution", "run support")
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Status != ReviewPreserved {
		t.Fatalf("status = %s, want %s", candidate.Status, ReviewPreserved)
	}
	if err := candidate.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProposeReviewProducesInspectableProposal(t *testing.T) {
	candidate, err := ProposeReview(reviewFeedback(t), "new-capability", "supports a declared boundary")
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Status != ReviewProposed {
		t.Fatalf("status = %s, want %s", candidate.Status, ReviewProposed)
	}
	if err := candidate.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProposeReviewRejectsTampering(t *testing.T) {
	candidate, err := ProposeReview(reviewFeedback(t), "new-capability", "supports a declared boundary")
	if err != nil {
		t.Fatal(err)
	}
	candidate.CandidateDigest = strings.Repeat("0", 64)
	if err := candidate.Validate(); err == nil {
		t.Fatal("tampered review candidate should be rejected")
	}
}
