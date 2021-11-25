package pagemate

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

type PageMateClient struct {
	BaseURL       string
	Username      string
	Password      string
	loggedIn      bool
	browserObject *browser.Browser
}

func NewPageMateClient(baseURL, username, password string) PageMateClient {
	return PageMateClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

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
	f.Input("loginSubscriber", p.Username)
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

func (p *PageMateClient) FindRecipientGroups(query string) (map[string]string, error) {
	r := map[string]string{}

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
	f.Input("objectname", query)

	if f.Submit() != nil {
		return r, err
	}

	if strings.Contains(p.browserObject.Body(), "Error: Invalid") {
		return r, errors.New("Invalid query")
	}

	p.browserObject.Dom().Find("table.labels tbody tr").Each(func(_ int, s *goquery.Selection) {
		var k, v string
		s.Find("td:nth-child(1) a").Each(func(_ int, s2 *goquery.Selection) {
			k = strings.TrimSpace(s2.Text())
		})
		s.Find("td:nth-child(3) a").Each(func(_ int, s3 *goquery.Selection) {
			v = strings.TrimSpace(s3.Text())
		})
		if k != "" && v != "" {
			r[k] = v
		}
	})

	return r, nil
}

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
