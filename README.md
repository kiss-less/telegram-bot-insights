# telegram-bot-insights
Telegram Bot Insights is a Go-based tool for parsing Telegram bot logs, tracking message deliveries, and collecting valuable user analytics. Simplify Telegram bot management and enhance user engagement with data-driven insights.

## Features

* Parse JSON files with Telegram history to fecth Chat IDs of the users that interact with the bot
* Send messages to all or the specified users of the Bot
* Track messages delivery
* Track whether users are still active (It is cosidered so if the user allows Bot to send messages to them)
* Run simple Server that generates HTML Report of the specified Message ID and Bot ID
* All the above data is stored in SQLite3 db
