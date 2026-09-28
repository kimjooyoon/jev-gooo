package envelope

import "testing"

func validCapabilityDiscoveryReceipt() CapabilityDiscoveryReceipt {
	return CapabilityDiscoveryReceipt{
		Version: CapabilityDiscoveryReceiptVersion,
		Status: CapabilityDiscoveryAvailable,
		SourceDigest: "sha256:" + "1" + "111111111111111111111111111111111111111111111111111111111111111",
		DeclarationDigest: "sha256:" + "2" + "222222222222222222222222222222222222222222222222222222222222222",
		QueryDigest: "sha256:" + "3" + "333333333333333333333333333333333333333333333333333333333333333",
		CatalogDigest: "sha256:" + "4" + "444444444444444444444444444444444444444444444444444444444444444",
		EvidenceDigest: "sha256:" + "5" + "555555555555555555555555555555555555555555555555555555555555555",
		ToolchainIdentity: "gooo-jev/capability-discovery/v1",
		CapabilityIDs: []string{"canonical_generation", "syntax_completion"},
		NextOperation: "inspect_next_non_executing_operation",
		NonExecuting: true,
		NonAuthorizing: true,
	}
}

func TestCapabilityDiscoveryReceiptValidatesAvailableResult(t *testing.T) {
	receipt := validCapabilityDiscoveryReceipt()
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if _, err := receipt.CanonicalDigest(); err != nil {
		t.Fatalf("CanonicalDigest() error = %v", err)
	}
}

func TestCapabilityDiscoveryReceiptPreservesUnknownBoundary(t *testing.T) {
	receipt := validCapabilityDiscoveryReceipt()
	receipt.Status = CapabilityDiscoveryUnknown
	receipt.DeclarationDigest = ""
	receipt.CapabilityIDs = nil
	receipt.FirstMissingStage = "evidence_digest"
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestCapabilityDiscoveryReceiptRejectsBoundaryCrossing(t *testing.T) {
	receipt := validCapabilityDiscoveryReceipt()
	receipt.NonExecuting = false
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() accepted an executing receipt")
	}
}
