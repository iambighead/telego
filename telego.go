package main

import (
	"fmt"
	"os"

	"github.com/iambighead/goutils/logger"
	"github.com/iambighead/telego/internal/config"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const VERSION = "v0.0.1"

// --------------------------------

var main_logger logger.Logger
var master_config config.MasterConfig

// --------------------------

func init() {
	logger.Init("telego.log", "TELEGO_LOG_LEVEL")
	main_logger = logger.NewLogger("main")

	var err error
	master_config, err = config.ReadConfig("config.yaml")
	if err != nil {
		main_logger.Error(fmt.Sprintf("failed to read config: %v", err))
	}

}

// --------------------------

func main() {

	main_logger.Info(fmt.Sprintf("Telego started. Version %s", VERSION))

	// Create bot and enable debugging info
	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	// (more on configuration in examples/configuration/main.go)
	bot, err := telego.NewBot(
		master_config.TeleConfigs[0].Token,
		telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	chatID := master_config.TeleConfigs[0].ChatId

	// Call method sendMessage.
	// Send a message to sender with the same text (echo bot).
	// (https://core.telegram.org/bots/api#sendmessage)
	bot.SendMessage(
		tu.Message(
			tu.ID(chatID),
			"hello from golang",
		),
	)
}
