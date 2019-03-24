package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Setup struct {
	baseUrl string
	client  http.Client
	reader  bufio.Reader
}

type InviteReq struct {
	Token string `json:"token"`
	Alias string `json:"alias"`
}

type InviteResp struct {
	UserName     string `json:"userName"`
	ApiKey       string `json:"apiKey"`
	ErrorMessage string `json:"errorMessage"`
}

func (s *Setup) Setup() string {

	in := UserInput{BufIoWrapper{s.reader}}

	println("Welcome! It looks like you are using Kchat for the first time. Press <enter> on your keyboard to continue.")

	in.getInput()

	code, err := in.Ask("Please enter an invite code: ", 2, validateInvite, "Invite code appears to be invalid. Codes look something like 'XXX-XXX-XXX'")

	if err != nil {
		println("Unable to get code. Thank you for choosing Kchat. Goodbye.")
		return ""
	}

	alias, err := in.Ask("Please enter an alias (username) to go by: ", 3, notEmpty, "(enter at least one non-whitespace character)")

	return s.setupInvite(code, alias)
}

func (s *Setup) setupInvite(code string, alias string) string {

	inv := InviteReq{code, alias}

	reqBody, err := json.Marshal(inv)

	if err != nil {
		println(err.Error())
		return ""
	}

	body, err := s.Post("invite/redeem", reqBody)

	if err != nil {
		println("!!!")
	}

	var resp InviteResp

	err = json.Unmarshal(body, &resp)

	if err != nil {
		println(err.Error())
	}

	return resp.ApiKey
}

func (s *Setup) Post(path string, postBody []byte) ([]byte, error) {

	println("Attempting POST to " + s.baseUrl + path + " with body " + string(postBody))

	req, err := http.NewRequest("POST", s.baseUrl+path, bytes.NewBuffer(postBody))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func validateInvite(i string) bool {
	if strings.ContainsAny(i, "'\",./?<>;:`~!@#$%^&*()_=+|\\ []{}aeoiuyAEOIUY") {
		return false
	}
	if len(i) == 9 {
		return true
	} else if len(i) == 11 {
		return i[3] == '-' && i[7] == '-'
	}
	return false
}

func notEmpty(i string) bool {
	if i == "" {
		return false
	}
	return true
}
