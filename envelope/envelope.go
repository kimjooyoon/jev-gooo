// Package envelope defines a provider-neutral, non-executing evidence boundary
// between JEV decisions and Gooo runtime observations.
package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Status string

const (
	StatusUnknown   Status = "UNKNOWN"
	StatusCompleted Status = "COMPLETED"
	StatusObserved  Status = "OBSERVED"
)

type Declaration struct {
	Source string
	Digest string
}

func BindDeclaration(source string) (Declaration, error) {
	if strings.TrimSpace(source) == "" {
		return Declaration{}, fmt.Errorf("gooo declaration source is required")
	}
	declaration := Declaration{Source: source, Digest: digest("gooo-declaration", source)}
	return declaration, declaration.Validate()
}

func (d Declaration) Validate() error {
	if strings.TrimSpace(d.Source) == "" || !validDigest(d.Digest) {
		return fmt.Errorf("gooo declaration binding is incomplete")
	}
	if digest("gooo-declaration", d.Source) != d.Digest {
		return fmt.Errorf("gooo declaration digest does not match source")
	}
	return nil
}

type CapabilityRequest struct {
	Subject           string
	Audience          string
	Capability        string
	DeclarationDigest string
}

func NewCapabilityRequest(declaration Declaration, subject, audience, capability string) (CapabilityRequest, error) {
	if err := declaration.Validate(); err != nil {
		return CapabilityRequest{}, err
	}
	request := CapabilityRequest{Subject: subject, Audience: audience, Capability: capability, DeclarationDigest: declaration.Digest}
	return request, request.Validate()
}

func (r CapabilityRequest) Validate() error {
	for name, value := range map[string]string{
		"subject":            r.Subject,
		"audience":           r.Audience,
		"capability":         r.Capability,
		"declaration_digest": r.DeclarationDigest,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("capability request %s is required", name)
		}
	}
	if !validDigest(r.DeclarationDigest) {
		return fmt.Errorf("capability request declaration digest is invalid")
	}
	return nil
}

func (r CapabilityRequest) Digest() string {
	return digest("capability-request", r.Subject, r.Audience, r.Capability, r.DeclarationDigest)
}

type WorkloadIdentity struct {
	SPIFFEID string
}

func (i WorkloadIdentity) Digest() string {
	return digest("workload-identity", i.SPIFFEID)
}

type DecisionReceipt struct {
	Digest         string
	NonAuthorizing bool
}

type ExecutionReceipt struct {
	RequestDigest         string
	GrantDigest           string
	IdentityDigest        string
	DecisionReceiptDigest string
	ResultDigest          string
	Terminal              bool
	Status                Status
	MissingStage          string
	Digest                string
}

// NewExecutionReceipt observes evidence only. It never executes a capability,
// issues a grant, or treats a JEV decision as authorization.
func NewExecutionReceipt(request CapabilityRequest, grantDigest string, identity WorkloadIdentity, decision DecisionReceipt, resultDigest string, terminal bool) (ExecutionReceipt, error) {
	if err := request.Validate(); err != nil {
		return ExecutionReceipt{}, err
	}
	if grantDigest != "" && !validDigest(grantDigest) {
		return ExecutionReceipt{}, fmt.Errorf("capability grant digest is invalid")
	}
	if decision.Digest != "" && !validDigest(decision.Digest) {
		return ExecutionReceipt{}, fmt.Errorf("decision receipt digest is invalid")
	}
	if resultDigest != "" && !validDigest(resultDigest) {
		return ExecutionReceipt{}, fmt.Errorf("result digest is invalid")
	}

	receipt := ExecutionReceipt{
		RequestDigest:         request.Digest(),
		GrantDigest:           grantDigest,
		IdentityDigest:        identity.Digest(),
		DecisionReceiptDigest: decision.Digest,
		ResultDigest:          resultDigest,
		Terminal:              terminal,
		Status:                StatusUnknown,
	}
	switch {
	case grantDigest == "":
		receipt.MissingStage = "capability_grant"
	case strings.TrimSpace(identity.SPIFFEID) == "":
		receipt.MissingStage = "workload_identity"
	case decision.Digest == "" || !decision.NonAuthorizing:
		receipt.MissingStage = "decision_receipt"
	case resultDigest == "" || !terminal:
		receipt.MissingStage = "terminal_result"
	default:
		receipt.Status = StatusCompleted
	}
	receipt.Digest = receipt.digest()
	return receipt, nil
}

