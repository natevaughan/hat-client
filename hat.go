package main

import (
    "os"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "math/rand"
)

func main() {

    host := os.Getenv("HAT_HOST")
    token := os.Getenv("HAT_TOKEN") 
    if (len(host) == 0) {
        printerr("you must have HAT_HOST set in your environment")
        return
    }
    if (len(token) == 0) {
        printerr("you must have HAT_TOKEN set in your environment")
        return
    }
    client := &http.Client{}
    colorMap := map[User]string {}
 
    body, err := HttpGet(host + "user", token, client)

    if err != nil {
        printerr(err.Error())
    }

    var identity User
    colorMap[identity] = ansi("RESET")

    err = json.Unmarshal(body, &identity)

    if err != nil {
        printerr(err.Error())
    }

    body, err = HttpGet(host + "message/previous/25", token, client)

    if err != nil {
        printerr(err.Error())
    }
    var m []Message

    err = json.Unmarshal(body, &m); 
    if err != nil {
        printerr(err.Error())
    }

    for i := 0; i < len(m); i++ {
        if val, ok := colorMap[m[i].Author]; ok {
            println(val + m[i].Text + ansi("RESET"))
        } else {
            colorMap[m[i].Author] = ansi("RANDOM")
        }

    }
}

func HttpGet(path string, token string, client *http.Client) ([]byte, error) {
    req, err := http.NewRequest("GET", path, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Add("api-token", token)
    
    resp, err := client.Do(req)
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

func printerr(error string) {
    println(ansi("RED") + error + ansi("RESET"))
}

func ansi(name string) string {

    colors := map[string]string{
        "RESET":  "\u001B[0m",
        "RED":    "\u001B[31m",
        "GREEN":  "\u001B[32m",
        "YELLOW": "\u001B[33m",
        "BLUE":   "\u001B[34m",
        "PURPLE": "\u001B[35m",
        "CYAN":   "\u001B[36m",
        "WHITE":  "\u001B[37m",
    }

    if name == "RANDOM" {
        all := []string {"\u001B[37m"}

        for k := range colors {
            all = append(all, colors[k])
        }
        return all[rand.Intn(len(all))]
    }

    return colors[name]
}

type Message struct {
    Id int64
    Timestamp int64
    Text string
    Author User
}

type User struct {
    Name string
    Role string
}

