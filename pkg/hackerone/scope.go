package hackerone

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kenjoe41/h1scope/pkg/options"
)

func ProcessScope(scope *Scope, opt options.Options) {
	// TODO: Consider using a chan to write output.

	// if scope == nil {
	// 	return nil
	// }

	ProcessProgramScope(*scope, opt)
}

func ProcessProgramScope(scope Scope, opt options.Options) {

	for _, asset := range scope.Relationships.StructuredScopes.Data {

		// Out of Scope, TODO: implement flag for it if there is ever a need.
		if asset.Attributes.EligibleForBounty == false {
			// Skip out of scope
			continue
		}

		identifier := asset.Attributes.Identifier
		assetType := asset.Attributes.AssetType

		if assetType == "URL" {
			if (opt.Wildcard || opt.ALL) && strings.HasPrefix(identifier, "*") {
				if opt.CleanWildcard {
					identifier = cleanDomain(identifier)
					fmt.Println(identifier)
				} else {
					fmt.Println(identifier)
				}
				continue
			} else if (opt.Domains || opt.ALL) && !strings.HasPrefix(identifier, "*") {
				fmt.Println(identifier)
			}
		} else if (opt.CIDR || opt.ALL) && assetType == "CIDR" {
			fmt.Println(identifier)
		} else if (opt.Code || opt.ALL) && assetType == "SOURCE_CODE" {
			fmt.Println(identifier)
		} else if (opt.Android || opt.ALL) && assetType == "GOOGLE_PLAY_APP_ID" {
			fmt.Println(identifier)
		} else if (opt.APK || opt.ALL) && assetType == "OTHER_APK" {
			fmt.Println(identifier)
		} else if (opt.IOS || opt.ALL) && assetType == "APPLE_STORE_APP_ID" {
			fmt.Println(identifier)
		} else if (opt.IPA || opt.ALL) && assetType == "OTHER_IPA" {
			fmt.Println(identifier)
		} else if (opt.Other || opt.ALL) && assetType == "OTHER" {
			fmt.Println(identifier)
		} else if (opt.Hardware || opt.ALL) && assetType == "HARDWARE" {
			fmt.Println(identifier)
		} else if (opt.Windows || opt.ALL) && assetType == "WINDOWS_APP_STORE_APP_ID" {
			fmt.Println(identifier)
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
