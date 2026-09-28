package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type capabilityReceiptReport struct {
	Valid             bool                                      `json:"valid"`
	Status            envelope.CapabilityDiscoveryReceiptState  `json:"status"`
	ReceiptDigest     string                                    `json:"receipt_digest,omitempty"`
	FirstMissingStage string                                    `json:"first_missing_stage,omitempty"`
	NextOperation     string                                    `json:"next_operation,omitempty"`
	NonExecuting      bool                                      `json:"non_executing"`
	NonAuthorizing    bool                                      `json:"non_authorizing"`
	Error             string                                    `json:"error,omitempty"`
}

func main() {
	payload, err := input()
	if err != nil {
		fail(err)
	}
	var receipt envelope.CapabilityDiscoveryReceipt
	if err := json.Unmarshal(payload, &receipt); err != nil {
		fail(err)
	}
	report := capabilityReceiptReport{
		Status: receipt.Status,
		FirstMissingStage: receipt.FirstMissingStage,
		NextOperation: receipt.NextOperation,
		NonExecuting: receipt.NonExecuting,
		NonAuthorizing: receipt.NonAuthorizing,
	}
	if err := receipt.Validate(); err != nil {
		report.Error = err.Error()
		write(report)
		os.Exit(2)
	}
	report.Valid = true
	report.ReceiptDigest, err = receipt.CanonicalDigest()
	if err != nil {
		fail(err)
	}
	write(report)
}

func input() ([]byte, error) {
	if len(os.Args) == 2 {
		return os.ReadFile(os.Args[1])
	}
	return io.ReadAll(os.Stdin)
}

func write(value capabilityReceiptReport) {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fail(err)
	}
	fmt.Println(string(payload))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
