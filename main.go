package main

import (
	"log"

	"telegram-bot-insights/database"
	"telegram-bot-insights/handlers"
)

func main() {
	defaultArgs := handlers.Arguments{
		false,                         // Parse flag
		false,                         // Debug flag
		"./bots_history",              // ParseDirectory flag Should contain path of the directory with Json files
		"2006-01-02T15:04:05",         // DateFormat flag to parse timestamp
		`(?i)(?:Id|chat_id):\s+(\d+)`, // Regex flag to parse chat_id
		"BOTS_API_KEYS",               // BotsApiEnvVar flag of env var name to map Bots and API Keys. Env Var's value should be in the following format: "<BOT_ID1>:<API_KEY1> <BOT_ID2>:<API_KEY2>..."
	}

	parsedArgs := handlers.ParseArgs(defaultArgs)

	if parsedArgs.Parse {
		db, err := database.InitDB("./tg-bot-insights.db", parsedArgs.Debug)
		if err != nil {
			log.Fatal("Error Initializing DB: %v\n", err)
		}
		defer db.Close()

		parsedJsons, err := handlers.ProcessJSONFilesInDirectory(parsedArgs)
		if err != nil {
			log.Fatal("Error processing JSON files: %v\n", err)
		} else {
			for _, item := range parsedJsons.ParsedJsons {
				for i, id := range item.ChatIDs {
					userExists, err := database.UserExists(db, id)
					if err != nil {
						log.Fatalf("Error checking if user exists: %v", err)
					}

					errAssoc := database.AssociateBotWithUser(db, item.BotID, id)
					if errAssoc != nil && parsedArgs.Debug {
						log.Println(errAssoc)
					}
					if !userExists {
						err := database.CreateUser(db, id, item.Usernames[i], item.Timestamps[i].Format(parsedArgs.DateFormat))
						if err != nil {
							log.Fatalf("Error creating user with Id %v: %v", id, err)
						}
					}
				}
			}
		}

		mappings, err := handlers.CreateBotsKeysMapping(parsedArgs)
		if err != nil {
			log.Printf("There was an error during Bots API Keys mapping: %v\n", err)
			log.Println("Skipping this step")
		} else {
			for _, item := range mappings {
				botExists, err := database.BotExists(db, item.ID)
				if err != nil {
					log.Fatalf("Error checking if bot exists: %v", err)
				}

				if !botExists {
					err := database.CreateBot(db, item.ID, item.APIKey, "", "")
					if err != nil {
						log.Fatalf("Error creating bot with Id %v: %v", item.ID, err)
					}
				}
			}
		}
	}
}
