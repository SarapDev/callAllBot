package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var telegramUrl = os.Getenv("TELEGRAM_BOT_URL")

func mentionAll(admins []Admin, chatId int64) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", telegramUrl+"/sendMessage", nil)

	if err != nil {
		fmt.Println(err)
	}

	c := make(chan string, len(admins))
	wg := sync.WaitGroup{}

	for _, admin := range admins {
		wg.Add(1)
		go mentionUser(admin.User, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var text string
	for i := range c {
		text += i + " "
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.FormatInt(chatId, 10))
	q.Add("text", text)
	q.Add("parse_mode", "MarkdownV2")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
}

func mentionUser(u User, c chan string, wg *sync.WaitGroup) {
	c <- "\\[" + u.Username + "\\](tg://user?id=" + strconv.FormatInt(u.Id, 10) + ")"
	wg.Done()
}

func getAllAdmin(m Message) []Admin {
	client := &http.Client{}
	req, err := http.NewRequest("GET", telegramUrl+"/getChatAdministrators", nil)
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.FormatInt(m.Chat.Id, 10))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	var data AllAdminResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data.Result
}

func GetUpdate() error {
	var data Updates
	resp, err := http.Get(telegramUrl + "/getUpdates")

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	resultLastElem := len(data.Result) - 1
	lastUpdate := data.Result[resultLastElem]

	if lastUpdate.Message.Text == "/all@call_all_users_bot" {
		admins := getAllAdmin(lastUpdate.Message)
		fmt.Println(admins[0])
		mentionAll(admins, lastUpdate.Message.Chat.Id)
	}

	return nil
}
