package main

import (
	"os"
)

func main() {
	app := MakeApp()
	app.Run(os.Args)
}
