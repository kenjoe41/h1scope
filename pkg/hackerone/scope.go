package hackerone

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/kenjoe41/h1scope/pkg/options"
)

func GetProgramsScope(programsChan chan string, output chan string, opt options.Options) error {
	link := fmt.Sprintf("https://api.hackerone.com/v1/hackers/programs")

	processPrograms(link, programsChan, output, opt)

	close(programsChan)

	return nil
}

func processPrograms(link string, programsChan chan string, output chan string, opt options.Options) {

	var programScopeWG sync.WaitGroup
	programScopeWG.Add(1)
	go func() {
		defer programScopeWG.Done()

		for h_program := range programsChan {
			opt.Handle = h_program
			GetProgramScope(output, opt)
		}

	}()

	for {
		programs, err := getPrograms(link, opt)
		if err != nil {
			continue
		}

		for _, program := range programs.ProgramsData {

			// Looks like we have VDP
			if opt.Private && program.ProgramsAttributes.OffersBounty != true {
				continue
			} else if opt.Vdp && program.ProgramsAttributes.OffersBounty == true {
				// For some God Almighty reason, someone wants only VDP.
				continue
			}

			// Do we want Private or Public programs
			if opt.Private && program.ProgramsAttributes.State == "public_mode" {
				continue
			} else if opt.Public && program.ProgramsAttributes.State != "public_mode" {
				continue
			}

			// Okay, we are all good for now.
			programsChan <- program.ProgramsAttributes.Handle

		}

		// Pagination, do we have another page or break out of the Infinite loop!?
		if programs.Links.Next == "" {
			break
		} else {
			link = programs.Links.Next
		}
	}
}

func getPrograms(link string, opt options.Options) (*Programs, error) {
	programs, err := processProgramsApiRequest(link, opt)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

func GetProgramScope(output chan string, opt options.Options) error {
	link := fmt.Sprintf("https://api.hackerone.com/v1/hackers/programs/%s", opt.Handle)

	scope, err := processAPIRequest(link, opt)
	if err != nil {
		return err
	}

	ProcessProgramScope(*scope, opt, output)

	return nil
}

func ProcessProgramScope(scope Scope, opt options.Options, output chan string) {

	for _, asset := range scope.Relationships.StructuredScopes.ScopeData {

		// Out of Scope, TODO: implement flag for it if there is ever a need.
		if asset.Attributes.EligibleForBounty == false {
			// Skip out of scope
			continue
		}

		identifier := asset.Attributes.Identifier
		assetType := asset.Attributes.AssetType

		if assetType == "URL" {
			if (opt.Wildcard || opt.ALL) && strings.HasPrefix(identifier, "*") {
				if opt.CleanWildcard { //TODO: run with handle security, seems it is not cleaning the *
					cleanidentifier := cleanDomain(identifier)
					output <- cleanidentifier
				} else {
					output <- identifier
				}
				continue
			} else if (opt.Domains || opt.ALL) && !strings.HasPrefix(identifier, "*") {

				output <- identifier
			}
		} else if (opt.CIDR || opt.ALL) && assetType == "CIDR" {
			output <- identifier
		} else if (opt.Code || opt.ALL) && assetType == "SOURCE_CODE" {
			output <- identifier
		} else if (opt.Android || opt.ALL) && assetType == "GOOGLE_PLAY_APP_ID" {
			output <- identifier
		} else if (opt.APK || opt.ALL) && assetType == "OTHER_APK" {
			output <- identifier
		} else if (opt.IOS || opt.ALL) && assetType == "APPLE_STORE_APP_ID" {
			output <- identifier
		} else if (opt.IPA || opt.ALL) && assetType == "OTHER_IPA" {
			output <- identifier
		} else if (opt.Other || opt.ALL) && assetType == "OTHER" {
			output <- identifier
		} else if (opt.Hardware || opt.ALL) && assetType == "HARDWARE" {
			output <- identifier
		} else if (opt.Windows || opt.ALL) && assetType == "WINDOWS_APP_STORE_APP_ID" {
			output <- identifier
		}
	}
	// For some fricking reason failed to solve the issue of noot printin last item on chan.
	output <- ""
}

func cleanDomain(domain string) string {

	pattern := `[\w]+[\w\-_~\.]+\.[a-zA-Z]+|$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		// Whatever happened, just return the original domain
		return domain
	}

	cDomain := r.FindString(domain)
	if cDomain != "" {
		return cDomain
	}
	return domain
}
