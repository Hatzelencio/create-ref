package main

import (
	"github.com/hatzelencio/crete-ref/remote"
	"log"
)

func main() {
	err := remote.ValidateInputs()

	if err != nil {
		log.Fatal(err)
	}

	remote.CreateGitRef()
}
