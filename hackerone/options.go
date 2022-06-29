package hackerone

import (
	"flag"
	"os"
)

type Options struct {
	Username string
	Apikey   string
	Handle   string
}

func ScanFlag() Options {
	usernamePtr := flag.String("u", os.Getenv("H1_USERNAME"), "Hackerone Username.")
	apikeyPtr := flag.String("apikey", os.Getenv("H1_APIKEY"), "Generate APIKEY from https://hackerone.com/settings/api_token/edit")
	handlePtr := flag.String("handle", "", "Handle for a specific program.")

	flag.Parse()

	result := Options{
		*usernamePtr,
		*apikeyPtr,
		*handlePtr,
	}
	return result
}

func Usage() {
	flag.Usage()
}
