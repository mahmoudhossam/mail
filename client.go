package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
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

func readConfig(conf *Config) {
	tomlData, err := ioutil.ReadFile("config.toml")
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

func main() {
	var config Config
	readConfig(&config)
	c, err := client.DialTLS(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Logout when done
	defer c.Logout()

	err = c.Login(config.Login.Username, config.Login.Password)

	if err != nil {
		log.Fatal(err)
	}

	listMailboxes(c)
}
