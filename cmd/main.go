package main

import (
	"os"
	"github.com/mahmoudhossam/mail"
)

func main() {
	app := mail.makeApp()
	app.Run(os.Args)
}
