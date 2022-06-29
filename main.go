package main

import (
	"fmt"
	"os"

	"github.com/kenjoe41/h1scope/hackerone"
)

func main() {
	flags := hackerone.ScanFlag()

	if flags.Username == "" || flags.Apikey == "" {
		fmt.Println("H1 Username and API key are needed.")
		hackerone.Usage()
		os.Exit(1)
	}

	if flags.Handle == "" {
		fmt.Println("No program handle specified.")
		hackerone.Usage()
		os.Exit(1)
	}

	scope, _ := hackerone.GetProgramScope(flags)

	for _, data := range scope.Relationships.StructuredScopes.Data {
		fmt.Println(data.Attributes.Identifier)
	}

}