func (r ExecutionReceipt) Validate() error {
	if !validDigest(r.RequestDigest) || !validDigest(r.IdentityDigest) || !validDigest(r.Digest) {
		return fmt.Errorf("execution receipt digest is invalid")
	}
	if r.Status != StatusUnknown && r.Status != StatusCompleted {
		return fmt.Errorf("execution receipt status %q is invalid", r.Status)
	}
	if r.Status == StatusUnknown && strings.TrimSpace(r.MissingStage) == "" {
		return fmt.Errorf("unknown execution receipt must preserve missing stage")
	}
	if r.Status == StatusCompleted {
		if r.MissingStage != "" || !r.Terminal || !validDigest(r.GrantDigest) || !validDigest(r.DecisionReceiptDigest) || !validDigest(r.ResultDigest) {
			return fmt.Errorf("completed execution receipt is incomplete")
		}
	}
	if r.digest() != r.Digest {
		return fmt.Errorf("execution receipt evidence digest does not match")
	}
	return nil
}

func (r ExecutionReceipt) digest() string {
	return digest("execution-receipt", r.RequestDigest, r.GrantDigest, r.IdentityDigest, r.DecisionReceiptDigest, r.ResultDigest, fmt.Sprintf("%t", r.Terminal), string(r.Status), r.MissingStage)
}

type ReverseObservation struct {
	ExecutionDigest      string
	ObservedResultDigest  string
	VerifierDigest       string
	Status               Status
	MissingStage         string
	Digest               string
}

// ObserveReverse records an external reverse observation without executing,
// mutating source, or promoting an incomplete receipt.
func ObserveReverse(receipt ExecutionReceipt, observedResultDigest, verifierDigest string) ReverseObservation {
	observation := ReverseObservation{
		ExecutionDigest:     receipt.Digest,
		ObservedResultDigest: observedResultDigest,
		VerifierDigest:      verifierDigest,
		Status:              StatusUnknown,
	}
	switch {
	case receipt.Validate() != nil:
		observation.MissingStage = "execution_receipt"
	case receipt.Status != StatusCompleted:
		observation.MissingStage = "execution_receipt:" + receipt.MissingStage
	case !validDigest(observedResultDigest):
		observation.MissingStage = "observed_result"
	case !validDigest(verifierDigest):
		observation.MissingStage = "verifier"
	default:
		observation.Status = StatusObserved
	}
	observation.Digest = observation.digest()
	return observation
}

func (o ReverseObservation) Validate() error {
	if !validDigest(o.ExecutionDigest) || !validDigest(o.Digest) {
		return fmt.Errorf("reverse observation digest is invalid")
	}
	if o.Status != StatusUnknown && o.Status != StatusObserved {
		return fmt.Errorf("reverse observation status %q is invalid", o.Status)
	}
	if o.Status == StatusUnknown && strings.TrimSpace(o.MissingStage) == "" {
		return fmt.Errorf("unknown reverse observation must preserve missing stage")
	}
	if o.Status == StatusObserved && (o.MissingStage != "" || !validDigest(o.ObservedResultDigest) || !validDigest(o.VerifierDigest)) {
		return fmt.Errorf("observed reverse observation is incomplete")
	}
	if o.digest() != o.Digest {
		return fmt.Errorf("reverse observation digest does not match")
	}
	return nil
}

func (o ReverseObservation) digest() string {
	return digest("reverse-observation", o.ExecutionDigest, o.ObservedResultDigest, o.VerifierDigest, string(o.Status), o.MissingStage)
}

func digest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
