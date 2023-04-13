package main

import (
	"log"
	"os"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
)

func logEvent(event deltachat.Event) {
	switch ev := event.(type) {
	case deltachat.EventInfo:
		log.Printf("INFO: %v", ev.Msg)
	case deltachat.EventWarning:
		log.Printf("WARNING: %v", ev.Msg)
	case deltachat.EventError:
		log.Printf("ERROR: %v", ev.Msg)
	}
}

func runEchoBot(bot *deltachat.Bot) {
	sysinfo, _ := bot.Account.Manager.SystemInfo()
	log.Println("Running deltachat core", sysinfo["deltachat_core_version"])

	bot.On(deltachat.EventInfo{}, logEvent)
	bot.On(deltachat.EventWarning{}, logEvent)
	bot.On(deltachat.EventError{}, logEvent)
	bot.OnNewMsg(func(msg *deltachat.Message) {
		snapshot, _ := msg.Snapshot()
		chat := &deltachat.Chat{bot.Account, snapshot.ChatId}
		chat.SendText(snapshot.Text)
	})

	if !bot.IsConfigured() {
		log.Println("Bot not configured, configuring...")
		err := bot.Configure(os.Args[1], os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
	}

	addr, _ := bot.GetConfig("configured_addr")
	log.Println("Listening at:", addr)
	bot.Run()
}

func main() {
	rpc := deltachat.NewRpcIO()
	rpc.Start()
	defer rpc.Stop()
	manager := &deltachat.AccountManager{rpc}
	runEchoBot(deltachat.NewBotFromAccountManager(manager))
}
