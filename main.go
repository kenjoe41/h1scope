package main

import (
	"fmt"
	"os"
	"sync"

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

	output := make(chan string)

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer close(output)

		for scopeAsset := range output {
			fmt.Println(scopeAsset)
		}
		outputWG.Done()
	}()

	if err := hackerone.GetProgramScope(output, flags); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when fetching scope: %s\n", err)
	}

	// close(output)

}
