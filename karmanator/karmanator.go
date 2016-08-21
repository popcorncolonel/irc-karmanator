package main

import (
    "fmt"
    "github.com/karmanator/ircutils"
    "github.com/thoj/go-ircevent"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "regexp"
    "sort"
    "strings"
)

type IrcSettings struct {
    Nick     string
    User     string
    Password string
    Server   string
    Channels []string
    Port     int
}
type KarmaObj struct {
    name     string
    rankings map[string]int
}
type KarmaList []KarmaObj

func (k KarmaList) Len() int {
    return len(k)
}

func (k KarmaList) Swap(i, j int) {
    k[i], k[j] = k[j], k[i]
}

func (k KarmaList) Less(i, j int) bool {
    totalKarmaI := k[i].rankings["++"] - k[i].rankings["--"]
    totalKarmaJ := k[j].rankings["++"] - k[j].rankings["--"]
    // > because we want greatest to be first
    return totalKarmaI > totalKarmaJ
}

func toKarmaList(karmaMap map[string]map[string]int) KarmaList {
    list := make(KarmaList, len(karmaMap), len(karmaMap))
    i := 0
    for key, val := range karmaMap {
        list[i] = KarmaObj{key, val}
        i++
    }
    return list
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

func writeKarmaMapToFile(karmaMap map[string]map[string]int) {
    contents, _ := yaml.Marshal(&karmaMap)
    _ = ioutil.WriteFile("karma.yaml", contents, 0666)
}

func getKarmaMap() map[string]map[string]int {
    m := make(map[string]map[string]int)
    data, _ := ioutil.ReadFile("karma.yaml")
    err := yaml.Unmarshal(data, &m)
    if err != nil {
        fmt.Printf("error un-yaml-ing karma.yaml!!: %+v", err)
        return nil
    } else {
        return m
    }
}

func displayTopKarma(con *irc.Connection, channel string) {
    // Get karma from dict
    karmaMap := getKarmaMap()
    karmaList := toKarmaList(karmaMap)
    if len(karmaList) < 3 {
        return
    }
    sort.Sort(karmaList)
    msg := fmt.Sprintf("Top 3 karma: %v (%v), %v (%v), %v (%v)",
        karmaList[0].name, karmaList[0].rankings["++"]-karmaList[0].rankings["--"],
        karmaList[1].name, karmaList[1].rankings["++"]-karmaList[1].rankings["--"],
        karmaList[2].name, karmaList[2].rankings["++"]-karmaList[2].rankings["--"],
    )
    con.Privmsg(channel, msg)
}

func displayKarma(con *irc.Connection, name string, channel string) {
    // Get karma from dict
    karmaMap := getKarmaMap()
    name = strings.ToLower(name)

    // Send privmsg to channel
    msg := fmt.Sprintf(
        "Karma for %v: %v (++: %v | --: %v | +-: %v)",
        name,
        karmaMap[name]["++"]-karmaMap[name]["--"],
        karmaMap[name]["++"],
        karmaMap[name]["--"],
        karmaMap[name]["+-"],
    )
    con.Privmsg(channel, msg)
}

func addKarma(name string, karmaType string) {
    name = strings.ToLower(name)
    // karmaType is either "++", "--", or "+-"
    karmaMap := getKarmaMap()
    if len(karmaMap[name]) == 0 {
        karmaMap[name] = make(map[string]int)
    }
    karmaMap[name][karmaType] += 1
    writeKarmaMapToFile(karmaMap)
}

func getCallback(con *irc.Connection) func(*irc.Event) {
    return func(e *irc.Event) {
        msg := e.Message()
        channel := e.Arguments[0]
        sender := e.Nick

        words := strings.Split(msg, " ")
        matched, _ := regexp.MatchString("^!karma [a-zA-Z0-9]+$", msg)
        if matched {
            displayKarma(con, words[1], channel)
            return
        }
        matched, _ = regexp.MatchString("^!topkarma$", msg)
        if matched {
            displayTopKarma(con, channel)
            return
        }
        for _, word := range words {
            plusPlusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+\\+\\+$", word)
            if plusPlusMatched {
                nameToAward := word[:len(word)-2]
                addKarma(nameToAward, "++")
                con.Privmsg(channel, sender+": Gave ++ to "+nameToAward)
            }
            plusMinusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+\\+-$", word)
            if plusMinusMatched {
                nameToAward := word[:len(word)-2]
                addKarma(nameToAward, "+-")
                con.Privmsg(channel, sender+": Gave +- to "+nameToAward)
            }
            minusMinusMatched, _ := regexp.MatchString("^[a-zA-Z0-9]+--$", word)
            if minusMinusMatched {
                nameToAward := word[:len(word)-2]
                addKarma(nameToAward, "--")
                con.Privmsg(channel, sender+": Gave -- to "+nameToAward)
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
        settings.Channels,
    )
    con.AddCallback("PRIVMSG", getCallback(con))

    con.Loop()
}
