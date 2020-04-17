package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/textproto"
	"time"

	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/urfave/cli"
)

func imapAccept(c *cli.Context) error {
	server := c.String("imap-server")
	port := c.String("imap-port")
	login := c.String("imap-login")

	if server == "" {
		return errors.New("IMAP server not specified")
	}

	if login == "" {
		return errors.New("IMAP login not specified")
	}

	fmt.Println("IMAP password for " + login + " on server " + server + " port " + port + ": ")
	password, err := passwordInput()
	if err != nil {
		return err
	}

	if password == "" {
		return errors.New("IMAP password not specified")
	}

	var client *imapClient.Client
	if c.Bool("imap-no-tls") == false {
		client, err = imapClient.DialTLS(server+":"+port, nil)
	} else {
		client, err = imapClient.Dial(server + ":" + port)
	}

	if err != nil {
		return err
	}

	defer client.Logout()

	if err := client.Login(login, password); err != nil {
		return err
	}

	_, err = client.Select("INBOX", false)
	if err != nil {
		return err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Header = textproto.MIMEHeader{
		"Subject": {"Modification des informations de contact"},
	}
	date := time.Now()
	criteria.Since = date.AddDate(0, 0, int(^(c.Uint("since-days") - 1)))

	ids, err := client.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		section := &imap.BodySectionName{}
		go func() {
			done <- client.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
		}()

		for msg := range messages {
			bodySection := msg.GetBody(section)
			mailMessage, _ := mail.ReadMessage(bodySection)
			body, _ := ioutil.ReadAll(mailMessage.Body)
			id, token := extractRequestIDAndToken(body)
			if id != "" && token != "" {
				err := apiAccept(id, token)
				if err != nil {
					fmt.Println("Error for request #" + id + " " + err.Error())
				} else {
					fmt.Println("Contact change request #" + id + " accepted")
				}
			}
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
