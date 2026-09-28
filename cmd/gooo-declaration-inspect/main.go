package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

type input struct {
	Source string `json:"source"`
}

type output struct {
	Shape capability.DeclarationShape `json:"shape"`
}

func main() {
	var reader io.Reader = os.Stdin
	var request input
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		fatal(err)
	}
	shape, err := capability.InspectDeclaration(request.Source)
	if err != nil {
		fatal(err)
	}
	if err := shape.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(output{Shape: shape}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
