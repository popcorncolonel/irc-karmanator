package ircutils

import (
    "fmt"
    "github.com/thoj/go-ircevent"
)

func GetConnection(
    server string, 
    port int, 
    nick, user, password, roomName string,
) *irc.Connection {
    con := irc.IRC(nick, user)
    con.Password = password
    err := con.Connect(fmt.Sprintf("%v:%v", server, port))
    if err != nil {
        fmt.Println("ERROR CONNECTING: %+v", err)
        return nil
    } else {
        fmt.Println("Connection successful.")
    }
    con.AddCallback("001", func (e *irc.Event) {
        fmt.Println("Got welcome message. Connecting...")
        con.Join("#" + roomName)
        fmt.Println("Connected.")
    })
    return con
}
