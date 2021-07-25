package main

import (
	"github.com/emersion/go-imap/client"
	"github.com/urfave/cli"
	"github.com/mahmoudhossam/mail"
	"log"
)

func getClient() *client.Client {
	var config mail.Config
	mail.ReadConfig(&config)
	c := mail.Connect(&config)
	err := c.Login(config.Login.Username, config.Login.Password)

	if err != nil {
		log.Fatal(err)
	}
	return c
}

var commands = []cli.Command{
	{
		Name:    "check",
		Aliases: []string{"c"},
		Action: func(ctx *cli.Context) {
			c := getClient()
			mail.CheckMail(c)
			err := c.Logout()
			if err != nil {
				log.Println(err.Error())
			}
		},
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Action: func(ctx *cli.Context) {
			c := getClient()
			mail.ListMailboxes(c)
			err := c.Logout()
			if err != nil {
				log.Println(err.Error())
			}
		},
	},
}

var flags = []cli.Flag{
	cli.StringFlag{Name: "config", Value: "config.toml", Destination: &mail.ConfigFilePath},
}

func MakeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Mail"
	app.Version = "0.1"
	app.Usage = "Reads email"
	app.UsageText = "mail COMMAND"
	app.Description = "A simple e-mail client."
	app.EnableBashCompletion = true
	app.Commands = commands
	app.Flags = flags
	return app
}
