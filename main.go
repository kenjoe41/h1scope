package main

import (
	"fmt"
	"os"

	"github.com/kenjoe41/h1scope/pkg/hackerone"
	"github.com/kenjoe41/h1scope/pkg/options"
)

func main() {
	flags := options.ScanFlag()

	if flags.Username == "" || flags.Apikey == "" {
		fmt.Fprintln(os.Stderr, "H1 Username and API key are needed.")
		options.Usage()
		os.Exit(1)
	}

	if flags.Handle == "" {
		fmt.Fprintln(os.Stderr, "No program handle specified.")
		options.Usage()
		os.Exit(1)
	}

	scope, err := hackerone.GetProgramScope(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when fetching scope: %s\n", err)
	}

	for _, data := range scope.Relationships.StructuredScopes.Data {
		fmt.Println(data.Attributes.Identifier)
	}

}
