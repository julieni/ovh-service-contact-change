package main

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {

	dotEnvError := godotenv.Load()
	if dotEnvError != nil {
	}

	extractRegexp, _ = regexp.Compile("contacts/([0-9]+)\\?.+token=(.+)")

	app := cli.NewApp()
	app.Name = "OVH Service contact change"
	app.Author = "Julien Issler"
	app.Email = "julien@issler.net"
	app.Version = "0.2.0"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "ovh-endpoint",
			Value:       "ovh-eu",
			Usage:       "OVH API endpoint",
			EnvVar:      "OVH_ENDPOINT",
			Destination: &flagOvhEP,
		},
		cli.StringFlag{
			Name:        "ovh-ak",
			Value:       "",
			Usage:       "OVH API Application Key",
			EnvVar:      "OVH_AK",
			Destination: &flagOvhAK,
		},
		cli.StringFlag{
			Name:        "ovh-as",
			Value:       "",
			Usage:       "OVH API Application Secret",
			EnvVar:      "OVH_AS",
			Destination: &flagOvhAS,
		},
		cli.StringFlag{
			Name:        "ovh-ck",
			Value:       "",
			Usage:       "OVH API Consumer Key",
			EnvVar:      "OVH_CK",
			Destination: &flagOvhCK,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "init",
			Action: func(c *cli.Context) error {
				err := createAPIToken()
				return err
			},
		},
		{
			Name: "file",
			Action: func(c *cli.Context) error {
				err := fileAccept(c)
				return err
			},
			Flags: append(flags,
				cli.StringFlag{
					Name:   "dir",
					Value:  "mails",
					Usage:  "directory where mails are stored as text files, relative to current directory",
					EnvVar: "MAIL_DIR",
				},
			),
		},
		{
			Name: "imap",
			Action: func(c *cli.Context) error {
				err := imapAccept(c)
				return err
			},
			Flags: append(flags,
				cli.StringFlag{
					Name:   "imap-server",
					Value:  "",
					Usage:  "imap server",
					EnvVar: "IMAP_SERVER",
				},
				cli.StringFlag{
					Name:   "imap-port",
					Value:  "993",
					Usage:  "imap port",
					EnvVar: "IMAP_PORT",
				},
				cli.BoolFlag{
					Name:   "imap-no-tls",
					Usage:  "imap do not use TLS",
					EnvVar: "IMAP_NO_TLS",
				},
				cli.StringFlag{
					Name:   "imap-login",
					Value:  "",
					Usage:  "imap login",
					EnvVar: "IMAP_LOGIN",
				},
				cli.UintFlag{
					Name:  "since-days",
					Value: 1,
					Usage: "Only consider emails received in the last N days",
				},
			),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
