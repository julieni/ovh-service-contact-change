package main

import "github.com/ovh/go-ovh/ovh"

type contactChangeRequest struct {
	Token string `json:"token"`
}

var (
	flagOvhEP string
	flagOvhAK string
	flagOvhAS string
	flagOvhCK string
)

func apiAccept(request string, token string) error {
	ovhapi, err := ovh.NewClient(
		flagOvhEP,
		flagOvhAK,
		flagOvhAS,
		flagOvhCK,
	)
	if err != nil {
		return err
	}

	params := &contactChangeRequest{
		Token: token,
	}

	err = ovhapi.Post("/me/task/contactChange/"+request+"/accept", params, nil)
	if err != nil {
		return err
	}

	return nil
}
