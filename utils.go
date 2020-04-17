package main

import (
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/pkg/browser"
	"golang.org/x/crypto/ssh/terminal"
)

func createAPIToken() error {
	err := browser.OpenURL("https://eu.api.ovh.com/createToken/?GET=/me/task/contactChange&GET=/me/task/contactChange/*&POST=/me/task/contactChange/*")
	return err
}

var extractRegexp *regexp.Regexp

func extractRequestIDAndToken(content []byte) (requestID string, token string) {
	matches := extractRegexp.FindAllSubmatch(content, -1)
	if matches != nil {
		requestID = string(matches[0][1])
		token = strings.TrimSpace(string(matches[0][2]))
	}
	return requestID, token
}

func passwordInput() (string, error) {
	initialState, err := terminal.GetState(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, os.Kill)
	go func() {
		<-channel
		_ = terminal.Restore(int(syscall.Stdin), initialState)
		os.Exit(0)
	}()

	password, err := terminal.ReadPassword(int(syscall.Stdin))
	return string(password), nil
}
