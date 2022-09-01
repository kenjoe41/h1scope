package hackerone

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/kenjoe41/h1scope/pkg/options"
)

func processProgramsApiRequest(link string, opt options.Options) (*Programs, error) {
	resBody, err := makeAPIRequest(link, opt)
	if err != nil {
		return nil, fmt.Errorf("error making AI HTTP request: %s", err)
	}

	programs, err := UnmarshalPrograms(resBody)
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func processAPIRequest(link string, opt options.Options) (*Scope, error) {
	resBody, err := makeAPIRequest(link, opt)
	if err != nil {
		return nil, fmt.Errorf("error making AI HTTP request: %s", err)
	}

	scope, err := Unmarshal(resBody)
	if err != nil {
		return nil, err
	}
	return scope, nil
}

func makeAPIRequest(link string, opt options.Options) ([]byte, error) {

	client := retryablehttp.NewClient()
	client.Logger = nil

	req, _ := retryablehttp.NewRequest("GET", link, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(opt))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-2XX response from server: %s", err)
	}

	return resBody, nil
}

func basicAuth(opt options.Options) string {
	auth := opt.Username + ":" + opt.Apikey
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Unmarshal(jsonBytes []byte) (*Scope, error) {
	scope := new(Scope)

	if err := json.Unmarshal(jsonBytes, &scope); err != nil {
		return nil, err
	}
	return scope, nil
}

func UnmarshalPrograms(jsonBytes []byte) (*Programs, error) {
	programs := new(Programs)

	if err := json.Unmarshal(jsonBytes, &programs); err != nil {
		return nil, err
	}
	return programs, nil
}
