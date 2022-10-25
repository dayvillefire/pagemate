// Licensed to Dayville Fire Company under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Dayville Fire Company licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dayvillefire/pagemate"
	"github.com/joho/godotenv"
)

var (
	action = flag.String("action", "send", "Action (send, list)")
	target = flag.String("to", "", "Default destination")
)

func main() {
	flag.Parse()

	switch *action {
	case "list":
		break
	case "send":
		break
	default:
		flag.PrintDefaults()
		return
	}

	var env map[string]string
	env, err := godotenv.Read()

	url := env["URL"]
	if url == "" {
		url = "http://pageme.qvec.org"
	}
	user := env["USER"]
	pass := env["PASS"]

	pm := pagemate.NewPageMateClient(url, user, pass)

	if *action == "list" {
		groups, err := pm.FindRecipientGroups("")
		if err != nil {
			panic(err)
		}
		for k, v := range groups {
			fmt.Printf("%s: %s\n", k, v)
		}
		return
	}

	message := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if message == "" {
		panic("no message!")
	}

	to := strings.TrimSpace(strings.ToUpper(*target))
	if to == "" {
		if _, ok := env["TO"]; !ok {
			panic("no target")
		}
		to = strings.TrimSpace(strings.ToUpper(env["TO"]))
	}

	groups, err := pm.FindRecipientGroups(to)
	if err != nil {
		panic(err)
	}
	err = pm.SendMessage(message, to, groups[to], "")
	if err != nil {
		panic(err)
	}
}
