package main

import (
	"fmt"
	"log"

	"telegram-bot-insights/handlers"
)

func main() {
	defaultArgs := handlers.Arguments{
		true,                          // Parse flag
		false,                         // Debug flag
		"./bots_history",              // ParseDirectory flag Should contain path of the directory with Json files
		"2006-01-02T15:04:05",         // DateFormat flag to parse timestamp
		`(?i)(?:Id|chat_id):\s+(\d+)`, // Regex flag to parse chat_id
		"BOTS_API_KEYS",               // BotsApiEnvVar flag of env var name to map Bots and API Keys. Env Var's value should be in the following format: "<BOT_ID1>:<API_KEY1> <BOT_ID2>:<API_KEY2>..."
	}

	parsedArgs := handlers.ParseArgs(defaultArgs)

	if parsedArgs.Parse {
		parsedJsons, err := handlers.ProcessJSONFilesInDirectory(parsedArgs)
		if err != nil {
			log.Fatal("Error processing JSON files: %v\n", err)
		}
		fmt.Printf("%v", parsedJsons)

		mappings, err := handlers.CreateBotsKeysMapping(parsedArgs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", mappings)
	}
}
