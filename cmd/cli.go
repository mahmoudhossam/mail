package main

import (
	"github.com/urfave/cli"
	"github.com/mahmoudhossam/mail"
	"log"
)

var commands = []cli.Command{
	{
		Name:    "check",
		Aliases: []string{"c"},
		Action: func(ctx *cli.Context) {
			c := mail.GetClient()
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
			c := mail.GetClient()
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
