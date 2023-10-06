package handlers

import (
	"flag"
	"log"
)

type Arguments struct {
	Parse                      bool
	Debug                      bool
	SendMessageToUser          bool
	SendMessageToAllUsersOfBot bool
	SendOnce                   bool
	MessageID                  int
	UserID                     int
	BotID                      int
	CreateMessage              string
	ParseDirectory             string
	DateFormat                 string
	Regex                      string
	BotsApiEnvVar              string
}

func ParseArgs(defaultArgs Arguments) Arguments {
	flag.BoolVar(&defaultArgs.Parse, "parse", defaultArgs.Parse, "Parse Json files and write result to DB")
	flag.BoolVar(&defaultArgs.Debug, "debug", defaultArgs.Debug, "Verbose output")
	flag.BoolVar(&defaultArgs.SendMessageToUser, "send-msg-to-user", defaultArgs.SendMessageToUser, "Send a message to the user (Specify msg-id, user-id and bot-id)")
	flag.BoolVar(&defaultArgs.SendMessageToAllUsersOfBot, "send-msg-to-all-users-of-bot", defaultArgs.SendMessageToAllUsersOfBot, "Send a message to all users of bot (Specify msg-id and bot-id)")
	flag.BoolVar(&defaultArgs.SendOnce, "send-once", defaultArgs.SendOnce, "Do not send a message to the user if it was already sent earlier")
	messageIDPtr := flag.Int("msg-id", 0, "ID Of the message to send")
	userIDPtr := flag.Int("user-id", 0, "ID Of the user to send the message to")
	botIDPtr := flag.Int("bot-id", 0, "ID Of the bot")
	createMessagePtr := flag.String("create-msg", defaultArgs.CreateMessage, "Should contain message text if passed")
	parseDirPtr := flag.String("dir", defaultArgs.ParseDirectory, "Should be passed if parse is true")
	customDateFormatPtr := flag.String("custom-date-fmt", defaultArgs.DateFormat, "Custom date format")
	customRegexPtr := flag.String("custom-regex", defaultArgs.Regex, "Custom regex to parse chat_id")
	customEnvVarPtr := flag.String("custom-env-var", defaultArgs.BotsApiEnvVar, "Custom Env Var")
	flag.Parse()

	messageID := *messageIDPtr
	userID := *userIDPtr
	botID := *botIDPtr
	createMessage := *createMessagePtr
	parseDir := *parseDirPtr
	customDateFormat := *customDateFormatPtr
	customRegex := *customRegexPtr
	customEnvVar := *customEnvVarPtr

	if defaultArgs.Parse && parseDir == "" {
		log.Fatal("Please provide --dir argument when setting --parse to true")
	} else if defaultArgs.Parse && parseDir != "" {
		defaultArgs.ParseDirectory = parseDir
	}

	if messageID != 0 {
		defaultArgs.MessageID = messageID
	}

	if userID != 0 {
		defaultArgs.UserID = userID
	}

	if botID != 0 {
		defaultArgs.BotID = botID
	}

	if createMessage != "" {
		defaultArgs.CreateMessage = createMessage
	}

	if customDateFormat != "" {
		defaultArgs.DateFormat = customDateFormat
	}

	if customRegex != "" {
		defaultArgs.Regex = customRegex
	}

	if customEnvVar != "" {
		defaultArgs.BotsApiEnvVar = customEnvVar
	}

	return defaultArgs
}
