package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ParsedUsers struct {
	ParsedUsers []ParsedUser
}

type ParsedUser struct {
	BotID      int
	ChatIDs    []string
	Usernames  []string
	Timestamps []time.Time
}

func parseJSONFile(filename string, dateFormat string, regex string) (int, []string, []string, []time.Time, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, []string{}, []string{}, []time.Time{}, err
	}

	var messageArray MessageArray
	if err := json.Unmarshal(data, &messageArray); err != nil {
		return 0, []string{}, []string{}, []time.Time{}, err
	}

	botID, chatIDs, usernames, timestamps, err := messageArray.ExtractData(dateFormat, regex)
	if err != nil {
		return 0, []string{}, []string{}, []time.Time{}, err
	}

	return botID, chatIDs, usernames, timestamps, nil
}

func ProcessJSONFilesInDirectory(directoryPath string, dateFormat string, regex string) (ParsedUsers, error) {
	ParsedUsers := ParsedUsers{}
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			fmt.Printf("Processing %v\n", path)
			if botID, chatIDs, usernames, timestamps, err := parseJSONFile(path, dateFormat, regex); err != nil {
				fmt.Printf("Error parsing JSON file %s: %v\n", path, err)
			} else {
				ParsedUser := ParsedUser{
					BotID:      botID,
					ChatIDs:    chatIDs,
					Usernames:  usernames,
					Timestamps: timestamps,
				}

				ParsedUsers.ParsedUsers = append(ParsedUsers.ParsedUsers, ParsedUser)
			}
		}
		return nil
	})

	if err != nil {
		return ParsedUsers, err
	}

	return ParsedUsers, nil
}
