package capability

import (
	"fmt"
	"strings"
)

type ReviewStatus string

const (
	ReviewPreserved       ReviewStatus = "PRESERVED"
	ReviewPendingProposal ReviewStatus = "PENDING_PROPOSAL"
	ReviewProposed        ReviewStatus = "PROPOSED"
)

type ReviewCandidate struct {
	FeedbackDigest   string       `json:"feedback_digest"`
	CapabilityID     string       `json:"capability_id"`
	RationaleDigest   string       `json:"rationale_digest"`
	RationalePresent  bool         `json:"rationale_present"`
	Status           ReviewStatus `json:"status"`
	MissingStage     string       `json:"missing_stage"`
	CandidateDigest  string       `json:"candidate_digest"`
}

func ProposeReview(feedback Feedback, capabilityID, rationale string) (ReviewCandidate, error) {
	if err := feedback.Validate(); err != nil {
		return ReviewCandidate{}, err
	}
	candidate := ReviewCandidate{
		FeedbackDigest:  feedback.FeedbackDigest,
		CapabilityID:    strings.TrimSpace(capabilityID),
		RationaleDigest: digest("review-rationale", rationale),
		RationalePresent: strings.TrimSpace(rationale) != "",
		Status:          ReviewPreserved,
		MissingStage:    feedback.MissingStage,
	}
	switch {
	case feedback.Disposition != DispositionEligibleForCatalogReview:
		candidate.Status = ReviewPreserved
	case candidate.CapabilityID == "" || !candidate.RationalePresent:
		candidate.Status = ReviewPendingProposal
		candidate.MissingStage = "review_proposal"
	default:
		candidate.Status = ReviewProposed
		candidate.MissingStage = ""
	}
	candidate.CandidateDigest = candidate.digest()
	return candidate, nil
}

func (c ReviewCandidate) Validate() error {
	if !validDigest(c.FeedbackDigest) || !validDigest(c.RationaleDigest) || !validDigest(c.CandidateDigest) {
		return fmt.Errorf("review candidate digest is invalid")
	}
	switch c.Status {
	case ReviewPreserved:
		if strings.TrimSpace(c.MissingStage) == "" {
			return fmt.Errorf("preserved review candidate must retain missing stage")
		}
	case ReviewPendingProposal:
		if c.MissingStage != "review_proposal" {
			return fmt.Errorf("pending review candidate must retain proposal stage")
		}
	case ReviewProposed:
		if strings.TrimSpace(c.CapabilityID) == "" || !c.RationalePresent || c.MissingStage != "" {
			return fmt.Errorf("proposed review candidate is incomplete")
		}
	default:
		return fmt.Errorf("review candidate status %q is invalid", c.Status)
	}
	if c.digest() != c.CandidateDigest {
		return fmt.Errorf("review candidate digest does not match")
	}
	return nil
}

func (c ReviewCandidate) digest() string {
	return digest(
		"capability-review-candidate",
		c.FeedbackDigest,
		c.CapabilityID,
		c.RationaleDigest,
		fmt.Sprintf("%t", c.RationalePresent),
		string(c.Status),
		c.MissingStage,
	)
}
