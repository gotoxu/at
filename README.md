[![Build Status](https://travis-ci.org/Jamesxql/at.svg?branch=master)](https://travis-ci.org/Jamesxql/at)

# At
A at library for Go, simulate the linux 'at' command. The at command is used to schedule a one-time task at a specific time

## Example
```Go
package main

import (
	"fmt"

	"github.com/Jamesxql/at"
)

func main() {
	at := at.New()
	at.AddFunc("now + 5 minutes", func() {
		fmt.Println("Hello world!")
	})
	at.Start()
	select {}
}
```

## Time definition
| Definition | Example | Description |
| ------| ------ | ------ |
| HH:mm:ss[HH:mm] | 14:30:52, 14:30 | For example, 14:30 specifies 14:30PM, If the time is already past, it is executed at the specified time the next day. |
| HH:mm:ss[HH:mm] yyyy-MM-dd | 14:30 2017-09-01 | If the time is already past, The function will execute immediately |
| hh:mm[AM/PM] Month Date | 04:30pm sep 1 | Currently only support 12-hour clock |
| hh:mm[am/pm]/now + number [seconds/minutes/hours/days/weeks/months] | now + 5 minutes | Currently only support 12-hour clock |
