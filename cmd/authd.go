//	authd.go - Kahlela project
//	Copyright 2023 by Kurt Duncan
//	Main code for authentication service

package main

import (
	"fmt"
	"kalehla/auth"
	"log"
)

func doInitialize(a *auth.Authenticator) {
	fmt.Println("Initializing authd database...")
	err := a.Initialize()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Initialization complete")
}

func main() {
	a := auth.Authenticator{}
	err := a.Open()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func(a *auth.Authenticator) {
		err := a.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(&a)

	doInitialize(&a) //TODO figure out what major function to run
}
