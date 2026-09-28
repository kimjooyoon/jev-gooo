package main

import (
	"encoding/json"
	"flag"
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
	query := flag.String("query", "", "natural-language capability question")
	declarationPath := flag.String("declaration", "", "path to a .gooo declaration")
	flag.Parse()

	if *query != "" || *declarationPath != "" {
		if *query == "" || *declarationPath == "" {
			fatal(fmt.Errorf("--query and --declaration must be provided together"))
		}
		source, err := os.ReadFile(*declarationPath)
		if err != nil {
			fatal(err)
		}
		declaration, err := envelope.BindDeclaration(string(source))
		if err != nil {
			fatal(err)
		}
		emit(*query, declaration)
		return
	}

	var reader io.Reader = os.Stdin
	if len(flag.Args()) > 1 {
		fatal(fmt.Errorf("at most one JSON input path is supported"))
	}
	if len(flag.Args()) == 1 {
		file, err := os.Open(flag.Args()[0])
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
	emit(request.Query, declaration)
}

func emit(query string, declaration envelope.Declaration) {
	assistant, err := capability.DiscoverAssistant(query, declaration)
	if err != nil {
		fatal(err)
	}
	if err := assistant.Validate(); err != nil {
		fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(assistant); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
