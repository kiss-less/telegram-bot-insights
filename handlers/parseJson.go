package handlers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ParsedJson struct {
	BotID      int
	ChatIDs    []int
	Usernames  []string
	Timestamps []time.Time
}

func parseJSONFile(filename string, dateFormat string, regex string, debug bool) (int, []int, []string, []time.Time, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, []int{}, []string{}, []time.Time{}, err
	}

	var messageArray MessageArray
	if err := json.Unmarshal(data, &messageArray); err != nil {
		return 0, []int{}, []string{}, []time.Time{}, err
	}

	botID, chatIDs, usernames, timestamps, err := messageArray.ExtractData(dateFormat, regex, debug)
	if err != nil {
		return 0, []int{}, []string{}, []time.Time{}, err
	}

	return botID, chatIDs, usernames, timestamps, nil
}

func ProcessJSONFilesInDirectory(args Arguments) ([]ParsedJson, error) {
	directoryPath := args.ParseDirectory
	dateFormat := args.DateFormat
	regex := args.Regex
	debug := args.Debug

	var ParsedJsons []ParsedJson
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			if args.Debug {
				log.Printf("Started Processing %v\n", path)
			}
			if botID, chatIDs, usernames, timestamps, err := parseJSONFile(path, dateFormat, regex, debug); err != nil {
				log.Printf("Error parsing JSON file %s: %v\n", path, err)
			} else {
				ParsedJson := ParsedJson{
					BotID:      botID,
					ChatIDs:    chatIDs,
					Usernames:  usernames,
					Timestamps: timestamps,
				}

				ParsedJsons = append(ParsedJsons, ParsedJson)
			}
		}

		if args.Debug {
			log.Printf("Finished Processing %v\n", path)
		}

		return nil
	})

	if err != nil {
		return ParsedJsons, err
	}

	return ParsedJsons, nil
}
