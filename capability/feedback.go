package capability

import (
	"fmt"
	"strings"
)

type FeedbackDisposition string

const (
	DispositionPendingRequiredEvidence   FeedbackDisposition = "PENDING_REQUIRED_EVIDENCE"
	DispositionPreserveUnknown            FeedbackDisposition = "PRESERVE_UNKNOWN"
	DispositionPreserveDeferred           FeedbackDisposition = "PRESERVE_DEFERRED"
	DispositionRequiresDeclarationBinding FeedbackDisposition = "REQUIRES_DECLARATION_BINDING"
	DispositionEligibleForCatalogReview   FeedbackDisposition = "ELIGIBLE_FOR_CATALOG_REVIEW"
)

type FeedbackInput struct {
	Source         string `json:"source"`
	EvidenceDigest string `json:"evidence_digest"`
	Verified       bool   `json:"verified"`
}

type Feedback struct {
	DiscoveryDigest string             `json:"discovery_digest"`
	DiscoveryStatus Status             `json:"discovery_status"`
	SourceDigest    string             `json:"source_digest"`
	EvidenceDigest  string             `json:"evidence_digest"`
	Verified        bool               `json:"verified"`
	Disposition     FeedbackDisposition `json:"disposition"`
	MissingStage    string             `json:"missing_stage"`
	FeedbackDigest  string             `json:"feedback_digest"`
}

func ObserveFeedback(discovery Discovery, input FeedbackInput) Feedback {
	result := Feedback{
		DiscoveryDigest: discovery.EvidenceDigest,
		DiscoveryStatus: discovery.Status,
		SourceDigest:    digest("capability-feedback-source", input.Source),
		EvidenceDigest:  input.EvidenceDigest,
		Verified:        input.Verified,
		Disposition:     DispositionPendingRequiredEvidence,
		MissingStage:    "verified_feedback_receipt",
	}

	switch {
	case discovery.Validate() != nil:
		result.Disposition = DispositionRequiresDeclarationBinding
		result.MissingStage = "capability_discovery"
	case discovery.Status == StatusUnknown:
		result.Disposition = DispositionPreserveUnknown
		result.MissingStage = "catalog_match"
	case discovery.Status == StatusDeferred:
		result.Disposition = DispositionPreserveDeferred
		result.MissingStage = discovery.MissingSignals[0]
	case strings.TrimSpace(input.Source) == "" || !validDigest(input.EvidenceDigest) || !input.Verified:
		result.Disposition = DispositionPendingRequiredEvidence
		result.MissingStage = "verified_feedback_receipt"
	default:
		result.Disposition = DispositionEligibleForCatalogReview
		result.MissingStage = ""
	}
	result.FeedbackDigest = result.digest()
	return result
}

func (f Feedback) Validate() error {
	if f.SourceDigest == "" || !validDigest(f.SourceDigest) {
		return fmt.Errorf("capability feedback source digest is invalid")
	}
	if f.EvidenceDigest != "" && !validDigest(f.EvidenceDigest) {
		return fmt.Errorf("capability feedback evidence digest is invalid")
	}
	if f.DiscoveryDigest != "" && !validDigest(f.DiscoveryDigest) {
		return fmt.Errorf("capability feedback discovery digest is invalid")
	}
	switch f.Disposition {
	case DispositionPendingRequiredEvidence:
		if strings.TrimSpace(f.MissingStage) == "" {
			return fmt.Errorf("pending capability feedback must preserve missing stage")
		}
	case DispositionPreserveUnknown:
		if f.DiscoveryStatus != StatusUnknown || f.MissingStage != "catalog_match" {
			return fmt.Errorf("unknown capability feedback changed discovery state")
		}
	case DispositionPreserveDeferred:
		if f.DiscoveryStatus != StatusDeferred || strings.TrimSpace(f.MissingStage) == "" {
			return fmt.Errorf("deferred capability feedback changed discovery state")
		}
	case DispositionRequiresDeclarationBinding:
		if strings.TrimSpace(f.MissingStage) == "" {
			return fmt.Errorf("unbound capability feedback must preserve missing stage")
		}
	case DispositionEligibleForCatalogReview:
		if f.DiscoveryStatus != StatusAvailable || !f.Verified || !validDigest(f.EvidenceDigest) || f.MissingStage != "" {
			return fmt.Errorf("catalog-review feedback is not evidence-bound")
		}
	default:
		return fmt.Errorf("capability feedback disposition %q is invalid", f.Disposition)
	}
	if f.FeedbackDigest == "" || !validDigest(f.FeedbackDigest) {
		return fmt.Errorf("capability feedback digest is invalid")
	}
	if f.digest() != f.FeedbackDigest {
		return fmt.Errorf("capability feedback digest does not match")
	}
	return nil
}

func (f Feedback) digest() string {
	return digest(
		"capability-feedback",
		f.DiscoveryDigest,
		string(f.DiscoveryStatus),
		f.SourceDigest,
		f.EvidenceDigest,
		fmt.Sprintf("%t", f.Verified),
		string(f.Disposition),
		f.MissingStage,
	)
}
