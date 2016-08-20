package main

import (
    "fmt"
    "github.com/karmanator/ircutils"
    "github.com/thoj/go-ircevent"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "strings"
	"regexp"
)

type IrcSettings struct {
	Nick string
	User string
	Password string
	Server string
	Channel string
	Port int
}

func getSettings() IrcSettings {
	data, err := ioutil.ReadFile("settings.yaml")
    if err != nil {
        fmt.Printf("error reading settings.yaml: %v", err)
    }
	ircObj := IrcSettings{}
    err2 := yaml.Unmarshal(data, &ircObj)
    if err2 != nil {
        fmt.Printf("error un-yaml-ing settings.yaml: %v", err2)
    }

    return ircObj
}

func displayKarma(con *irc.Connection, name string, channel string) {
	// Get karma from dict
	// Send privmsg to channel 
}

func addKarma(name string, karmaType string) {
	// karmaType is either "++", "--", or "+-"

}

func getCallback(con *irc.Connection) func(*irc.Event) {
	return func (e *irc.Event) {
		msg := e.Message()
		channel := e.Arguments[0]
		words := strings.Split(msg, " ")
		matched, _ := regexp.MatchString("^!karma [a-zA-Z0-9]+$", msg) 
		if matched {
			displayKarma(con, words[1], channel)
			return
		}
		for _, word := range words {
			plusPlusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+\\+\\+$", word)
			if plusPlusMatched {
				nameToAward := word[:len(word)-2]
				addKarma(nameToAward, "++")
			}
			plusMinusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+\\+-$", word)
			if plusMinusMatched {
				nameToAward := word[:len(word)-2]
				addKarma(nameToAward, "+-")
			}
			minusMinusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+--$", word)
			if minusMinusMatched {
				nameToAward := word[:len(word)-2]
				addKarma(nameToAward, "--")
			}
		}
	}
}

func main() {
	settings := getSettings()
	con := ircutils.GetConnection(
		settings.Server,
		settings.Port,
		settings.Nick,
		settings.User,
		settings.Password,
		settings.Channel,
	)
    con.AddCallback("PRIVMSG", getCallback(con))

	con.Loop()
}
