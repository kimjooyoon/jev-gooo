package envelope

import (
	"fmt"
	"strings"
)

type CapabilityObservationInvestment struct {
	ObservationDigest  string                     `json:"observation_digest"`
	State              CapabilityObservationState `json:"state"`
	Decision           string                     `json:"decision"`
	NextAction         string                     `json:"next_action"`
	InvestmentClass    string                     `json:"investment_class"`
	CompletenessClaimed bool                      `json:"completeness_claimed"`
	PlanDigest         string                     `json:"plan_digest"`
}

func (r CapabilityObservationReport) BuildInvestment() (CapabilityObservationInvestment, error) {
	if err := r.Validate(); err != nil {
		return CapabilityObservationInvestment{}, err
	}
	investment := CapabilityObservationInvestment{
		ObservationDigest:  r.ObservationDigest,
		State:              r.State,
		CompletenessClaimed: false,
	}
	switch r.State {
	case CapabilityObservationUnknown:
		investment.Decision = "DISCOVERY_REQUIRED"
		investment.NextAction = "clarify_or_extend_query"
		investment.InvestmentClass = "NO_INVESTMENT_UNTIL_EVIDENCE"
	case CapabilityObservationDeferred:
		investment.Decision = "DECLARATION_EVIDENCE_REQUIRED"
		investment.NextAction = "declare_missing_signals"
		investment.InvestmentClass = "BOUNDED_DECLARATION_EVIDENCE"
	case CapabilityObservationAvailable:
		investment.Decision = "CATALOG_REVIEW_READY"
		investment.NextAction = "review_candidate_without_catalog_mutation"
		investment.InvestmentClass = "REVIEW_BEFORE_CATALOG_CHANGE"
	default:
		return CapabilityObservationInvestment{}, fmt.Errorf("unsupported observation state %q", r.State)
	}
	investment.PlanDigest = digest(
		"capability-observation-investment",
		investment.ObservationDigest,
		string(investment.State),
		investment.Decision,
		investment.NextAction,
		investment.InvestmentClass,
		fmt.Sprintf("%t", investment.CompletenessClaimed),
	)
	return investment, investment.Validate()
}

func (i CapabilityObservationInvestment) Validate() error {
	if !validDigest(i.ObservationDigest) || !validDigest(i.PlanDigest) {
		return fmt.Errorf("capability observation investment digest is invalid")
	}
	if i.CompletenessClaimed {
		return fmt.Errorf("capability observation investment overclaims completeness")
	}
	if strings.TrimSpace(i.Decision) == "" || strings.TrimSpace(i.NextAction) == "" || strings.TrimSpace(i.InvestmentClass) == "" {
		return fmt.Errorf("capability observation investment is incomplete")
	}
	switch i.State {
	case CapabilityObservationUnknown:
		if i.Decision != "DISCOVERY_REQUIRED" || i.NextAction != "clarify_or_extend_query" || i.InvestmentClass != "NO_INVESTMENT_UNTIL_EVIDENCE" {
			return fmt.Errorf("unknown capability observation investment is inconsistent")
		}
	case CapabilityObservationDeferred:
		if i.Decision != "DECLARATION_EVIDENCE_REQUIRED" || i.NextAction != "declare_missing_signals" || i.InvestmentClass != "BOUNDED_DECLARATION_EVIDENCE" {
			return fmt.Errorf("deferred capability observation investment is inconsistent")
		}
	case CapabilityObservationAvailable:
		if i.Decision != "CATALOG_REVIEW_READY" || i.NextAction != "review_candidate_without_catalog_mutation" || i.InvestmentClass != "REVIEW_BEFORE_CATALOG_CHANGE" {
			return fmt.Errorf("available capability observation investment is inconsistent")
		}
	default:
		return fmt.Errorf("capability observation investment state %q is invalid", i.State)
	}
	expected := digest(
		"capability-observation-investment",
		i.ObservationDigest,
		string(i.State),
		i.Decision,
		i.NextAction,
		i.InvestmentClass,
		fmt.Sprintf("%t", i.CompletenessClaimed),
	)
	if expected != i.PlanDigest {
		return fmt.Errorf("capability observation investment digest does not match")
	}
	return nil
}
