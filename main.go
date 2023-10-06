package main

import (
	"log"

	"telegram-bot-insights/database"
	"telegram-bot-insights/handlers"
	"telegram-bot-insights/telegramAPI"
)

func main() {
	defaultArgs := handlers.Arguments{
		false,                         // html-report
		false,                         // parse
		false,                         // debug
		false,                         // send-msg-to-user
		false,                         // send-msg-to-all-users-of-bot
		false,                         // send-once
		0,                             // msg-id
		0,                             // user-id
		0,                             // bot-id
		"",                            // create-msg
		"./tg-bot-insights.db",        //db-path
		"./bots_history",              // parse-dir
		"2006-01-02T15:04:05",         // custom-date-fmt
		`(?i)(?:Id|chat_id):\s+(\d+)`, // custom-regex
		"BOTS_API_KEYS",               // custom-env-var
	}

	parsedArgs := handlers.ParseArgs(defaultArgs)

	db, err := database.InitDB(parsedArgs.DBPath, parsedArgs.Debug)
	if err != nil {
		log.Fatal("Error Initializing DB: %v\n", err)
	}
	defer db.Close()

	if parsedArgs.HtmlReport {
		handlers.StartHTMLReportServer(db, parsedArgs.MessageID, parsedArgs.BotID)
	}

	if parsedArgs.Parse {
		if parsedArgs.Debug {
			log.Println("Started Parsing")
		}
		parsedJsons, err := handlers.ProcessJSONFilesInDirectory(parsedArgs)
		if err != nil {
			log.Fatal("Error processing JSON files: %v\n", err)
		} else {
			for _, item := range parsedJsons {
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

		if parsedArgs.Debug {
			log.Println("Finished parsing")
			log.Println("Started BotsMapping creation")
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
		if parsedArgs.Debug {
			log.Println("Finished BotsMapping creation")
		}
	}

	if parsedArgs.CreateMessage != "" {
		if parsedArgs.Debug {
			log.Println("Started Message creation")
		}
		err := database.CreateMessageText(db, parsedArgs.CreateMessage)
		if err != nil {
			log.Fatalf("Error creating message with text %v: %v", parsedArgs.CreateMessage, err)
		}
		if parsedArgs.Debug {
			log.Println("Finished Message creation")
		}
	}

	if parsedArgs.SendMessageToUser {
		if parsedArgs.Debug {
			log.Println("Started Sending message to the single user")
		}

		message, err := database.GetMessageTextByTextId(db, parsedArgs.MessageID)
		if err != nil {
			log.Fatalf("Error getting message with ID %v: %v", parsedArgs.MessageID, err)
		}

		bot, err := database.GetBotByID(db, parsedArgs.BotID)
		if err != nil {
			log.Fatalf("Error getting bot with ID %v: %v", parsedArgs.BotID, err)
		}

		res, err := telegramAPI.SendMessageToUser(parsedArgs.UserID, bot.ID, bot.APIKey, message, true)
		if err != nil {
			log.Fatalf("Error sending message to User %v: %v", parsedArgs.UserID, err)
		}

		errLog := database.LogMessage(db, bot.ID, parsedArgs.UserID, parsedArgs.MessageID, res)
		if errLog != nil {
			log.Fatalf("Error logging message %v: %v", parsedArgs.MessageID, errLog)
		}

		if res == 401 || res == 403 {
			err := database.MakeUserInactiveByBotID(db, bot.ID, parsedArgs.UserID)
			if err != nil {
				log.Fatalf("Error setting user %v as inactive: %v", parsedArgs.UserID, err)
			}
		}

		if parsedArgs.Debug {
			log.Println("Finished Sending message to the single user")
		}
	}

	if parsedArgs.SendMessageToAllUsersOfBot {
		if parsedArgs.Debug {
			log.Println("Started Sending message to all the users of the Bot")
		}

		message, err := database.GetMessageTextByTextId(db, parsedArgs.MessageID)
		if err != nil {
			log.Fatalf("Error getting message with ID %v: %v", parsedArgs.MessageID, err)
		}

		bot, err := database.GetBotByID(db, parsedArgs.BotID)
		if err != nil {
			log.Fatalf("Error getting bot with ID %v: %v", parsedArgs.BotID, err)
		}

		users, err := database.GetUsersByBotID(db, bot.ID)
		if err != nil {
			log.Fatalf("Error getting users with botID %v: %v", bot.ID, err)
		}

		step := 10
		for i, user := range users {
			skip := false
			if parsedArgs.Debug {
				pct := handlers.PercentCompleted(i, len(users))
				if pct > 0 && pct%step == 0 {
					log.Printf("%v%% Completed. Processing %v out of %v", pct, i+1, len(users))
					step += 10
				}
			}
			if parsedArgs.SendOnce {
				messageAlreadySent, err := database.MessageWasSentEarlier(db, bot.ID, user.ID, parsedArgs.MessageID)
				if err != nil {
					log.Fatalf("Error checking if message was sent: %v", err)
				}

				if messageAlreadySent {
					skip = true
				}
			}

			userIsActive, err := database.UserIsActive(db, user.ID, bot.ID)
			if err != nil {
				log.Fatalf("Error checking if user is active: %v", err)
			}
			if !skip && userIsActive {
				res, err := telegramAPI.SendMessageToUser(user.ID, bot.ID, bot.APIKey, message, parsedArgs.Debug)
				if err != nil {
					log.Fatalf("Error sending message to User %v: %v", user.ID, err)
				}

				errLog := database.LogMessage(db, bot.ID, user.ID, parsedArgs.MessageID, res)
				if errLog != nil {
					log.Fatalf("Error logging message %v: %v", parsedArgs.MessageID, errLog)
				}

				if res == 401 || res == 403 {
					err := database.MakeUserInactiveByBotID(db, bot.ID, user.ID)
					if err != nil {
						log.Fatalf("Error setting user %v as inactive: %v", parsedArgs.UserID, err)
					}
				}
			}
		}

		if parsedArgs.Debug {
			log.Println("Finished Sending message to all the users of the Bot")
		}
	}
}
