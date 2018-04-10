[![Build Status](https://travis-ci.org/gotoxu/at.svg?branch=master)](https://travis-ci.org/gotoxu/at)

# At
A at library for Go, simulate the linux 'at' command. The at command is used to schedule a one-time task at a specific time

## Example
```Go
package main

import (
	"fmt"
	"time"

	"github.com/Jamesxql/at"
)

func main() {
	at := at.New()
	at.AddFunc(time.Now().Add(5*time.Minute), func() {
		fmt.Println("Hello world!")
	})

	at.Start()
	select {}
}
```