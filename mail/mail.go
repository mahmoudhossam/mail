package mail

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

var ConfigFilePath = "config.toml"

func ReadConfig(conf *Config) {
	tomlData, err := ioutil.ReadFile(ConfigFilePath)
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

func ListMailboxes(client *client.Client) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.List("", "*", mailboxes)
	}()

	fmt.Println("Mailboxes:")
	for m := range mailboxes {
		fmt.Println("  * " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

func CheckMail(client *client.Client) {
	items := []imap.StatusItem{"MESSAGES", "UNSEEN"}
	status, err := client.Status("INBOX", items)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You have %v new emails, %v total.\n", status.Unseen, status.Messages)
}

func Connect(config *Config) (c *client.Client) {
	c, err := client.DialTLS(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}
