package main

import (
	"check_nfi_upate/src"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

var config src.Config

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	// 初始化代码
	exists, err := src.CheckFolderExists("github")
	if err != nil {
		log.Fatalf("Error checking github folder: %v", err)
	}

	if !exists {
		err = src.GitClone(config.Github, "github")
		if err != nil {
			log.Fatalf("Error cloning github: %v", err)
		}
	} else {
		err = src.GitPull("github")
		if err != nil {
			log.Fatalf("Error pulling github: %v", err)
		}
	}

	// 初始化 Telegram Bot
	bot, err := tgbotapi.NewBotAPI(config.TgBot)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	fmt.Println("Message sent successfully to Telegram user.")

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/show_config"),
			tgbotapi.NewKeyboardButton("/sync_config"),
			tgbotapi.NewKeyboardButton("/start"),
			tgbotapi.NewKeyboardButton("/stop"),
		),
	)

	// 发送消息并附带Reply Keyboard
	msg := tgbotapi.NewMessage(config.TgUserID, "update_nfi_bot init successful!")
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// 设置更新监听
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updatesTgMsg := bot.GetUpdatesChan(u)
		for tgMsg := range updatesTgMsg {
			if tgMsg.Message != nil {
				switch tgMsg.Message.Text {
				case "/show_config":
					src.ShowConfig(bot, config.TgUserID, config)
				case "/start":
					config.Status = true
					src.Stop(bot, config.TgUserID)
				case "/stop":
					config.Status = false
					src.Start(bot, config.TgUserID)
				case "/sync_config":
					src.SyncConfig(bot, config.TgUserID)
					task(bot, true)
					src.ShowConfig(bot, config.TgUserID, config)
				default:
					msgDefault := tgbotapi.NewMessage(config.TgUserID, "Unknown command. Try /show_config or /sync_config.")
					bot.Send(msgDefault)
				}
			}
		}
	}()

	// 创建一个 15 分钟的 ticker
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	// 立即执行一次任务
	task(bot, config.Status)
	// 使用 for 循环不断监听 ticker
	for {
		select {
		case <-ticker.C:
			// 每 15 分钟触发一次
			task(bot, config.Status)
		}
	}
}

func task(bot *tgbotapi.BotAPI, canChange bool) {
	hasChange := false
	for i, info := range config.UpdateInfos {
		isChangedFlag := false
		for j, file := range info.Files {
			err, isChange := src.CheckAndReplaceFile(file.FilePath, file.GithubFilePath, canChange)
			if err != nil {
				fmt.Printf("Error updating %v file %v: %v", info.Name, file.FilePath, err)
				continue
			}
			if isChange {
				info.Files[j].Update = true
				config.UpdateInfos[i] = info
				isChangedFlag = true
			} else {
				info.Files[j].Update = false
				config.UpdateInfos[i] = info
			}
		}
		// 改变了， call 一下重启接口
		if isChangedFlag {
			hasChange = true
			if canChange {
				err := src.RestartBot(info.RestartAPI, config.Username, config.Password)
				if err != nil {
					fmt.Printf("Error restarting bot: %v", err)
					continue
				}
				msg := tgbotapi.NewMessage(config.TgUserID, fmt.Sprintf("success update %v!", info.Name))
				_, err = bot.Send(msg)
			}
		}
	}

	// 检查一下本地的版本号， 并更新
	for i, info := range config.UpdateInfos {
		for _, file := range info.Files {
			if file.HasVersion {
				versionNum := src.GetVersions(file.FilePath)
				config.UpdateInfos[i].Version = versionNum
				versionGithubNum := src.GetVersions(file.GithubFilePath)
				config.UpdateInfos[i].GitVersion = versionGithubNum
			}
		}
	}

	if hasChange {
		src.ShowConfig(bot, config.TgUserID, config)
	}
}
