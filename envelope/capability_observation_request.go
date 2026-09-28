package envelope

import "fmt"

// CapabilityRequestFor projects an AVAILABLE observation into a request-shaped
// value while preserving the declaration binding. It never grants authority.
func (e CapabilityObservationEnvelope) CapabilityRequestFor(subject, audience, capabilityID string) (CapabilityRequest, error) {
	if err := e.Validate(); err != nil {
		return CapabilityRequest{}, err
	}
	for _, capability := range e.Capabilities {
		if capability.ID != capabilityID {
			continue
		}
		if capability.State != CapabilityObservationAvailable {
			return CapabilityRequest{}, fmt.Errorf("capability %q is observed as %s", capabilityID, capability.State)
		}
		request := CapabilityRequest{
			Subject:           subject,
			Audience:          audience,
			Capability:        capabilityID,
			DeclarationDigest: e.DeclarationDigest,
		}
		return request, request.Validate()
	}
	return CapabilityRequest{}, fmt.Errorf("capability %q was not observed", capabilityID)
}
