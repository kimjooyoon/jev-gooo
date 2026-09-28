package main

import "testing"

func TestBuildReceiptPreservesDeclarationBoundary(t *testing.T) {
	document := goooDiscoveryDocument{
		Trail: goooTrail{
			Response: goooResponse{
				Status:         "AVAILABLE",
				QueryDigest:    "sha256:" + "1" + "111111111111111111111111111111111111111111111111111111111111111",
				Capabilities:   []goooCapability{{ID: "ir_generation", NextOperation: "generate_ir"}},
				Declaration:    &goooDeclaration{Bound: true, SourceDigest: "sha256:" + "2" + "222222222222222222222222222222222222222222222222222222222222222"},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			EvidenceDigest: "sha256:" + "3" + "333333333333333333333333333333333333333333333333333333333333333",
		},
		Guide: goooGuide{NextOperations: []string{"generate_ir"}},
	}
	receipt, err := buildReceipt(document, "sha256:"+"4"+"444444444444444444444444444444444444444444444444444444444444444", "sha256:"+"5"+"555555555555555555555555555555555555555555555555555555555555555", "gooo-jev/capability-discovery/v1", "")
	if err != nil {
		t.Fatalf("buildReceipt() error = %v", err)
	}
	if receipt.DeclarationDigest != document.Trail.Response.Declaration.SourceDigest {
		t.Fatalf("declaration digest = %q, want bound digest", receipt.DeclarationDigest)
	}
}

func TestBuildReceiptRejectsUnboundDiscovery(t *testing.T) {
	_, err := buildReceipt(goooDiscoveryDocument{}, "sha256:"+"1"+"111111111111111111111111111111111111111111111111111111111111111", "sha256:"+"2"+"222222222222222222222222222222222222222222222222222222222222222", "toolchain", "inspect_next")
	if err == nil {
		t.Fatal("buildReceipt() accepted an unbound discovery")
	}
}