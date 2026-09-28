package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

type input struct {
	Discovery capability.Discovery   `json:"discovery"`
	Feedback  capability.FeedbackInput `json:"feedback"`
}

type output struct {
	Feedback capability.Feedback `json:"feedback"`
}

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fatal(err)
		}
		defer file.Close()
		reader = file
	}

	var request input
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		fatal(err)
	}
	feedback := capability.ObserveFeedback(request.Discovery, request.Feedback)
	if err := feedback.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(output{Feedback: feedback}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
