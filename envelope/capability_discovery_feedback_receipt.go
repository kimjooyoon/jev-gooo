package envelope

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
)

// CapabilityDiscoveryFeedbackReceipt preserves user feedback about a discovery suggestion.
// It is an observation envelope and never authorizes or executes an operation.
type CapabilityDiscoveryFeedbackReceipt struct {
    Version                 string `json:"version"`
    SourceDigest            string `json:"source_digest"`
    DeclarationDigest       string `json:"declaration_digest"`
    CapabilityReceiptDigest string `json:"capability_receipt_digest"`
    EvidenceDigest          string `json:"evidence_digest"`
    OriginalQuery           string `json:"original_query"`
    SuggestedQuery          string `json:"suggested_query"`
    Outcome                 string `json:"outcome"`
    NonExecuting            bool   `json:"non_executing"`
    NonAuthorizing          bool   `json:"non_authorizing"`
}

// Validate enforces provenance and safety boundaries for feedback receipts.
func (r CapabilityDiscoveryFeedbackReceipt) Validate() error {
    if r.Version == "" || r.SourceDigest == "" || r.DeclarationDigest == "" {
        return fmt.Errorf("feedback receipt identity fields are required")
    }
    if r.CapabilityReceiptDigest == "" || r.EvidenceDigest == "" {
        return fmt.Errorf("feedback receipt evidence bindings are required")
    }
    if r.OriginalQuery == "" || r.SuggestedQuery == "" {
        return fmt.Errorf("feedback receipt queries are required")
    }
    switch r.Outcome {
    case "useful", "not_useful", "unresolved":
    default:
        return fmt.Errorf("unsupported feedback receipt outcome %q", r.Outcome)
    }
    if !r.NonExecuting || !r.NonAuthorizing {
        return fmt.Errorf("feedback receipt must be non-executing and non-authorizing")
    }
    return nil
}

// CanonicalDigest returns a stable digest for replay and domain measurement.
func (r CapabilityDiscoveryFeedbackReceipt) CanonicalDigest() (string, error) {
    if err := r.Validate(); err != nil {
        return "", err
    }
    encoded, err := json.Marshal(r)
    if err != nil {
        return "", fmt.Errorf("marshal feedback receipt: %w", err)
    }
    digest := sha256.Sum256(encoded)
    return hex.EncodeToString(digest[:]), nil
}