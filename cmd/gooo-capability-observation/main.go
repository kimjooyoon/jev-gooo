package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kimjooyoon/jev-gooo/capability"
	"github.com/kimjooyoon/jev-gooo/envelope"
)

type input struct {
	Query               string `json:"query"`
	Declaration         string `json:"declaration"`
	CorrelationID       string `json:"correlation_id"`
	CatalogDigest       string `json:"catalog_digest"`
	SchemaDigest        string `json:"schema_digest"`
	FirstMissingStage   int    `json:"first_missing_stage"`
	RequestCapabilityID string `json:"request_capability_id,omitempty"`
	Subject             string `json:"subject,omitempty"`
	Audience            string `json:"audience,omitempty"`
}

type output struct {
	Discovery   capability.Discovery                  `json:"discovery"`
	Observation envelope.CapabilityObservationEnvelope `json:"observation"`
	Request     *envelope.CapabilityRequest            `json:"request,omitempty"`
}

func main() {
	reader, err := inputReader()
	if err != nil {
		fatal(err)
	}
	defer reader.close()

	var request input
	if err := json.NewDecoder(reader.reader).Decode(&request); err != nil {
		fatal(err)
	}
	declaration, err := envelope.BindDeclaration(request.Declaration)
	if err != nil {
		fatal(err)
	}
	discovery, err := capability.Discover(request.Query, declaration)
	if err != nil {
		fatal(err)
	}
	if err := discovery.Validate(); err != nil {
		fatal(err)
	}

	state := envelope.CapabilityObservationUnknown
	capabilityID := discovery.CapabilityID
	firstMissingStage := request.FirstMissingStage
	nextOperation := "ask_narrower_question"
	switch discovery.Status {
	case capability.StatusAvailable:
		state = envelope.CapabilityObservationAvailable
		firstMissingStage = -1
		nextOperation = ""
	case capability.StatusDeferred:
		state = envelope.CapabilityObservationDeferred
		if firstMissingStage < 0 {
			fatal(fmt.Errorf("deferred discovery requires first_missing_stage"))
		}
		nextOperation = "provide_missing_declaration_signal"
	case capability.StatusUnknown:
		capabilityID = "capability.discovery"
		if firstMissingStage < 0 {
			fatal(fmt.Errorf("unknown discovery requires first_missing_stage"))
		}
	default:
		fatal(fmt.Errorf("unsupported discovery status %q", discovery.Status))
	}

	observation, err := envelope.NewCapabilityObservationEnvelope(
		request.CorrelationID,
		declaration.Digest,
		declaration.Digest,
		request.CatalogDigest,
		request.SchemaDigest,
		discovery.EvidenceDigest,
		[]envelope.CapabilityObservation{{ID: capabilityID, State: state}},
		firstMissingStage,
		nextOperation,
	)
	if err != nil {
		fatal(err)
	}

	var capabilityRequest *envelope.CapabilityRequest
	if strings.TrimSpace(request.RequestCapabilityID) != "" {
		if discovery.Status != capability.StatusAvailable {
			fatal(fmt.Errorf("request-shaped output requires AVAILABLE discovery"))
		}
		value, err := observation.CapabilityRequestFor(request.Subject, request.Audience, request.RequestCapabilityID)
		if err != nil {
			fatal(err)
		}
		capabilityRequest = &value
	}

	if err := json.NewEncoder(os.Stdout).Encode(output{
		Discovery:   discovery,
		Observation: observation,
		Request:     capabilityRequest,
	}); err != nil {
		fatal(err)
	}
}

type readerHandle struct {
	reader io.Reader
	close  func()
}

func inputReader() (readerHandle, error) {
	if len(os.Args) == 1 {
		return readerHandle{reader: os.Stdin, close: func() {}}, nil
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		return readerHandle{}, err
	}
	return readerHandle{reader: file, close: func() { _ = file.Close() }}, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
