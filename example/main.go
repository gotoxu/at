package main

import (
	"fmt"
	"time"

	"github.com/gotoxu/at"
)

func main() {
	at := at.New()
	at.AddFunc(time.Now().Add(5*time.Minute), func() {
		fmt.Println("Hello world!")
	})

	at.Start()
	select {}
}
