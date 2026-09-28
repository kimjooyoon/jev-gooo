package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type goooCapability struct {
	ID            string `json:"id"`
	NextOperation string `json:"next_operation"`
}

type goooDeclaration struct {
	Bound        bool   `json:"bound"`
	SourceDigest string `json:"source_digest"`
}

type goooResponse struct {
	Status          string             `json:"status"`
	QueryDigest     string             `json:"query_digest"`
	Capabilities    []goooCapability   `json:"capabilities"`
	Declaration     *goooDeclaration   `json:"declaration,omitempty"`
	MissingStage    string             `json:"missing_stage"`
	NonExecuting    bool               `json:"non_executing"`
	NonAuthorizing  bool               `json:"non_authorizing"`
}

type goooTrail struct {
	Response       goooResponse `json:"response"`
	EvidenceDigest string       `json:"evidence_digest"`
}

type goooGuide struct {
	NextOperations []string `json:"next_operations"`
}

type goooDiscoveryDocument struct {
	Trail goooTrail `json:"trail"`
	Guide goooGuide `json:"guide"`
}

func main() {
	inputPath := flag.String("input", "-", "gooo capability discovery JSON path, or - for stdin")
	sourceDigest := flag.String("source-digest", "", "digest of the source artifact that produced the discovery")
	catalogDigest := flag.String("catalog-digest", "", "digest of the capability catalog")
	toolchainIdentity := flag.String("toolchain-identity", "", "exact gooo/JEV toolchain identity")
	nextOperation := flag.String("next-operation", "", "explicit next non-executing operation when the discovery does not provide one")
	flag.Parse()

	payload, err := readInput(*inputPath)
	if err != nil {
		fail(err)
	}
	var document goooDiscoveryDocument
	if err := json.Unmarshal(payload, &document); err != nil {
		fail(fmt.Errorf("decode gooo capability discovery: %w", err))
	}
	receipt, err := buildReceipt(document, *sourceDigest, *catalogDigest, *toolchainIdentity, *nextOperation)
	if err != nil {
		fail(err)
	}
	encoded, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		fail(err)
	}
	fmt.Println(string(encoded))
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filepath.Clean(path))
}

func buildReceipt(document goooDiscoveryDocument, sourceDigest, catalogDigest, toolchainIdentity, explicitNextOperation string) (envelope.CapabilityDiscoveryReceipt, error) {
	response := document.Trail.Response
	if response.Declaration == nil || !response.Declaration.Bound || strings.TrimSpace(response.Declaration.SourceDigest) == "" {
		return envelope.CapabilityDiscoveryReceipt{}, fmt.Errorf("gooo discovery is not bound to a declaration")
	}
	nextOperation := strings.TrimSpace(explicitNextOperation)
	if nextOperation == "" && len(document.Guide.NextOperations) > 0 {
		nextOperation = strings.TrimSpace(document.Guide.NextOperations[0])
	}
	if nextOperation == "" {
		for _, capability := range response.Capabilities {
			if strings.TrimSpace(capability.NextOperation) != "" {
				nextOperation = strings.TrimSpace(capability.NextOperation)
				break
			}
		}
	}
	if nextOperation == "" {
		return envelope.CapabilityDiscoveryReceipt{}, fmt.Errorf("gooo discovery has no next operation; provide -next-operation explicitly")
	}

	capabilityIDs := make([]string, 0, len(response.Capabilities))
	for _, capability := range response.Capabilities {
		capabilityIDs = append(capabilityIDs, capability.ID)
	}
	receipt := envelope.CapabilityDiscoveryReceipt{
		Version:           envelope.CapabilityDiscoveryReceiptVersion,
		Status:            envelope.CapabilityDiscoveryReceiptState(response.Status),
		SourceDigest:      sourceDigest,
		DeclarationDigest: response.Declaration.SourceDigest,
		QueryDigest:       response.QueryDigest,
		CatalogDigest:     catalogDigest,
		EvidenceDigest:    document.Trail.EvidenceDigest,
		ToolchainIdentity: toolchainIdentity,
		CapabilityIDs:     capabilityIDs,
		FirstMissingStage: response.MissingStage,
		NextOperation:     nextOperation,
		NonExecuting:      response.NonExecuting,
		NonAuthorizing:    response.NonAuthorizing,
	}
	if err := receipt.Validate(); err != nil {
		return envelope.CapabilityDiscoveryReceipt{}, fmt.Errorf("validate converted capability receipt: %w", err)
	}
	return receipt, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}