package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

func fileAccept(c *cli.Context) error {
	cwd, _ := os.Getwd()
	dir := c.String("dir")

	files, err := ioutil.ReadDir(cwd + "/" + dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		c, _ := ioutil.ReadFile(cwd + "/" + dir + "/" + file.Name())
		id, token := extractRequestIDAndToken(c)
		if id != "" && token != "" {
			err := apiAccept(id, token)
			if err != nil {
				fmt.Println("Error for request #" + id + " " + err.Error())
			} else {
				fmt.Println("Contact change request #" + id + " accepted")
			}
		}
	}

	return nil
}
