package main

import (
    "fmt"
    "os"
    "net/http"
    "io/ioutil"
)

func main() {

    host := os.Getenv("HAT_HOST")
    token := os.Getenv("HAT_TOKEN") 
    if (len(host) == 0) {
        fmt.Printf("you must have HAT_HOST set in your environment\n")
        return
    }
    if (len(token) == 0) {
        fmt.Printf("you must have HAT_TOKEN set in your environment\n")
        return
    }
    client := &http.Client{}
 
    req, err := http.NewRequest("GET", host + "message/previous/20", nil)

    req.Header.Add("api-token", token)

    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf(err.Error() + "\n")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    fmt.Printf(string(body))
}
