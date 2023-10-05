package handlers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

type MessageArray struct {
	BotID    int       `json:"id"`
	Messages []Message `json:"messages"`
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

func (ma *MessageArray) ExtractData(dateFormat string, regex string, debug bool) (int, []int, []string, []time.Time, error) {
	var chatIDs []int
	var usernames []string
	var timestamps []time.Time

	if debug {
		log.Printf("Total number of messages to process: %v", len(ma.Messages))
	}
	step := 10
	for i, message := range ma.Messages {
		if debug {
			pct := percentCompleted(i, len(ma.Messages))
			if pct > 0 && pct%step == 0 {
				log.Printf("%v%% Completed. Processing %v out of %v", pct, i+1, len(ma.Messages))
				step += 10
			}
		}
		if message.Type == "message" {
			chatIDStr, username, timestamp, err := extractChatIDUsernameAndTimestampFromMessage(message, dateFormat, regex)
			if err == nil {
				chatID, err := strconv.Atoi(chatIDStr)
				if err != nil {
					log.Fatal("Error converting ChatId %v to Int: %v", chatIDStr, err)
				}
				if !containsValue(chatIDs, chatID) {
					chatIDs = append(chatIDs, chatID)
					usernames = append(usernames, username)
					timestamps = append(timestamps, timestamp)
				}
			}
		}
	}

	if debug {
		log.Print("100% Completed")
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

func containsValue(slice []int, target int) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func percentCompleted(current int, total int) int {
	return int(float64(current) / float64(total) * 100)
}
