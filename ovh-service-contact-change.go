package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/ovh/go-ovh/ovh"
	"github.com/pkg/browser"
	"github.com/urfave/cli"
)

func createAPIToken() error {
	err := browser.OpenURL("https://eu.api.ovh.com/createToken/?GET=/me/task/contactChange&GET=/me/task/contactChange/*&POST=/me/task/contactChange/*")
	return err
}

// ContactChangeRequest request
type ContactChangeRequest struct {
	Token string `json:"token"`
}

func contactChangeAccept(c *cli.Context) error {
	cwd, _ := os.Getwd()
	ovhEndpoint := c.String("ovh-endpoint")
	ovhAk := c.String("ovh-ak")
	ovhAs := c.String("ovh-as")
	ovhCk := c.String("ovh-ck")
	dir := c.String("dir")

	ovhapi, _ := ovh.NewClient(
		ovhEndpoint,
		ovhAk,
		ovhAs,
		ovhCk,
	)

	files, err := ioutil.ReadDir(cwd + "/" + dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	r, _ := regexp.Compile("contacts/([0-9]+)\\?.+token=(.+)")
	for _, file := range files {
		c, _ := ioutil.ReadFile(cwd + "/" + dir + "/" + file.Name())
		matches := r.FindAllSubmatch(c, -1)
		if matches != nil {
			id := string(matches[0][1])
			token := strings.TrimSpace(string(matches[0][2]))
			params := &ContactChangeRequest{Token: token}
			err := ovhapi.Post("/me/task/contactChange/"+id+"/accept", params, nil)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Contact change request #" + id + " accepted")
			}
		}
	}

	return nil
}

func main() {

	dotEnvError := godotenv.Load()
	if dotEnvError != nil {
	}

	app := cli.NewApp()
	app.Name = "OVH Service contact change"
	app.Author = "Julien Issler"
	app.Email = "julien@issler.net"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name: "init",
			Action: func(c *cli.Context) error {
				err := createAPIToken()
				return err
			},
		},
		{
			Name: "accept",
			Action: func(c *cli.Context) error {
				err := contactChangeAccept(c)
				return err
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "ovh-endpoint",
					Value:  "ovh-eu",
					Usage:  "OVH API endpoint",
					EnvVar: "OVH_ENDPOINT",
				},
				cli.StringFlag{
					Name:   "ovh-ak",
					Value:  "",
					Usage:  "OVH API Application Key",
					EnvVar: "OVH_AK",
				},
				cli.StringFlag{
					Name:   "ovh-as",
					Value:  "",
					Usage:  "OVH API Application Secret",
					EnvVar: "OVH_AS",
				},
				cli.StringFlag{
					Name:   "ovh-ck",
					Value:  "",
					Usage:  "OVH API Consumer Key",
					EnvVar: "OVH_CK",
				},
				cli.StringFlag{
					Name:   "dir",
					Value:  "mails",
					Usage:  "directory where mails are stored as text files, relative to current directory",
					EnvVar: "MAIL_DIR",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
