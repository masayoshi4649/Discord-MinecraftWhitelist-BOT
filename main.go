package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	conf := getConf()

	dg, err := discordgo.New("Bot " + conf.Bot.Token)
	if err != nil {
		fmt.Println("error:start\n", err)
		return
	}

	dg.UpdateGameStatus(0, "!mc をプレイ中")

	//on message
	dg.AddHandler(conf.messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error:wss\n", err)
		return
	}

	//シグナル受け取り可にしてチャネル受け取りを待つ（受け取ったら終了）
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func (conf tomlconf) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	cmdtext := strings.Fields(m.Content)

	if includeChannel(conf.Discord.Channel, m.ChannelID) && cmdtext[0] == "!mc" {
		if len(cmdtext) > 2 {
			switch cmdtext[1] {
			case "add":
				cmd := exec.Command("screen", "-S", "minecraft", "-X", "stuff", "whitelist add "+cmdtext[2]+" \r")
				err := cmd.Run()

				if err != nil {
					log.Fatal(err)
				} else {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ADD %s TO WHITELIST", cmdtext[2]))
				}

			case "remove":
				cmd := exec.Command("screen", "-S", "minecraft", "-X", "stuff", "whitelist remove "+cmdtext[2]+" \r")
				err := cmd.Run()

				if err != nil {
					log.Fatal(err)
				} else {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("REMOVE %s FROM WHITELIST", cmdtext[2]))
				}

			case "reload":
				cmd := exec.Command("screen", "-S", "minecraft", "-X", "stuff", "whitelist reload \r")
				err := cmd.Run()

				if err != nil {
					log.Fatal(err)
				} else {
					s.ChannelMessageSend(m.ChannelID, "RELOAD WHITELIST")
				}

			default:
				s.ChannelMessageSend(m.ChannelID, "HELP\n\nユーザ追加 !mc add <userID>\nユーザ削除 !mc remove <userID>\nホワイトリスト再読み込み !mc reload")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "HELP\n\nユーザ追加 !mc add <userID>\nユーザ削除 !mc remove <userID>\nホワイトリスト再読み込み !mc reload")
		}
	}
}

func includeChannel(slice []string, target string) bool {
	for _, num := range slice {
		if num == target {
			return true
		}
	}
	return false
}
