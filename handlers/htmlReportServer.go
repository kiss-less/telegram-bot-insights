package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"telegram-bot-insights/database"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	MessageID                int
	BotID                    int
	MessageText              string
	TotalUsersPerBot         int
	TotalActiveUsersPerBot   int
	TotalSent                int
	TotalDelivered           int
	LastSentTime             string
	LastSentAndDeliveredTime string
}

func StartHTMLReportServer(db *sql.DB, messageID int, botID int) {

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {

		messageText, err := database.GetMessageTextByTextId(db, messageID)
		if err != nil {
			log.Fatalf("Error getting message with ID %v: %v", messageID, err)
		}

		users, err := database.GetUsersByBotID(db, botID)
		if err != nil {
			log.Fatalf("Error getting users with botID %v: %v", botID, err)
		}

		activeUsers, err := database.TotalNumberOfActiveUsersPerBot(db, botID)
		if err != nil {
			log.Fatalf("Error checking if user is active: %v", err)
		}

		messagesSent, err := database.TotalNumberOfMessagesSentPerBot(db, botID, messageID)
		if err != nil {
			log.Fatalf("Error checking if messages were sent: %v", err)
		}

		messagesDelivered, err := database.TotalNumberOfMessagesPerBotByStatus(db, botID, messageID, 200)
		if err != nil {
			log.Fatalf("Error checking if message was delivered: %v", err)
		}

		lastSentAndDeliveredTime, err := database.LastMessageSentAndDeliveredTimePerBot(db, botID, messageID)
		if err != nil {
			log.Fatalf("Error checking last sent time: %v", err)
		}

		lastSentTime, err := database.LastMessageSentTimePerBot(db, botID, messageID)
		if err != nil {
			log.Fatalf("Error checking last sent time: %v", err)
		}

		data := PageData{
			MessageID:                messageID,
			BotID:                    botID,
			MessageText:              messageText,
			TotalUsersPerBot:         len(users),
			TotalActiveUsersPerBot:   activeUsers,
			TotalSent:                messagesSent,
			TotalDelivered:           messagesDelivered,
			LastSentTime:             lastSentTime,
			LastSentAndDeliveredTime: lastSentAndDeliveredTime,
		}

		tmpl := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Report Page</title>
        </head>
        <body>
            <h1>Report for Message ID #{{.MessageID}}, Bot ID #{{.BotID}}</h1>
			Message Text: {{.MessageText}}<br><br>
			Total number of Bot users: {{.TotalUsersPerBot}}<br>
			Total number of Bot active users: {{.TotalUsersPerBot}}<br>
			Total number of messages sent: {{.TotalSent}}<br>
			Total number of messages delivered: {{.TotalDelivered}}<br><br>
			Last message sent at: {{.LastSentTime}}<br>
			Last successful message sent at: {{.LastSentAndDeliveredTime}}<br>
        </body>
        </html>
        `

		t := template.Must(template.New("report").Parse(tmpl))
		t.Execute(w, data)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
