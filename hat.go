package main

import (
    "os"
    "net/http"
	"fmt"
	"bufio"
	"time"
)

func main() {

    host := os.Getenv("HAT_HOST")
    token := os.Getenv("HAT_TOKEN")

    if len(host) == 0 {
        printErr("you must have HAT_HOST set in your environment")
        return
    }

    if len(token) == 0 {
        printErr("you must have HAT_TOKEN set in your environment")
        return
    }

	client := ChatClient{ host, token, http.Client{} }
	identity := client.identity()

	println("Connected to chat server as " + identity.Name)

    var cm ColorManager
    cm.init(identity)
    hats := client.listHats()

	println("Please select a chat:" + identity.Name)

	for i := 0; i < len(hats); i++ {
		names := []string { hats[i].Creator.Name }

		for j := 0; j < len(hats[i].Participants); j++ {
			names = append(names, hats[i].Participants[j].Name)
		}

		fmt.Printf("%d) %s %s\n",i + 1, hats[i].Name, names)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select a hat: ")
	text, _ := reader.ReadString('\n')

	var inputHatId int64
	if _, err := fmt.Sscan(text, &inputHatId); err == nil {
		fmt.Printf("i=%d, type: %T\n", inputHatId, inputHatId)
	}

	hatId := hats[inputHatId - 1].Id

	m := client.previous(hatId, 20)

	printMessages(m, identity, cm)

	w := boolWrapper{ true }

	go poll(hatId, cm, &client, &w, identity)

	for w.running {
		text, _ = reader.ReadString('\n')

		if text == ":q" {
			w.running = false
		} else {
			client.sendMessage(hatId, text)
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

func poll(hatId int64, cm ColorManager, cc *ChatClient, monitor *boolWrapper, identity User) {
	last := time.Now()
	for monitor.running {
		var messages []Message
		messages = cc.since(hatId, last.UnixNano() / 1000000)
		last = time.Now()

		printMessages(messages, identity, cm)

		time.Sleep(2000000000)

	}
}

type boolWrapper struct {
	running bool
}