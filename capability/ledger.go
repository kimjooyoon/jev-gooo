package capability

import "fmt"

type InvestmentClass string

const (
	InvestmentNoInvestmentUntilEvidence InvestmentClass = "NO_INVESTMENT_UNTIL_EVIDENCE"
	InvestmentBoundedDeclarationEvidence InvestmentClass = "BOUNDED_DECLARATION_EVIDENCE"
	InvestmentBoundedVerification        InvestmentClass = "BOUNDED_VERIFICATION"
	InvestmentReviewBeforeCatalogChange  InvestmentClass = "REVIEW_BEFORE_CATALOG_CHANGE"
)

type ImprovementLedger struct {
	CoverageDigest         string                 `json:"coverage_digest"`
	BindingDigest          string                 `json:"binding_digest"`
	PlanDigest             string                 `json:"plan_digest"`
	SourceDigest           string                 `json:"source_digest"`
	EvidenceDigest         string                 `json:"evidence_digest,omitempty"`
	Scope                  string                 `json:"scope"`
	FeedbackStatus         CoverageFeedbackStatus `json:"feedback_status"`
	Decision               FocusDecision          `json:"decision"`
	NextAction             string                 `json:"next_action"`
	InvestmentClass        InvestmentClass        `json:"investment_class"`
	ObservedMatchedCount   int                    `json:"observed_matched_count"`
	ObservedAvailableCount int                    `json:"observed_available_count"`
	ObservedDeferredCount  int                    `json:"observed_deferred_count"`
	CompletenessClaimed    bool                   `json:"completeness_claimed"`
	LedgerDigest           string                 `json:"ledger_digest"`
}

func BuildImprovementLedger(coverage Coverage, binding CoverageFeedback, plan FocusPlan) (ImprovementLedger, error) {
	if err := coverage.Validate(); err != nil {
		return ImprovementLedger{}, err
	}
	if err := binding.Validate(); err != nil {
		return ImprovementLedger{}, err
	}
	if err := plan.Validate(); err != nil {
		return ImprovementLedger{}, err
	}
	if binding.CoverageDigest != coverage.CoverageDigest || plan.CoverageDigest != coverage.CoverageDigest || plan.BindingDigest != binding.BindingDigest {
		return ImprovementLedger{}, fmt.Errorf("improvement ledger inputs do not share evidence identity")
	}
	result := ImprovementLedger{
		CoverageDigest:         coverage.CoverageDigest,
		BindingDigest:          binding.BindingDigest,
		PlanDigest:             plan.PlanDigest,
		SourceDigest:           binding.SourceDigest,
		EvidenceDigest:         binding.EvidenceDigest,
		Scope:                  coverage.Scope,
		FeedbackStatus:         binding.Status,
		Decision:               plan.Decision,
		NextAction:             plan.NextAction,
		ObservedMatchedCount:   coverage.MatchedCount,
		ObservedAvailableCount: coverage.AvailableCount,
		ObservedDeferredCount:  coverage.DeferredCount,
		CompletenessClaimed:    false,
	}
	switch binding.Status {
	case CoverageFeedbackPreserveUnknown:
		result.InvestmentClass = InvestmentNoInvestmentUntilEvidence
	case CoverageFeedbackPreserveDeferred:
		result.InvestmentClass = InvestmentBoundedDeclarationEvidence
	case CoverageFeedbackPendingRequiredEvidence, CoverageFeedbackRequiresDeclaration:
		result.InvestmentClass = InvestmentBoundedVerification
	case CoverageFeedbackEligibleForCatalogReview:
		result.InvestmentClass = InvestmentReviewBeforeCatalogChange
	default:
		return ImprovementLedger{}, fmt.Errorf("unsupported improvement ledger status %q", binding.Status)
	}
	result.LedgerDigest = result.digest()
	return result, nil
}

func (ledger ImprovementLedger) Validate() error {
	for name, value := range map[string]string{
		"coverage": ledger.CoverageDigest,
		"binding":  ledger.BindingDigest,
		"plan":     ledger.PlanDigest,
		"source":   ledger.SourceDigest,
		"ledger":   ledger.LedgerDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("improvement ledger %s digest is invalid", name)
		}
	}
	if ledger.EvidenceDigest != "" && !validDigest(ledger.EvidenceDigest) {
		return fmt.Errorf("improvement ledger evidence digest is invalid")
	}
	if ledger.Scope != CoverageScope || ledger.CompletenessClaimed {
		return fmt.Errorf("improvement ledger overclaims its scope")
	}
	if ledger.ObservedMatchedCount < 0 || ledger.ObservedAvailableCount < 0 || ledger.ObservedDeferredCount < 0 || ledger.ObservedMatchedCount != ledger.ObservedAvailableCount+ledger.ObservedDeferredCount {
		return fmt.Errorf("improvement ledger observed counts are inconsistent")
	}
	switch ledger.FeedbackStatus {
	case CoverageFeedbackPreserveUnknown:
		if ledger.Decision != FocusDiscoveryRequired || ledger.NextAction != "clarify_or_extend_query" || ledger.InvestmentClass != InvestmentNoInvestmentUntilEvidence {
			return fmt.Errorf("unknown improvement ledger action is inconsistent")
		}
	case CoverageFeedbackPreserveDeferred:
		if ledger.Decision != FocusDeclarationEvidenceRequired || ledger.NextAction != "declare_missing_signals" || ledger.InvestmentClass != InvestmentBoundedDeclarationEvidence {
			return fmt.Errorf("deferred improvement ledger action is inconsistent")
		}
	case CoverageFeedbackPendingRequiredEvidence, CoverageFeedbackRequiresDeclaration:
		if ledger.Decision != FocusVerificationRequired || ledger.NextAction != "collect_required_evidence" || ledger.InvestmentClass != InvestmentBoundedVerification {
			return fmt.Errorf("pending improvement ledger action is inconsistent")
		}
	case CoverageFeedbackEligibleForCatalogReview:
		if ledger.Decision != FocusCatalogReviewReady || ledger.NextAction != "review_candidate_without_catalog_mutation" || ledger.InvestmentClass != InvestmentReviewBeforeCatalogChange {
			return fmt.Errorf("review improvement ledger action is inconsistent")
		}
	default:
		return fmt.Errorf("improvement ledger status %q is invalid", ledger.FeedbackStatus)
	}
	if ledger.digest() != ledger.LedgerDigest {
		return fmt.Errorf("improvement ledger digest does not match")
	}
	return nil
}

func (ledger ImprovementLedger) digest() string {
	return digest(
		"capability-improvement-ledger",
		ledger.CoverageDigest,
		ledger.BindingDigest,
		ledger.PlanDigest,
		ledger.SourceDigest,
		ledger.EvidenceDigest,
		ledger.Scope,
		string(ledger.FeedbackStatus),
		string(ledger.Decision),
		ledger.NextAction,
		string(ledger.InvestmentClass),
		fmt.Sprintf("%d", ledger.ObservedMatchedCount),
		fmt.Sprintf("%d", ledger.ObservedAvailableCount),
		fmt.Sprintf("%d", ledger.ObservedDeferredCount),
		fmt.Sprintf("%t", ledger.CompletenessClaimed),
	)
}
