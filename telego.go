package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/iambighead/goutils/logger"
	"github.com/iambighead/goutils/utils"
	"github.com/iambighead/telego/internal/config"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const VERSION = "v0.0.3"

// --------------------------------

var main_logger logger.Logger
var master_config config.MasterConfig

// --------------------------

func init() {
	logger.Init("telego.log", "TELEGO_LOG_LEVEL")
	main_logger = logger.NewLogger("main")

	ex, err := os.Executable()
	if err != nil {
		main_logger.Error("unable to get executable path")
		os.Exit(1)
	}

	{
		var err error
		config_path := filepath.Join(filepath.Dir(ex), "config.yaml")
		master_config, err = config.ReadConfig(config_path)
		if err != nil {
			main_logger.Error(fmt.Sprintf("failed to read config: %v", err))
		}
	}
}

func processFile(file string, token string, chat_id int64) error {
	// Create bot and enable debugging info
	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	// (more on configuration in examples/configuration/main.go)
	bot, err := telego.NewBot(token)
	if err != nil {
		return err
	}

	msg, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Call method sendMessage.
	// Send a message to sender with the same text (echo bot).
	// (https://core.telegram.org/bots/api#sendmessage)
	bot.SendMessage(
		tu.Message(
			tu.ID(chat_id),
			string(msg),
		),
	)

	return nil
}

func cleanupFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

func monitorFolder(folder string, tele config.TeleConfig) {

	// chatID := master_config.TeleConfigs[0].ChatId

	for {
		main_logger.Debug(fmt.Sprintf("checking folder %s", folder))
		filelist, err := utils.ReadFilelist(folder)

		if err != nil {
			main_logger.Error(err.Error())
			continue
		}

		for _, f := range filelist {
			{
				err := processFile(f, tele.Token, tele.ChatId)
				if err != nil {
					main_logger.Error(fmt.Sprintf("failed to processe %s: %s", f, err.Error()))
				} else {
					main_logger.Info(fmt.Sprintf("processed %s", f))
				}

			}
			{
				err := cleanupFile(f)
				if err != nil {
					main_logger.Error(fmt.Sprintf("failed to remove %s: %s", f, err.Error()))
				}
			}
		}

		time.Sleep(2 * time.Duration(time.Second))
	}
}

// --------------------------

func main() {

	main_logger.Info(fmt.Sprintf("Telego started. Version %s", VERSION))

	for _, sender := range master_config.Senders {
		go monitorFolder(sender.Folder, *sender.TeleConfig)
	}

	for {
		time.Sleep(60 * time.Duration(time.Second))
	}
}
