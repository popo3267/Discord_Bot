package main

import (
	controller "Discord_Bot/Controller"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token string
)

func main() {
	enverr := godotenv.Load() // 載入 .env 檔案
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}

	//建立新的Token連線
	dis, err := discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		log.Fatal(err)
	}
	// 設定 Bot Intents
	dis.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	//註冊功能處理器
	dis.AddHandler(controller.HandleSaveImage) //儲存圖片
	//dis.AddHandler(handleSlashCommand)
	dis.AddHandler(controller.HandleCallImage) //呼叫圖片
	//dis.AddHandler(handleDeleteImage) //刪除圖片
	dis.AddHandler(controller.HandleListImage) //列出圖片
	//Bot建立Websocket連線
	err = dis.Open()
	if err != nil {
		fmt.Println("Error opening connection", err)
		return
	}

	// 確保 State.User 已初始化
	if dis.State.User == nil {
		log.Fatalf("Bot User 尚未初始化，請確認是否已成功連接至 Discord。")
	}
	// 註冊 Slash Command
	// cmd := &discordgo.ApplicationCommand{
	// 	Name:        "weather",
	// 	Description: "查詢天氣",
	// }
	// _, err = dis.ApplicationCommandCreate(dis.State.User.ID, "", cmd)
	// if err != nil {
	// 	log.Fatalf("無法註冊slash指令:%v", err)
	// }

	//建立sc通道，使得程式可接收SIGINT以及SIGTERM訊號
	fmt.Println("Bot is now running. Press CTRL+C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	//關閉連線
	dis.Close()

}

// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) { //s為Discord連線會話 m為訊息內容or訊息作者
// 	if m.Author.ID == s.State.User.ID { //忽略機器人本身發送的訊息
// 		return
// 	}
// 	if m.Content == "ping" {
// 		s.ChannelMessageSend(m.ChannelID, "Pong!")
// 	}
// 	if m.Content == "pong" {
// 		s.ChannelMessageSend(m.ChannelID, "Ping!")
// 	}

// }
