package main

import (
    "fmt"

    "github.com/karmanator/irc"
)

func main() {
    irc.Connect()
    fmt.Println("HEY! DONE!")
}
