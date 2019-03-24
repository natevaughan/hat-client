package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const KCHAT_SETTINGS_LOC = "~/Library/Application Support/kchat"
const KCHAT_SETTINGS = KCHAT_SETTINGS_LOC + "/kchat_settings.yml"

func main() {

	host := os.Getenv("KCHAT_HOST")
	apiKey := os.Getenv("KCHAT_TOKEN")

	if len(host) == 0 {
		host = "http://localhost:5000/"
		//printErr("you must have KCHAT_HOST set in your environment")
		//return
	}

	println("*****************************************")
	println(" ___   _ _______ __   __ _______ _______ ")
	println("|   | | |       |  | |  |   _   |       |")
	println("|   |_| |       |  |_|  |  |_|  |_     _|")
	println("|      _|       |       |       | |   |  ")
	println("|     |_|      _|       |       | |   |  ")
	println("|    _  |     |_|   _   |   _   | |   |  ")
	println("|___| |_|_______|__| |__|__| |__| |___|  ")
	println(" ")
	println("*****************************************")
	println(" ")

	reader := bufio.NewReader(os.Stdin)

	file, err := ioutil.ReadFile(KCHAT_SETTINGS)

	var props YamlProps

	if err == nil {
		err = yaml.Unmarshal(file, &props)
		if err == nil {
			apiKey = props.ApiKey
		}
	}

	hc := http.Client{}

	if apiKey == "" {
		setup := Setup{host, hc, *reader}
		apiKey = setup.Setup()
		if apiKey == "" {
			println("Unable to setup invite")
			return
		}
		props = YamlProps{apiKey}

		bytes, err := yaml.Marshal(props)

		if err != nil {
			println("Could not marshal yaml: " + err.Error())
		}

		if err == nil {
			err = os.Mkdir(KCHAT_SETTINGS_LOC, 0644)

			if err != nil {
				println("Could not write settings to file: " + err.Error())
			} else {

				err = ioutil.WriteFile(KCHAT_SETTINGS, bytes, 0644)

				if err != nil {
					println("Could not write settings to file: " + err.Error())
				}
			}
		}
	}

	client := ChatClient{host, apiKey, hc, *reader, nil}
	identity := client.identity()

	ui := UserInput{BufIoWrapper{*reader}}

	//if identity == nil {
	//	setup := Setup{host, hc, *reader}
	//	apiKey = setup.Setup()
	//	if apiKey == "" {
	//		println("Unable to setup invite")
	//		return
	//	}
	//}

	spaces := client.listSpaces()
	var space Space

	if len(spaces) == 0 {
		println("no space found, exiting")
	} else if len(spaces) > 1 {

		names := make([]string, len(spaces))
		for i, item := range spaces {
			names[i] = item.Name
		}

		idx, err := ui.AskWithOptions("Please select a space: ", names, 2)

		if err != nil {
			println(err.Error())
			return
		}
		space = spaces[idx]
	} else {
		space = spaces[0]
	}

	println("Connected to " + space.Name + " as " + identity.Name)

	var cm ColorManager
	cm.init(identity)

	chats := client.listChats(space)
	var chatId string
	var text string

	if len(chats) > 0 {

		chatNames := make([]string, len(chats))
		for i, chat := range chats {
			names := chat.Creator.Name

			for _, participant := range chat.Participants {
				names += ", " + participant.Name
			}

			chatNames[i] = chat.Name + " " + string(names) + " " + chat.Type
		}

		idx, err := ui.AskWithOptions("Please select a chat: ", chatNames, 2)

		if err != nil {
			println(err.Error())
			return
		}
		space = spaces[idx]

		var inputChatId int64
		if _, err := fmt.Sscan(text, &inputChatId); err == nil {
			fmt.Printf("i=%d, type: %T\n", inputChatId, inputChatId)
		}

		chatId = chats[inputChatId-1].Id
	} else {
		println("There are no chats available. Why not start one?")
		chatType, err := ui.AskWithOptions("What kind of chat would you like to start?", []string{"public chat channel", "private chat channel", "private message conversation"}, 2)

		if err != nil {
			println(err.Error())
			return
		}

		users, err := client.getUsersForSpace(space)

		if err != nil {
			println("Error getting participants: " + err.Error())
			return
		}

		userStrings := make([]string, len(users))

		for i, u := range users {
			userStrings[i] = u.User.Name
		}

		idx, err := ui.AskWithOptions("Please invite another user to this chat: ", userStrings, 2)

		ids := []string{identity.Id, users[idx].User.Id}

		if chatType < 3 {
			print("Give your chat a name: ")
			text, _ = reader.ReadString('\n')
		} else {
			text = identity.Name + ", " + users[idx].User.Name + " (private)"
		}

		if text != "" {
			chatId = client.createChat(space, text, chatType, ids).Id
		}
	}

	m := client.previous(chatId, 20)

	printMessages(m, identity, cm)

	w := boolWrapper{true}

	go poll(space, chatId, cm, &client, &w, identity)

	for w.running {
		text, _ = reader.ReadString('\n')

		if text == ":q" {
			w.running = false
		} else {
			client.sendMessage(chatId, text)
		}
	}
}

func printMessages(m []Message, identity User, cm ColorManager) {
	for i := 0; i < len(m); i++ {
		user := m[i].Author
		name := user.Name
		if user == identity {
			name = "(you)"
		}
		fmt.Printf("%s%s: %s%s\n", cm.getColor(user), name, m[i].Text, RESET)
	}
}

func printErr(error string) {
	println(RED + error + RESET)
}

func poll(space Space, chatId string, cm ColorManager, cc *ChatClient, monitor *boolWrapper, identity User) {
	last := time.Now()
	for monitor.running {
		var messages []Message
		messages = cc.since(space, chatId, last.UnixNano()/1000000)
		last = time.Now()

		printMessages(messages, identity, cm)

		time.Sleep(2000000000)

	}
}

type boolWrapper struct {
	running bool
}

type YamlProps struct {
	ApiKey string `yml:"apiKey"`
}
