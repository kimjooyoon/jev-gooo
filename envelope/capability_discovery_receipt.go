package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const CapabilityDiscoveryReceiptVersion = "capability.discovery.v1"

type CapabilityDiscoveryReceiptState string

const (
	CapabilityDiscoveryAvailable CapabilityDiscoveryReceiptState = "AVAILABLE"
	CapabilityDiscoveryDeferred  CapabilityDiscoveryReceiptState = "DEFERRED"
	CapabilityDiscoveryUnknown   CapabilityDiscoveryReceiptState = "UNKNOWN"
)

// CapabilityDiscoveryReceipt transports a bounded .gooo capability result.
// It is descriptive only: it cannot execute code or authorize access.
type CapabilityDiscoveryReceipt struct {
	Version           string                        `json:"schema_version"`
	Status            CapabilityDiscoveryReceiptState `json:"status"`
	SourceDigest      string                        `json:"source_digest"`
	DeclarationDigest string                        `json:"declaration_digest,omitempty"`
	QueryDigest       string                        `json:"query_digest"`
	CatalogDigest     string                        `json:"catalog_digest"`
	EvidenceDigest    string                        `json:"evidence_digest"`
	ToolchainIdentity string                        `json:"toolchain_identity"`
	CapabilityIDs     []string                      `json:"capability_ids,omitempty"`
	FirstMissingStage string                        `json:"first_missing_stage,omitempty"`
	NextOperation     string                        `json:"next_operation"`
	NonExecuting      bool                          `json:"non_executing"`
	NonAuthorizing    bool                          `json:"non_authorizing"`
}

func (receipt CapabilityDiscoveryReceipt) Validate() error {
	if receipt.Version != CapabilityDiscoveryReceiptVersion {
		return fmt.Errorf("unsupported capability discovery receipt version %q", receipt.Version)
	}
	switch receipt.Status {
	case CapabilityDiscoveryAvailable, CapabilityDiscoveryDeferred, CapabilityDiscoveryUnknown:
	default:
		return fmt.Errorf("unsupported capability discovery status %q", receipt.Status)
	}
	for name, value := range map[string]string{
		"source_digest": receipt.SourceDigest,
		"query_digest": receipt.QueryDigest,
		"catalog_digest": receipt.CatalogDigest,
		"evidence_digest": receipt.EvidenceDigest,
	} {
		if err := validateCapabilityDiscoveryDigest(name, value); err != nil {
			return err
		}
	}
	if receipt.Status == CapabilityDiscoveryAvailable {
		if err := validateCapabilityDiscoveryDigest("declaration_digest", receipt.DeclarationDigest); err != nil {
			return err
		}
		if len(receipt.CapabilityIDs) == 0 || receipt.FirstMissingStage != "" {
			return fmt.Errorf("available capability discovery receipt has incomplete result")
		}
	} else if receipt.FirstMissingStage == "" {
		return fmt.Errorf("non-available capability discovery receipt lost first missing stage")
	}
	if receipt.ToolchainIdentity == "" || receipt.NextOperation == "" {
		return fmt.Errorf("capability discovery receipt is missing toolchain or next operation")
	}
	if !receipt.NonExecuting || !receipt.NonAuthorizing {
		return fmt.Errorf("capability discovery receipt crossed a capability boundary")
	}
	seen := make(map[string]bool, len(receipt.CapabilityIDs))
	for _, id := range receipt.CapabilityIDs {
		if strings.TrimSpace(id) == "" || seen[id] {
			return fmt.Errorf("capability discovery receipt has duplicate or empty capability id")
		}
		seen[id] = true
	}
	return nil
}

func (receipt CapabilityDiscoveryReceipt) CanonicalDigest() (string, error) {
	if err := receipt.Validate(); err != nil {
		return "", err
	}
	canonical := receipt
	canonical.CapabilityIDs = append([]string(nil), receipt.CapabilityIDs...)
	sort.Strings(canonical.CapabilityIDs)
	payload, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func validateCapabilityDiscoveryDigest(name, value string) error {
	if len(value) != len("sha256:")+64 || !strings.HasPrefix(value, "sha256:") {
		return fmt.Errorf("invalid %s", name)
	}
	for _, character := range value[len("sha256:"):] {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return fmt.Errorf("invalid %s", name)
		}
	}
	return nil
}
