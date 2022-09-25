package options

import (
	"flag"
	"os"
)

type Options struct {
	Username      string
	Apikey        string
	Concurrency   uint16
	Handle        string
	Wildcard      bool
	CleanWildcard bool
	Domains       bool
	CIDR          bool
	Android       bool
	IOS           bool
	Code          bool
	Other         bool
	APK           bool
	IPA           bool
	Hardware      bool
	Windows       bool
	ALL           bool
	Private       bool
	Public        bool
	Vdp           bool
	Paid          bool
	IncludeHandle bool
}

func ScanFlag() Options {
	usernamePtr := flag.String("u", os.Getenv("H1_USERNAME"), "Hackerone Username.")
	apikeyPtr := flag.String("apikey", os.Getenv("H1_APIKEY"), "Generate APIKEY from https://hackerone.com/settings/api_token/edit")
	concurrencyPtr := flag.Uint64("c", 5, "Concurrent programs to process.")
	handlePtr := flag.String("handle", "", "Handle for a specific program.")
	wildcardPtr := flag.Bool("wildcard", false, "Get wildcard domains.")
	cleanwildcardPtr := flag.Bool("cw", false, "Clean wildcard domains to pipe to recon tools, *.hackerone.com => hackerone.com")
	domainsPtr := flag.Bool("domains", false, "")
	cidrPtr := flag.Bool("cidr", false, "")
	androidPtr := flag.Bool("android", false, "")
	iosPtr := flag.Bool("ios", false, "")
	codePtr := flag.Bool("code", false, "")
	otherPtr := flag.Bool("other", false, "")
	apkPtr := flag.Bool("apk", false, "")
	ipaPtr := flag.Bool("ipa", false, "")
	hardwarePtr := flag.Bool("hardware", false, "")
	windowsPtr := flag.Bool("windows", false, "")
	allPtr := flag.Bool("all", false, "Get all scopes for a program.")
	privatePtr := flag.Bool("private", false, "Get scope for private programs.")
	publicPtr := flag.Bool("public", false, "Get scope for public programs.")
	vdpPtr := flag.Bool("vdp", false, "Get scope for free VDP programs.")
	paidPtr := flag.Bool("paid", false, "Get scope for Paid Programs.")
	includeHandlePtr := flag.Bool("include-handle", false, "Include Hadle of the program per asset printout e.g. hackerone, vpn.hackerone.com")

	flag.Parse()

	result := Options{
		*usernamePtr,
		*apikeyPtr,
		uint16(*concurrencyPtr),
		*handlePtr,
		*wildcardPtr,
		*cleanwildcardPtr,
		*domainsPtr,
		*cidrPtr,
		*androidPtr,
		*iosPtr,
		*codePtr,
		*otherPtr,
		*apkPtr,
		*ipaPtr,
		*hardwarePtr,
		*windowsPtr,
		*allPtr,
		*privatePtr,
		*publicPtr,
		*vdpPtr,
		*paidPtr,
		*includeHandlePtr,
	}

	if !result.Wildcard && !result.Domains && !result.CIDR && !result.Android && !result.IOS && !result.Code && !result.Other && !result.APK && !result.IPA && !result.Hardware && !result.Windows {
		result.ALL = true
	}

	return result
}

func Usage() {
	flag.Usage()
}
