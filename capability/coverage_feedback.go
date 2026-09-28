package capability

import "fmt"

type CoverageFeedbackStatus string

const (
	CoverageFeedbackPreserveUnknown            CoverageFeedbackStatus = "PRESERVE_UNKNOWN"
	CoverageFeedbackPreserveDeferred           CoverageFeedbackStatus = "PRESERVE_DEFERRED"
	CoverageFeedbackPendingRequiredEvidence   CoverageFeedbackStatus = "PENDING_REQUIRED_EVIDENCE"
	CoverageFeedbackRequiresDeclaration       CoverageFeedbackStatus = "REQUIRES_DECLARATION_BINDING"
	CoverageFeedbackEligibleForCatalogReview   CoverageFeedbackStatus = "ELIGIBLE_FOR_CATALOG_REVIEW"
)

type CoverageFeedback struct {
	CoverageDigest          string                 `json:"coverage_digest"`
	FeedbackDigest          string                 `json:"feedback_digest"`
	SourceDigest            string                 `json:"source_digest"`
	FeedbackDisposition     FeedbackDisposition    `json:"feedback_disposition"`
	CoverageMatchedCount    int                    `json:"coverage_matched_count"`
	CoverageAvailableCount  int                    `json:"coverage_available_count"`
	CoverageDeferredCount   int                    `json:"coverage_deferred_count"`
	Status                  CoverageFeedbackStatus `json:"status"`
	MissingStage            string                 `json:"missing_stage"`
	EvidenceDigest          string                 `json:"evidence_digest"`
	BindingDigest           string                 `json:"binding_digest"`
}

func BindCoverageFeedback(coverage Coverage, feedback Feedback) (CoverageFeedback, error) {
	if err := coverage.Validate(); err != nil {
		return CoverageFeedback{}, err
	}
	if err := feedback.Validate(); err != nil {
		return CoverageFeedback{}, err
	}
	result := CoverageFeedback{
		CoverageDigest:         coverage.CoverageDigest,
		FeedbackDigest:         feedback.FeedbackDigest,
		SourceDigest:           feedback.SourceDigest,
		FeedbackDisposition:    feedback.Disposition,
		CoverageMatchedCount:   coverage.MatchedCount,
		CoverageAvailableCount: coverage.AvailableCount,
		CoverageDeferredCount:  coverage.DeferredCount,
		Status:                 CoverageFeedbackPendingRequiredEvidence,
		MissingStage:           "coverage-feedback-binding",
		EvidenceDigest:         feedback.EvidenceDigest,
	}
	switch feedback.Disposition {
	case DispositionPreserveUnknown:
		result.Status = CoverageFeedbackPreserveUnknown
		result.MissingStage = feedback.MissingStage
	case DispositionPreserveDeferred:
		result.Status = CoverageFeedbackPreserveDeferred
		result.MissingStage = feedback.MissingStage
	case DispositionPendingRequiredEvidence:
		result.Status = CoverageFeedbackPendingRequiredEvidence
		result.MissingStage = feedback.MissingStage
	case DispositionRequiresDeclarationBinding:
		result.Status = CoverageFeedbackRequiresDeclaration
		result.MissingStage = feedback.MissingStage
	case DispositionEligibleForCatalogReview:
		switch {
		case coverage.UnknownQuery || coverage.MatchedCount == 0:
			result.Status = CoverageFeedbackPreserveUnknown
			result.MissingStage = "coverage_match"
		case coverage.AvailableCount == 0:
			result.Status = CoverageFeedbackPreserveDeferred
			result.MissingStage = "available_capability"
		default:
			result.Status = CoverageFeedbackEligibleForCatalogReview
			result.MissingStage = ""
		}
	default:
		return CoverageFeedback{}, fmt.Errorf("unsupported capability feedback disposition %q", feedback.Disposition)
	}
	result.BindingDigest = result.digest()
	return result, nil
}

func (b CoverageFeedback) Validate() error {
	for name, value := range map[string]string{
		"coverage": b.CoverageDigest,
		"feedback": b.FeedbackDigest,
		"source": b.SourceDigest,
		"binding": b.BindingDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("coverage feedback %s digest is invalid", name)
		}
	}
	if b.EvidenceDigest != "" && !validDigest(b.EvidenceDigest) {
		return fmt.Errorf("coverage feedback evidence digest is invalid")
	}
	if b.CoverageMatchedCount < 0 || b.CoverageAvailableCount < 0 || b.CoverageDeferredCount < 0 || b.CoverageMatchedCount != b.CoverageAvailableCount+b.CoverageDeferredCount {
		return fmt.Errorf("coverage feedback counts are inconsistent")
	}
	switch b.Status {
	case CoverageFeedbackPreserveUnknown, CoverageFeedbackPreserveDeferred, CoverageFeedbackPendingRequiredEvidence, CoverageFeedbackRequiresDeclaration:
		if b.MissingStage == "" {
			return fmt.Errorf("preserved coverage feedback must retain missing stage")
		}
	case CoverageFeedbackEligibleForCatalogReview:
		if b.FeedbackDisposition != DispositionEligibleForCatalogReview || b.CoverageMatchedCount == 0 || b.CoverageAvailableCount == 0 || b.MissingStage != "" {
			return fmt.Errorf("eligible coverage feedback is not evidence-bound")
		}
	default:
		return fmt.Errorf("coverage feedback status %q is invalid", b.Status)
	}
	if b.digest() != b.BindingDigest {
		return fmt.Errorf("coverage feedback binding digest does not match")
	}
	return nil
}

func (b CoverageFeedback) digest() string {
	return digest(
		"capability-coverage-feedback",
		b.CoverageDigest,
		b.FeedbackDigest,
		b.SourceDigest,
		string(b.FeedbackDisposition),
		fmt.Sprintf("%d", b.CoverageMatchedCount),
		fmt.Sprintf("%d", b.CoverageAvailableCount),
		fmt.Sprintf("%d", b.CoverageDeferredCount),
		string(b.Status),
		b.MissingStage,
		b.EvidenceDigest,
	)
}