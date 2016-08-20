package irc

import (
    "fmt"
)

type IrcClient struct {
    server, channel, user, password string
    port int
}

func Connect() {
    fmt.Println("irc connection maaaaddeee broes")
}

func Listen(client IrcClient) {
}

