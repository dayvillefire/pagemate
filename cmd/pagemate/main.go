package main

import (
	"flag"
	"strings"

	"github.com/dayvillefire/pagemate"
	"github.com/joho/godotenv"
)

var (
	target = flag.String("to", "", "Default destination")
)

func main() {
	flag.Parse()
	message := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if message == "" {
		panic("no message!")
	}
	var env map[string]string
	env, err := godotenv.Read()

	to := strings.TrimSpace(strings.ToUpper(*target))
	if to == "" {
		if _, ok := env["TO"]; !ok {
			panic("no target")
		}
		to = strings.TrimSpace(strings.ToUpper(env["TO"]))
	}

	pm := pagemate.NewPageMateClient(env["URL"], env["USER"], env["PASS"])
	groups, err := pm.FindRecipientGroups(to)
	if err != nil {
		panic(err)
	}
	err = pm.SendMessage(message, to, groups[to], "")
	if err != nil {
		panic(err)
	}
}
