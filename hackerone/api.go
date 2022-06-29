package hackerone

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetProgramScope(opt Options) (*Scope, error) {
	link := fmt.Sprintf("https://api.hackerone.com/v1/hackers/programs/%s", opt.Handle)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", link, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(opt))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making http request: %s\n", err)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s\n", err)
	}
	// fmt.Printf("client: response body: %s\n", resBody)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-2XX response from server: %s", err)
	}
	return Unmarshal(resBody), nil
}

func basicAuth(opt Options) string {
	auth := opt.Username + ":" + opt.Apikey
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Unmarshal(jsonBytes []byte) *Scope {
	scope := new(Scope)

	if err := json.Unmarshal(jsonBytes, &scope); err != nil {
		return nil
	}
	return scope
}
