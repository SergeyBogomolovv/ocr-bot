package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

var (
	userLangs sync.Map
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	if err := checkTesseractAvailable(); err != nil {
		log.Fatalf("Tesseract check failed: %v", err)
	}

	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatalf("failed to create new bot: %v", err)
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update", "error", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	dispatcher.AddHandler(handlers.NewCommand("start", handleStart))
	dispatcher.AddHandler(handlers.NewCommand("lang", handleLang))

	dispatcher.AddHandler(handlers.NewMessage(message.Photo, handlePhoto))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("lang:"), handleLanguageSelection))

	updater := ext.NewUpdater(dispatcher, nil)
	updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:     9,
			RequestOpts: &gotgbot.RequestOpts{Timeout: time.Second * 10},
		},
	})

	log.Printf("%s has been started...\n", bot.User.Username)
	<-ctx.Done()
	updater.Stop()
	log.Println("bot stopped")
}

func checkTesseractAvailable() error {
	path, err := exec.LookPath("tesseract")
	if err != nil {
		return fmt.Errorf("tesseract not found in PATH")
	}
	cmd := exec.Command(path, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute tesseract: %v, output: %s", err, out.String())
	}
	return nil
}

func init() {
	godotenv.Load()
}
