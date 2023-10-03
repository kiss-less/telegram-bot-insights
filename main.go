package main

import (
	"fmt"

	"telegram-bot-insights/handlers"
)

func main() {
	directoryPath := "../path/to/json/files"
	dateFormat := "2006-01-02T15:04:05"
	regex := `(?i)(?:Id|chat_id):\s+(\d+)`

	users, err := handlers.ProcessJSONFilesInDirectory(directoryPath, dateFormat, regex)
	if err != nil {
		fmt.Printf("Error processing JSON files: %v\n", err)
	}

	fmt.Printf("%v\n", users)
}
