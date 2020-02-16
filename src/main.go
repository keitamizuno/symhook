package main

const configPATH = "../config/config.json"

func main() {

	botClient := initBot(configPATH)
	go initWebhook(botClient)
	initDatafeedEventsService(botClient)
}
