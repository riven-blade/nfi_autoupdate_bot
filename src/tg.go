package src

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

// ShowConfig 显示配置的功能
func ShowConfig(bot *tgbotapi.BotAPI, chatID int64, config Config) {
	configOutInfo := make([]ConfigOutput, 0)
	for _, info := range config.UpdateInfos {
		updateFile := make([]string, 0)
		for _, file := range info.Files {
			if file.Update {
				updateFile = append(updateFile, file.FilePath)
			}
		}
		configOutInfo = append(configOutInfo, ConfigOutput{
			Name:       info.Name,
			Version:    info.Version,
			UpdateFile: updateFile,
		})
	}
	outInfo, _ := yaml.Marshal(&configOutInfo)
	message := string(outInfo)
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

// SyncConfig 同步配置的功能
func SyncConfig(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Syncing..."
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)

}

// Start 开启机器人
func Start(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Started"
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

// Stop 关闭机器人
func Stop(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Stopped"
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}
