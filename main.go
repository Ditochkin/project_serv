package main

import (
	"crypto/sha1"
	"db_lab7/API"
	"fmt"
)

func main() {
	hash := sha1.New()
	hash.Write([]byte("1"))

	pas := fmt.Sprintf("%x", hash.Sum([]byte("asjhdjahsdjahsdas")))

	fmt.Println(pas)

	API, err := API.InitApi()

	if err != nil {
		fmt.Println(err)
	}

	err = API.Start()

	if err != nil {
		fmt.Println(err)
	}
}
