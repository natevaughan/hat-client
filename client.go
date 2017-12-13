package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"strconv"
	"fmt"
	"bytes"
)

type ChatClient struct {
	baseUrl string
	token string
	client http.Client
}

func (cc *ChatClient) identity() User {

	body, err := cc.get("user")

	if err != nil {
		printErr(err.Error())
		os.Exit(-1)
	}

	var identity User

	err = json.Unmarshal(body, &identity)

	if err != nil {
		printErr(err.Error())
		os.Exit(-1)
	}
	return identity
}

func (cc *ChatClient) listHats() []Hat {

	body, err := cc.get("hat")

	if err != nil {
		printErr(err.Error())
	}
	println(string(body))

	var m []Hat

	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			printErr(err.Error())
		}
	}

	return m
}

func (cc *ChatClient) previous(hatId int64, count int) []Message {

	path := "hat/" + fmt.Sprintf("%v", hatId) + "/previous/" + strconv.Itoa(count)
	body, err := cc.get(path)

	if err != nil {
		printErr(err.Error())
	}
	var m []Message

	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			printErr(err.Error())
		}
	}

	return m
}


func (cc *ChatClient) since(hatId int64, unix int64) []Message {

	path := "hat/" + fmt.Sprintf("%v", hatId) + "/since/" + strconv.FormatInt(unix, 10)

	body, err := cc.get(path)

	if err != nil {
		printErr(err.Error())
	}

	var m []Message

	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			printErr(err.Error())
		}
	}

	return m
}

func (cc *ChatClient) sendMessage(hatId int64, message string) []byte {
	path := "hat/" + fmt.Sprintf("%v", hatId)

	msg := MessageReq{0, message}

	postBody, err := json.Marshal(msg)

	if err != nil {
		printErr(err.Error())
	}

	body, err := cc.post(path, postBody)

	if err != nil {
		printErr(err.Error())
	}

	return body
}

func (cc *ChatClient) get(path string) ([]byte, error) {

	req, err := http.NewRequest("GET", cc.baseUrl + path, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("api-token", cc.token)
	req.Header.Add("Accept", "application/json")

	resp, err := cc.client.Do(req)
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


func (cc *ChatClient) post(path string, postBody []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", cc.baseUrl + path, bytes.NewBuffer(postBody))

	if err != nil {
		return nil, err
	}

	req.Header.Add("api-token", cc.token)
	req.Header.Add("Accept", "application/json")

	resp, err := cc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}


type Message struct {
	Id int64
	Timestamp int64
	Text string
	Author User
	LastEdited int64
}

type MessageReq struct {
	Id int64
	Text string
}

type User struct {
	Id int64
	Name string
}


type Hat struct {
	Id int64
	Name string
	Creator User
	Participants []User
}

