package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
	"github.com/kimjooyoon/jev-gooo/envelope"
)

type input struct {
	Query       string `json:"query"`
	Declaration string `json:"declaration"`
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
	declaration, err := envelope.BindDeclaration(request.Declaration)
	if err != nil {
		fatal(err)
	}
	overview, err := capability.DiscoverOverview(request.Query, declaration)
	if err != nil {
		fatal(err)
	}
	if err := overview.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(overview); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
