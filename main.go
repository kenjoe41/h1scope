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

	programsChan := make(chan string)
	output := make(chan string)

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer outputWG.Done()

		for scopeAsset := range output {
			fmt.Println(scopeAsset)
		}

	}()

	if flags.Handle != "" {
		hackerone.GetProgramScope(output, flags)
	} else {
		hackerone.GetProgramsScope(programsChan, output, flags)
	}

	go func() {
		outputWG.Wait()
		close(output)
	}()
}
