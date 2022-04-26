package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var offset int64

func mentionAll(admins []Admin, chatId int64, telegramUrl string) {
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

	query := make(map[string]string)
	query["chat_id"] = strconv.FormatInt(chatId, 10)
	query["text"] = text
	query["parse_mode"] = "HTML"

	resp := makeRequest("POST", telegramUrl+"/sendMessage", query)

	defer resp.Body.Close()

	var data interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
}

func mentionUser(u User, c chan string, wg *sync.WaitGroup) {
	//<a href="tg://user?id=123456789">inline mention of a user</a>
	c <- "<a href=\"tg://user?id=" + strconv.FormatInt(u.Id, 10) + "\">" + u.Username + "</a>"

	wg.Done()
}

func getAllAdmin(m Message, telegramUrl string) []Admin {
	query := make(map[string]string)
	query["chat_id"] = strconv.FormatInt(m.Chat.Id, 10)

	resp := makeRequest("GET", telegramUrl+"/getChatAdministrators", query)

	defer resp.Body.Close()

	var data AllAdminResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	return data.Result
}

func sendJoke(chatId int64, telegramUrl string) {
	text := "Мама собирает сыну обед в школу:\n— Вот, положила тебе в ранец хлеб, колбасу и гвозди.\n— Мам, нафига??\n— Ну как же, берешь хлеб, кладешь на него колбасу и ешь.\n— А гвозди?\n— Так вот же они!"

	query := make(map[string]string)
	query["chat_id"] = strconv.FormatInt(chatId, 10)
	query["text"] = text
	query["parse_mode"] = "HTML"

	resp := makeRequest("POST", telegramUrl+"/sendMessage", query)

	defer resp.Body.Close()

	var data interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
	}
}

func GetUpdate(telegramUrl string) error {
	var data Updates
	query := make(map[string]string)
	if offset != 0 {
		query["offset"] = strconv.FormatInt(offset, 10)
	}

	resp := makeRequest("GET", telegramUrl+"/getUpdates", query)

	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	resultLastElem := len(data.Result) - 1
	if resultLastElem == -1 {
		return nil
	}

	lastUpdate := data.Result[resultLastElem]
	offset = lastUpdate.Id + 1

	if lastUpdate.Message.Text == "/all@call_all_users_bot" {
		admins := getAllAdmin(lastUpdate.Message, telegramUrl)
		mentionAll(admins, lastUpdate.Message.Chat.Id, telegramUrl)
	}

	if lastUpdate.Message.Text == "/joke@call_all_users_bot" {
		sendJoke(lastUpdate.Message.Chat.Id, telegramUrl)
	}

	return nil
}

func makeRequest(method string, url string, query map[string]string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
