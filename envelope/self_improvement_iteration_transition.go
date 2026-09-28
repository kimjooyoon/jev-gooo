package envelope

import "fmt"

func (i SelfImprovementIteration) Accept(verificationDigest, rollbackDigest string) (SelfImprovementIteration, error) {
	if i.State != SelfImprovementIterationProposed {
		return SelfImprovementIteration{}, fmt.Errorf("only PROPOSED iterations can be accepted")
	}
	i.VerificationDigest = verificationDigest
	i.RollbackDigest = rollbackDigest
	i.State = SelfImprovementIterationAccepted
	if err := i.Validate(); err != nil {
		return SelfImprovementIteration{}, err
	}
	return i, nil
}

func (i SelfImprovementIteration) Reject(verificationDigest, rollbackDigest string) (SelfImprovementIteration, error) {
	if i.State != SelfImprovementIterationProposed {
		return SelfImprovementIteration{}, fmt.Errorf("only PROPOSED iterations can be rejected")
	}
	i.VerificationDigest = verificationDigest
	i.RollbackDigest = rollbackDigest
	i.State = SelfImprovementIterationRejected
	if err := i.Validate(); err != nil {
		return SelfImprovementIteration{}, err
	}
	return i, nil
}

