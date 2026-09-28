package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/jev-gooo/capability"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(capability.Catalog()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
