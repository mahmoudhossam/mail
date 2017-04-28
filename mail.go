package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/urfave/cli"
)

// Config user configuration
type Config struct {
	Login  LoginConfig
	Server ServerConfig
}

// ServerConfig the server section of user configuration
type ServerConfig struct {
	Host string
	Port int
}

// LoginConfig the login section of user configuration
type LoginConfig struct {
	Username string
	Password string
}

var configFilePath = "config.toml"

func readConfig(conf *Config) {
	tomlData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Could not find configuration file, exiting...")
		}
		log.Fatal(err)
	}
	if _, err = toml.Decode(string(tomlData), conf); err != nil {
		log.Fatal(err)
	}
}

func listMailboxes(client *client.Client) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

func checkMail(client *client.Client) {
	items := []string{"MESSAGES", "UNSEEN"}
	status, err := client.Status("INBOX", items)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You have %v new emails, %v total.\n", status.Unseen, status.Messages)
}

func connect(config *Config) (c *client.Client) {
	c, err := client.DialTLS(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")
	return
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Mail"
	app.Version = "0.1"
	app.Usage = "Reads email"
	app.UsageText = "mail COMMAND"
	app.Description = "A simple e-mail client."
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Action: func(ctx *cli.Context) {
				var config Config
				readConfig(&config)
				c := connect(&config)
				// Logout when done
				defer c.Logout()

				err := c.Login(config.Login.Username, config.Login.Password)

				if err != nil {
					log.Fatal(err)
				}
				checkMail(c)
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config", Value: "config.toml", Destination: &configFilePath},
	}
	return app
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
