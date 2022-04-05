package main

import (
	"log"

	"github.com/fabienbellanger/echo-boilerplate/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}
