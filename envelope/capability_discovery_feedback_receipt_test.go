package envelope

import "testing"

func TestCapabilityDiscoveryFeedbackReceiptDigest(t *testing.T) {
    receipt := CapabilityDiscoveryFeedbackReceipt{
        Version:                 "capability.discovery.feedback.v1",
        SourceDigest:            "source-1",
        DeclarationDigest:       "declaration-1",
        CapabilityReceiptDigest: "capability-1",
        EvidenceDigest:          "evidence-1",
        OriginalQuery:           "What can gooo do with this declaration?",
        SuggestedQuery:          "Which constraints can be checked next?",
        Outcome:                 "useful",
        NonExecuting:            true,
        NonAuthorizing:          true,
    }

    first, err := receipt.CanonicalDigest()
    if err != nil {
        t.Fatalf("canonical digest: %v", err)
    }
    second, err := receipt.CanonicalDigest()
    if err != nil {
        t.Fatalf("canonical digest repeat: %v", err)
    }
    if first == "" || first != second {
        t.Fatalf("digest is not stable: %q %q", first, second)
    }
}

func TestCapabilityDiscoveryFeedbackReceiptRejectsAuthorization(t *testing.T) {
    receipt := CapabilityDiscoveryFeedbackReceipt{
        Version:                 "capability.discovery.feedback.v1",
        SourceDigest:            "source-1",
        DeclarationDigest:       "declaration-1",
        CapabilityReceiptDigest: "capability-1",
        EvidenceDigest:          "evidence-1",
        OriginalQuery:           "What can gooo do?",
        SuggestedQuery:          "What can be inspected?",
        Outcome:                 "unresolved",
        NonExecuting:            true,
        NonAuthorizing:          false,
    }

    if err := receipt.Validate(); err == nil {
        t.Fatal("expected authorizing receipt to be rejected")
    }
}