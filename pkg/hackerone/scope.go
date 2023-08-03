package hackerone

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/kenjoe41/h1scope/pkg/options"
)

func GetProgramsScope(programsChan chan string, outputChan chan string, opt options.Options) {
	link := "https://api.hackerone.com/v1/hackers/programs"

	processPrograms(link, programsChan, outputChan, opt)

}

func processPrograms(link string, programsChan chan string, outputChan chan string, opt options.Options) {
	var programScopeWG sync.WaitGroup
	if opt.Handle != "" {
		programScopeWG.Add(1)
		goProcess(&programScopeWG, programsChan, outputChan, opt)
	} else {
		for i := 1; i <= int(opt.Concurrency); i++ {
			programScopeWG.Add(1)
			goProcess(&programScopeWG, programsChan, outputChan, opt)
		}
	}

	for {
		programs, err := getPrograms(link, opt)
		if err != nil {
			if programs.Links.Next == "" {
				break
			} else {
				continue
			}
		}

		for _, program := range programs.ProgramsData {

			// Looks like we have VDP
			if opt.Private && !program.ProgramsAttributes.OffersBounty {
				continue
			} else if opt.Vdp && program.ProgramsAttributes.OffersBounty {
				// For some God Almighty reason, someone wants only VDP.
				continue
			}

			// Do we want Private or Public programs
			if (opt.Private && program.ProgramsAttributes.State == "public_mode") {
				continue
			} else if (opt.Public && program.ProgramsAttributes.State != "public_mode") {
				continue
			}

			// Okay, we are all good for now.
			programsChan <- program.ProgramsAttributes.Handle

		}

		// Pagination, do we have another page or break out of the Infinite loop!?
		if programs.Links.Next == "" {
			close(programsChan)
			break

		} else {
			link = programs.Links.Next
		}
	}

	programScopeWG.Wait()
}

func goProcess(programScopeWG *sync.WaitGroup, programsChan chan string, outputChan chan string, opt options.Options) {
	go func() {
		defer programScopeWG.Done()

		for h_program := range programsChan {
			opt.Handle = h_program
			GetProgramScope(outputChan, opt)
		}

	}()
}

func getPrograms(link string, opt options.Options) (*Programs, error) {
	programs, err := processProgramsApiRequest(link, opt)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

func GetProgramScope(outputChan chan string, opt options.Options) error {
	link := fmt.Sprintf("https://api.hackerone.com/v1/hackers/programs/%s", opt.Handle)

	scope, err := processAPIRequest(link, opt)
	if err != nil {
		return err
	}

	ProcessProgramScope(*scope, opt, outputChan)

	return nil
}

func ProcessProgramScope(scope Scope, opt options.Options, outputChan chan string) {

	for _, asset := range scope.Relationships.StructuredScopes.ScopeData {

		// Out of Scope, TODO: implement flag for it if there is ever a need.
		if !asset.Attributes.EligibleForBounty {
			// Skip out of scope
			continue
		}

		identifier := asset.Attributes.Identifier
		assetType := asset.Attributes.AssetType

		if assetType == "URL" {
			if (opt.Wildcard || opt.ALL) && strings.HasPrefix(identifier, "*") {

				handleDomainIdentifier(identifier, outputChan, opt)

				continue
			} else if (opt.Domains || opt.ALL) && !strings.HasPrefix(identifier, "*") {

				handleDomainIdentifier(identifier, outputChan, opt)
			}
		} else if (opt.CIDR || opt.ALL) && assetType == "CIDR" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.Code || opt.ALL) && assetType == "SOURCE_CODE" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.Android || opt.ALL) && assetType == "GOOGLE_PLAY_APP_ID" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.APK || opt.ALL) && assetType == "OTHER_APK" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.IOS || opt.ALL) && assetType == "APPLE_STORE_APP_ID" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.IPA || opt.ALL) && assetType == "OTHER_IPA" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.Other || opt.ALL) && assetType == "OTHER" {
			if strings.HasPrefix(identifier, "*") {
				handleDomainIdentifier(identifier, outputChan, opt)
			} else {
				handleAsset(identifier, outputChan, opt)
			}
		} else if (opt.Hardware || opt.ALL) && assetType == "HARDWARE" {
			handleAsset(identifier, outputChan, opt)
		} else if (opt.Windows || opt.ALL) && assetType == "WINDOWS_APP_STORE_APP_ID" {
			handleAsset(identifier, outputChan, opt)
		}
	}
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

func domainSplitTrimSpace(domain string) []string {
	domainSlice := strings.Split(domain, ",")
	for i := range domainSlice {
		domainSlice[i] = strings.TrimSpace(domainSlice[i])
	}

	return domainSlice
}

func handleAsset(identifier string, outputChan chan string, opt options.Options) {
	domainsSlice := domainSplitTrimSpace(identifier)
	for _, identifier := range domainsSlice {
		if opt.IncludeHandle {
			outputChan <- fmt.Sprintf("%s, %s", opt.Handle, identifier)
		} else {
			outputChan <- identifier
		}
	}
}
func handleDomainIdentifier(identifier string, outputChan chan string, opt options.Options) {
	if opt.CleanWildcard {
		identifier = cleanDomain(identifier)
		domainsSlice := domainSplitTrimSpace(identifier)
		for _, identifier := range domainsSlice {
			if opt.IncludeHandle {
				outputChan <- fmt.Sprintf("%s, %s", opt.Handle, identifier)
			} else {
				outputChan <- identifier
			}
		}
	} else {
		// Lets clean up the domain abit e.g. "site.com, site2.com"
		domainsSlice := domainSplitTrimSpace(identifier)
		for _, identifier := range domainsSlice {
			if opt.IncludeHandle {
				outputChan <- fmt.Sprintf("%s, %s", opt.Handle, identifier)
			} else {
				outputChan <- identifier
			}
		}
	}

}
