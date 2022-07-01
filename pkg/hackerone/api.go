package hackerone

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/kenjoe41/h1scope/pkg/options"
)

func GetProgramScope(output chan string, opt options.Options) error {
	link := fmt.Sprintf("https://api.hackerone.com/v1/hackers/programs/%s", opt.Handle)

	resp, err := makeAPIRequest(link, opt)
	if err != nil {
		return fmt.Errorf("Error making AI HTTP request: %s\n", err)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("client: could not read response body: %s\n", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-2XX response from server: %s", err)
	}
	scope, err := Unmarshal(resBody)
	if err != nil {
		return err
	}
	ProcessScope(*scope, opt, output)

	return nil
}

func makeAPIRequest(link string, opt options.Options) (*http.Response, error) {

	client := retryablehttp.NewClient()
	client.Logger = nil

	req, _ := retryablehttp.NewRequest("GET", link, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(opt))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
