package capability

import "fmt"

type FocusDecision string

const (
	FocusDiscoveryRequired          FocusDecision = "DISCOVERY_REQUIRED"
	FocusDeclarationEvidenceRequired FocusDecision = "DECLARATION_EVIDENCE_REQUIRED"
	FocusVerificationRequired       FocusDecision = "VERIFICATION_REQUIRED"
	FocusCatalogReviewReady         FocusDecision = "CATALOG_REVIEW_READY"
)

type FocusPlan struct {
	CoverageDigest       string        `json:"coverage_digest"`
	BindingDigest        string        `json:"binding_digest"`
	Scope                string        `json:"scope"`
	Decision             FocusDecision `json:"decision"`
	NextAction           string        `json:"next_action"`
	CompletenessClaimed  bool          `json:"completeness_claimed"`
	PlanDigest           string        `json:"plan_digest"`
}

func BuildFocusPlan(coverage Coverage, binding CoverageFeedback) (FocusPlan, error) {
	if err := coverage.Validate(); err != nil {
		return FocusPlan{}, err
	}
	if err := binding.Validate(); err != nil {
		return FocusPlan{}, err
	}
	if binding.CoverageDigest != coverage.CoverageDigest {
		return FocusPlan{}, fmt.Errorf("focus plan inputs do not share coverage identity")
	}
	plan := FocusPlan{
		CoverageDigest:      coverage.CoverageDigest,
		BindingDigest:       binding.BindingDigest,
		Scope:               coverage.Scope,
		CompletenessClaimed: false,
	}
	switch binding.Status {
	case CoverageFeedbackPreserveUnknown:
		plan.Decision = FocusDiscoveryRequired
		plan.NextAction = "clarify_or_extend_query"
	case CoverageFeedbackPreserveDeferred:
		plan.Decision = FocusDeclarationEvidenceRequired
		plan.NextAction = "declare_missing_signals"
	case CoverageFeedbackPendingRequiredEvidence, CoverageFeedbackRequiresDeclaration:
		plan.Decision = FocusVerificationRequired
		plan.NextAction = "collect_required_evidence"
	case CoverageFeedbackEligibleForCatalogReview:
		plan.Decision = FocusCatalogReviewReady
		plan.NextAction = "review_candidate_without_catalog_mutation"
	default:
		return FocusPlan{}, fmt.Errorf("unsupported coverage feedback status %q", binding.Status)
	}
	plan.PlanDigest = plan.digest()
	return plan, nil
}

func (p FocusPlan) Validate() error {
	if !validDigest(p.CoverageDigest) || !validDigest(p.BindingDigest) || !validDigest(p.PlanDigest) {
		return fmt.Errorf("focus plan digest is invalid")
	}
	if p.Scope != CoverageScope || p.CompletenessClaimed {
		return fmt.Errorf("focus plan overclaims its scope")
	}
	switch p.Decision {
	case FocusDiscoveryRequired:
		if p.NextAction != "clarify_or_extend_query" {
			return fmt.Errorf("focus discovery action is invalid")
		}
	case FocusDeclarationEvidenceRequired:
		if p.NextAction != "declare_missing_signals" {
			return fmt.Errorf("focus declaration action is invalid")
		}
	case FocusVerificationRequired:
		if p.NextAction != "collect_required_evidence" {
			return fmt.Errorf("focus verification action is invalid")
		}
	case FocusCatalogReviewReady:
		if p.NextAction != "review_candidate_without_catalog_mutation" {
			return fmt.Errorf("focus review action is invalid")
		}
	default:
		return fmt.Errorf("focus decision %q is invalid", p.Decision)
	}
	if p.digest() != p.PlanDigest {
		return fmt.Errorf("focus plan digest does not match")
	}
	return nil
}

func (p FocusPlan) digest() string {
	return digest(
		"capability-focus-plan",
		p.CoverageDigest,
		p.BindingDigest,
		p.Scope,
		string(p.Decision),
		p.NextAction,
		fmt.Sprintf("%t", p.CompletenessClaimed),
	)
}