package controller

import (
	model "Discord_Bot/Models"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleSaveImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 確保只回應非機器人訊息
	if m.Author.Bot {
		return
	}

	//檢測!save指令
	if len(m.Content) == 2 && m.Content[:2] == "!s" {
		s.ChannelMessageSend(m.ChannelID, "請提供Keyword")
		return
	}

	if len(m.Content) > 2 && m.Content[:2] == "!s" {
		keyword := m.Content[3:]
		if keyword == "" {
			s.ChannelMessageSend(m.ChannelID, "請提供Keyword")
			return
		}
		if len(m.Attachments) == 0 {
			s.ChannelMessageSend(m.ChannelID, "請附上儲存的圖片")
			return
		}
		attachments := m.Attachments[0]
		imageURL := attachments.URL
		user := m.Author.Username
		CreateAt := getTaipeiTime()
		s.ChannelMessage(m.ChannelID, imageURL)

		//先查詢資料庫裡有無關鍵字
		filter := bson.M{"keyword": keyword}
		ImageCount, err := ImageCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			log.Fatalf("Save Collection Count Error:%v", err)
		}
		//不管有無資料皆新增至DB
		imageInfo := model.ImageManage{
			ImageID:    ImageCount + 1,
			Keyword:    keyword,
			ImageURL:   imageURL,
			UserName:   user,
			CreateTime: CreateAt,
		}
		//儲存至MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, insertErr := ImageCollection.InsertOne(ctx, imageInfo)
		if insertErr != nil {
			s.ChannelMessageSend(m.ChannelID, "圖片儲存失敗")
			return
		}
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
		if err != nil {
			log.Printf("添加Emoji失敗:%v", err)
		}
		return
	}
}

func HandleCallImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 確保只回應非機器人訊息
	if m.Author.Bot {
		return
	}

	//檢測!keyword指令
	if len(m.Content) >= 1 && m.Content[:1] == "!" {
		keyword := m.Content[1:]
		//先查詢資料庫裡有無關鍵字
		filter := bson.M{"keyword": keyword}
		var result model.ImageManage
		ImageConut, err := ImageCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			log.Printf("Find Collection Count Error:%v", err)
			return
		}
		//無資料時的處理
		if ImageConut == 0 {
			return
		} else if ImageConut == 1 { //單一資料時的處理
			err := ImageCollection.FindOne(context.Background(), filter).Decode(&result)
			if err != nil {
				log.Printf("Error:%v", err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, result.ImageURL)
			return
		} else if ImageConut > 1 {
			//多資料時的處理
			filter1 := bson.M{
				"keyword": keyword,
				"imageid": int64(rand.Intn(int(ImageConut)) + 1),
			}
			err := ImageCollection.FindOne(context.Background(), filter1).Decode(&result)
			if err != nil {
				log.Printf("Error:%v", err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, result.ImageURL)
			return
		}
	}
}

func HandleDeleteImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 確保只回應非機器人訊息
	if m.Author.Bot {
		return
	}

	//檢測!delete指令
	if len(m.Content) > 2 && m.Content[:2] == "!d" {
		keyword := m.Content[3:]
		//先查詢資料庫裡有無關鍵字
		filter := bson.M{"keyword": keyword}
		var result model.ImageManage
		err := ImageCollection.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			//無資料時的處理
			return
		}
		//有資料時的處理
		//至MongoDB刪除資料
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, deleteErr := ImageCollection.DeleteOne(ctx, filter)
		if deleteErr != nil {
			s.ChannelMessageSend(m.ChannelID, "")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "")
		return
	}
}

func HandleListImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 確保只回應非機器人訊息
	if m.Author.Bot {
		return
	}

	//檢測!keyword指令
	if len(m.Content) >= 5 && m.Content[:5] == "!list" {
		keyword := m.Content[6:]
		//先查詢資料庫裡有無關鍵字
		filter := bson.M{"keyword": keyword}
		//設定搜尋時間
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ListImage, err := ImageCollection.Find(ctx, filter)
		if err != nil {
			log.Printf("List Collection Error:%v", err)
		}
		defer ListImage.Close(ctx)
		var results []model.ImageManage
		if err := ListImage.All(ctx, &results); err != nil {
			log.Printf("List解析失敗:%v", err)
		}
		var ListResult string

		for _, result := range results {
			// 避免 Discord Markdown 語法解析
			ListResult += fmt.Sprintf("編號:%v 上傳者:%v: %v 上傳時間:%v\n", result.ImageID, result.UserName, result.ImageURL, result.CreateTime)
		}
		s.ChannelMessageSend(m.ChannelID, ListResult)
		return
	}
}

func getTaipeiTime() time.Time {
	taipeiOffset := 8 * 60 * 60 // UTC+8 的秒數
	return time.Now().UTC().Add(time.Second * time.Duration(taipeiOffset))
}
