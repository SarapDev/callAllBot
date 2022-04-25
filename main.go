package main

import (
	"fmt"
	"github.com/SarapDev/callAllBot/internal"
	"time"
)

func main() {
	for {
		time.Sleep(5 * time.Second)

		err := internal.GetUpdate()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
