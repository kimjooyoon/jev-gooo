package envelope

import "fmt"

func (i SelfImprovementIteration) ProposeNext(iterationID, candidateDigest string) (SelfImprovementIteration, error) {
	if i.State != SelfImprovementIterationUnknown {
		return SelfImprovementIteration{}, fmt.Errorf("only UNKNOWN iterations can propose a next candidate")
	}
	if iterationID == "" {
		return SelfImprovementIteration{}, fmt.Errorf("iteration_id is required for the next proposal")
	}
	if !validSelfImprovementDigest(candidateDigest) {
		return SelfImprovementIteration{}, fmt.Errorf("candidate_digest must be a sha256 digest")
	}
	i.IterationID = iterationID
	i.CandidateDigest = candidateDigest
	i.VerificationDigest = ""
	i.RollbackDigest = ""
	i.FirstMissingStage = ""
	i.NextOperation = ""
	i.State = SelfImprovementIterationProposed
	if err := i.Validate(); err != nil {
		return SelfImprovementIteration{}, err
	}
	return i, nil
}

