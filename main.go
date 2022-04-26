package main

import (
	"fmt"
	"github.com/SarapDev/callAllBot/internal"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return
	}

	telegramUrl := os.Getenv("TG_BOT_URL")

	for {
		time.Sleep(5 * time.Second)

		err := internal.GetUpdate(telegramUrl)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
