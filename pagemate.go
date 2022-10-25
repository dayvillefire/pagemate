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

package pagemate

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

// PageMateClient is the base class for accessing a PageMate paging interface
type PageMateClient struct {
	// BaseURL is the base URL under which the web interface is accessed.
	// This should not have a trailing slash
	BaseURL string
	// Username represents the username for the client
	Username string
	// Password represents the password for the lcient
	Password string
	// loggedIn is the internal logged in status of the browser
	loggedIn bool
	// browserObject is the internap surf representation
	browserObject *browser.Browser
}

// NewPageMateClient creates a new client with the specified base URL,
// username, and password.
func NewPageMateClient(baseURL, username, password string) PageMateClient {
	return PageMateClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

// Login logs into the web instance unless it has already happened.
// Returns an error if anything goes wrong.
func (p *PageMateClient) Login() error {
	if p.loggedIn {
		return nil
	}
	b := surf.NewBrowser()
	p.browserObject = b

	b.SetUserAgent(agent.Chrome())

	// Required to not have ASP.NET garbage yak all over me
	b.AddRequestHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	err := b.Open(p.BaseURL + "/")
	if err != nil {
		return err
	}

	if len(b.Forms()) < 1 {
		return errors.New("Form does not exist")
	}
	f, err := b.Form("form")
	if err != nil {
		return err
	}
	f.Input("loginSubscriber", strings.ToUpper(p.Username))
	f.Input("password", p.Password)

	if f.Submit() != nil {
		return err
	}

	err = b.Bookmark("main")
	if err != nil {
		return fmt.Errorf("Unable to bookmark main: %w", err)
	}

	return nil
}

// FindRecipientGroups scrapes recipient groups and descriptions from the list
// kept in PageMate. A blank parameter will list all groups.
func (p *PageMateClient) FindRecipientGroups(query string) (map[string]string, error) {
	r := map[string]string{}

	q := query
	if q == "" {
		q = "*"
	}

	if !p.loggedIn {
		err := p.Login()
		if err != nil {
			return r, err
		}
	}

	p.browserObject.Open(p.BaseURL + "/message/find_recipients.asp")

	f, err := p.browserObject.Form("form")
	if err != nil {
		return r, err
	}
	f.Input("objectname", q)

	if f.Submit() != nil {
		return r, err
	}

	if strings.Contains(p.browserObject.Body(), "Error: Invalid") {
		return r, errors.New("Invalid query")
	}

	rLock := &sync.Mutex{}

	p.browserObject.Dom().Find("table.labels tbody tr").Each(func(_ int, s *goquery.Selection) {
		var k, v string
		s.Find("td:nth-child(1) a").Each(func(_ int, s2 *goquery.Selection) {
			k = strings.TrimSpace(s2.Text())
		})
		s.Find("td:nth-child(3) a").Each(func(_ int, s3 *goquery.Selection) {
			v = strings.TrimSpace(s3.Text())
		})
		if k != "" && v != "" {
			rLock.Lock()
			r[k] = v
			rLock.Unlock()
		}
	})

	return r, nil
}

// SendMessage sends given a message, group, groupDescription, and an optional
// comment.
func (p *PageMateClient) SendMessage(message string, group, groupDescription string, comment string) error {
	if !p.loggedIn {
		err := p.Login()
		if err != nil {
			return err
		}
	}

	p.browserObject.Open(fmt.Sprintf("%s/message/send.asp?objectname=%s", p.BaseURL, url.QueryEscape(group)))

	f, err := p.browserObject.Form("form")
	if err != nil {
		return err
	}
	f.Input("display_objectname", group)
	f.Input("description", groupDescription)
	f.Input("comments", comment)
	f.Input("message", message)

	if f.Submit() != nil {
		return err
	}

	if strings.Contains(p.browserObject.Body(), "Error: Invalid") {
		return errors.New("Invalid query")
	}

	return nil
}
