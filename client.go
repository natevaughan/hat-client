package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ChatClient struct {
	baseUrl    string
	token      string
	client     http.Client
	reader     bufio.Reader
	outboundIP net.IP
}

const PUBLIC_CHANNEL = "PUBLIC_CHANNEL"
const PRIVATE_CHANNEL = "PRIVATE_CHANNEL"
const CONVERSATION = "CONVERSATION"

func (cc *ChatClient) getInput() string {
	text, _ := cc.reader.ReadString('\n')
	return strings.Join(strings.Fields(text), "")
}

func (cc *ChatClient) identity() User {

	body, err := cc.Get("user")

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

func (cc *ChatClient) listChats(space Space) []Chat {

	body, err := cc.Get("space/" + space.Id + "/chat")

	if err != nil {
		printErr(err.Error())
	}
	println(string(body))

	var m []Chat

	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			printErr(err.Error())
		}
	}

	return m
}

func (cc *ChatClient) createChat(space Space, name string, chatType int, participants []string) Chat {

	var t string
	switch chatType {
	case 1:
		t = PUBLIC_CHANNEL
		break
	case 2:
		t = PRIVATE_CHANNEL
		break
	default:
		t = CONVERSATION
		break
	}

	b, err := json.Marshal(CreateChatRequest{name, t, participants})

	if err != nil {
		printErr(err.Error())
	}

	println(string(b))

	body, err := cc.Post("space/"+space.Id+"/chat", b)

	if err != nil {
		printErr(err.Error())
	}

	println(string(body))

	var m Chat

	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			printErr(err.Error())
		}
	}

	return m
}

func (cc *ChatClient) listSpaces() []Space {

	body, err := cc.Get("space")

	if err != nil {
		printErr(err.Error())
	}
	println(string(body))

	var s []Space

	if len(body) > 0 {
		err = json.Unmarshal(body, &s)
		if err != nil {
			printErr(err.Error())
			return []Space{}
		}
	}

	return s
}

func (cc *ChatClient) previous(chatId string, count int) []Message {

	path := "hat/" + fmt.Sprintf("%v", chatId) + "/previous/" + strconv.Itoa(count)
	body, err := cc.Get(path)

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

func (cc *ChatClient) since(space Space, chatId string, unix int64) []Message {

	path := "space/" + space.Id + "/chat/" + chatId + "/since/" + strconv.FormatInt(unix, 10)

	body, err := cc.Get(path)

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

func (cc *ChatClient) sendMessage(chatId string, message string) []byte {
	path := "chat/" + chatId

	msg := MessageReq{"", message}

	postBody, err := json.Marshal(msg)

	if err != nil {
		printErr(err.Error())
	}

	body, err := cc.Post(path, postBody)

	if err != nil {
		printErr(err.Error())
	}

	return body
}

func (cc *ChatClient) getUsersForSpace(space Space) ([]Participant, error) {
	var p []Participant
	path := "user/space/" + space.Id

	body, err := cc.Get(path)

	if err != nil {
		return p, err
	}

	err = json.Unmarshal(body, &p)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (cc *ChatClient) Get(path string) ([]byte, error) {

	if cc.outboundIP == nil {
		cc.outboundIP = cc.getOutboundIP()
	}

	req, err := http.NewRequest("GET", cc.baseUrl+path, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("api-token", cc.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Forwarded-For", cc.outboundIP.String())

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

func (cc *ChatClient) Post(path string, postBody []byte) ([]byte, error) {

	if cc.outboundIP == nil {
		cc.outboundIP = cc.getOutboundIP()
	}

	req, err := http.NewRequest("POST", cc.baseUrl+path, bytes.NewBuffer(postBody))

	if err != nil {
		return nil, err
	}

	req.Header.Add("api-token", cc.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Forwarded-For", cc.outboundIP.String())

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

func (cc *ChatClient) getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type Message struct {
	Id         string
	Timestamp  int64
	Text       string
	Author     User
	LastEdited int64
}

type MessageReq struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type Participant struct {
	User User
	Role string
}

type User struct {
	Id   string
	Name string
}

type Chat struct {
	Id           string
	Name         string
	Type         string
	Space        Space
	Creator      User
	Participants []User
}

type Space struct {
	Id   string
	Name string
}

type CreateChatRequest struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Participants []string `json:"participants"`
}
