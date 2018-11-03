package main

import (
	"fmt"
	"log"
	"rockpaperscissor/http"
)

func main() {
	c := http.NewClient("localhost:3000", "asd")
	res, err := c.CreateGame("whatever")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Created:%v\n", res.GameName)

	err = c.Open("whatever")
	if err != nil {
		log.Fatalln(err)
	}
}
