package main

import (
	"db_lab7/API"
	"fmt"
)

func main() {
	API, err := API.InitApi()

	if err != nil {
		fmt.Println(err)
	}

	err = API.Start()

	if err != nil {
		fmt.Println(err)
	}
}
