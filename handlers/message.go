package handlers

import (
	"fmt"
	"regexp"
	"time"
)

type MessageArray struct {
	BotID      int       `json:"id"`
	Messages   []Message `json:"messages"`
	DateFormat string    `json:"-"`
}

type Message struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	Date         string `json:"date"`
	From         string `json:"from"`
	FromID       string `json:"from_id"`
	TextEntities []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"text_entities"`
}

func (ma *MessageArray) ExtractData(dateFormat string, regex string) (int, []string, []string, []time.Time, error) {
	var chatIDs []string
	var usernames []string
	var timestamps []time.Time

	for i, message := range ma.Messages {
		fmt.Printf("Processing message #%v\n", i)
		if message.Type == "message" {
			chatID, username, timestamp, err := extractChatIDUsernameAndTimestampFromMessage(message, dateFormat, regex)
			if err == nil {
				if !containsValue(chatIDs, chatID) {
					chatIDs = append(chatIDs, chatID)
					usernames = append(usernames, username)
					timestamps = append(timestamps, timestamp)
				}
			}
		}
	}

	return ma.BotID, chatIDs, usernames, timestamps, nil
}

func extractChatIDUsernameAndTimestampFromMessage(message Message, dateFormat string, regex string) (string, string, time.Time, error) {
	var chatID string
	var username string
	var timestamp time.Time

	for _, entity := range message.TextEntities {
		switch entity.Type {
		case "plain":
			re := regexp.MustCompile(regex)
			match := re.FindStringSubmatch(entity.Text)

			if len(match) >= 2 {
				chatID = match[1]
			}
		case "mention":
			username = entity.Text
		case "phone":
			chatID = entity.Text
		}
	}

	if chatID == "" {
		return "", "", timestamp, fmt.Errorf("chat ID not found in message: %v", message)
	}

	parsedTimestamp, err := time.Parse(dateFormat, message.Date)
	if err != nil {
		return "", "", timestamp, err
	}

	return chatID, username, parsedTimestamp, nil
}

func containsValue(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
