package envelope

import "testing"

func testRequest(t *testing.T) (CapabilityRequest, Declaration) {
	t.Helper()
	declaration, err := BindDeclaration("package billing\nentity Invoice\n")
	if err != nil {
		t.Fatalf("bind declaration: %v", err)
	}
	request, err := NewCapabilityRequest(declaration, "spiffe://example.test/workload", "billing", "invoice.read")
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	return request, declaration
}

func TestDeclarationBindsRequest(t *testing.T) {
	request, declaration := testRequest(t)
	if request.DeclarationDigest != declaration.Digest {
		t.Fatalf("request declaration digest = %q, want %q", request.DeclarationDigest, declaration.Digest)
	}
	declaration.Source = "tampered"
	if err := declaration.Validate(); err == nil {
		t.Fatal("tampered declaration was accepted")
	}
}

func TestExecutionReceiptPreservesFirstMissingStage(t *testing.T) {
	request, _ := testRequest(t)
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

func TestReverseObservationPreservesNestedUnknown(t *testing.T) {
	request, _ := testRequest(t)
	receipt, err := NewExecutionReceipt(request, "", WorkloadIdentity{SPIFFEID: "spiffe://example.test/workload"}, DecisionReceipt{Digest: digest("decision"), NonAuthorizing: true}, digest("result"), true)
	if err != nil {
		t.Fatalf("new receipt: %v", err)
	}
	observation := ObserveReverse(receipt, digest("observed"), digest("verifier"))
	if observation.Status != StatusUnknown || observation.MissingStage != "execution_receipt:capability_grant" {
		t.Fatalf("observation = %#v, want nested unknown", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("validate observation: %v", err)
	}
}

func TestReverseObservationRequiresTerminalEvidence(t *testing.T) {
	request, _ := testRequest(t)
	receipt, err := NewExecutionReceipt(request, digest("grant"), WorkloadIdentity{SPIFFEID: "spiffe://example.test/workload"}, DecisionReceipt{Digest: digest("decision"), NonAuthorizing: true}, digest("result"), true)
	if err != nil {
		t.Fatalf("new receipt: %v", err)
	}
	observation := ObserveReverse(receipt, digest("observed"), digest("verifier"))
	if observation.Status != StatusObserved || observation.MissingStage != "" {
		t.Fatalf("observation = %#v, want observed", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("validate observation: %v", err)
	}
	observation.VerifierDigest = digest("tampered")
	if err := observation.Validate(); err == nil {
		t.Fatal("tampered reverse observation was accepted")
	}
}
