package envelope

import "testing"

func TestExecutionReceiptPreservesFirstMissingStage(t *testing.T) {
	request := CapabilityRequest{
		Subject:           "spiffe://example.test/workload",
		Audience:          "billing",
		Capability:        "invoice.read",
		DeclarationDigest: digest("declaration"),
	}
	receipt, err := NewExecutionReceipt(request, "", WorkloadIdentity{SPIFFEID: "spiffe://example.test/workload"}, DecisionReceipt{Digest: digest("decision"), NonAuthorizing: true}, digest("result"), true)
	if err != nil {
		t.Fatalf("new receipt: %v", err)
	}
	if receipt.Status != StatusUnknown || receipt.MissingStage != "capability_grant" {
		t.Fatalf("receipt = %#v, want unknown at capability grant", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("validate receipt: %v", err)
	}
}

func TestExecutionReceiptCompletesOnlyWithTerminalEvidence(t *testing.T) {
	request := CapabilityRequest{
		Subject:           "spiffe://example.test/workload",
		Audience:          "billing",
		Capability:        "invoice.read",
		DeclarationDigest: digest("declaration"),
	}
	receipt, err := NewExecutionReceipt(request, digest("grant"), WorkloadIdentity{SPIFFEID: "spiffe://example.test/workload"}, DecisionReceipt{Digest: digest("decision"), NonAuthorizing: true}, digest("result"), true)
	if err != nil {
		t.Fatalf("new receipt: %v", err)
	}
	if receipt.Status != StatusCompleted || receipt.MissingStage != "" {
		t.Fatalf("receipt = %#v, want completed", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("validate receipt: %v", err)
	}
	receipt.ResultDigest = digest("tampered")
	if err := receipt.Validate(); err == nil {
		t.Fatal("tampered receipt was accepted")
	}
}
