package telegramAPI

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func SendMessageToUser(userId int, botID int, apiKey string, text string, debugFlag bool) (int, error) {
	if len(text) > 4095 {
		log.Fatal("Message size exceeds 4096 symbols - this is a limitation of the Telegram API")
	}
	main_url := fmt.Sprintf("https://api.telegram.org/bot%v:%v/sendMessage", botID, apiKey)

	resp, err := http.Get(fmt.Sprintf("%s?chat_id=%v&text=%s&parse_mode=Markdown", main_url, userId, url.QueryEscape(text)))
	if err != nil {
		return 503, err
	}

	if resp.StatusCode != http.StatusOK {
		if debugFlag {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return resp.StatusCode, err
			}
			log.Printf("Result: %v", string(body))
		}
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
