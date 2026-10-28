package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// The tests send real mail to Mailpit from docker-compose.yml and read it back through its API.
const mailpit = "http://127.0.0.1:8025/api/v1"

type Mail struct {
	Subject string
	Text    string
	HTML    string
}

func ClearMail() error {
	request, err := http.NewRequest(http.MethodDelete, mailpit+"/messages", nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	return response.Body.Close()
}

func MailTo(address string) ([]Mail, error) {
	var found struct {
		Messages []struct {
			ID string
		}
	}
	if err := getJSON(mailpit+"/search?query="+url.QueryEscape(fmt.Sprintf("to:%q", address)), &found); err != nil {
		return nil, err
	}

	mails := make([]Mail, 0, len(found.Messages))
	for _, message := range found.Messages {
		var mail Mail
		if err := getJSON(mailpit+"/message/"+message.ID, &mail); err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}

	return mails, nil
}

func getJSON(address string, into any) error {
	response, err := http.Get(address)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(into)
}
