package bots

import (
	"OpenAITest/model"
	"OpenAITest/utils"
	"context"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	userChatHistory = make(map[string][]openai.ChatCompletionMessage)
}

var userChatHistory map[string][]openai.ChatCompletionMessage

func StartWhatsappBot(chatContent model.ChatContent, userNumber string, userName string, client *openai.Client) error {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	whatsappClient := whatsmeow.NewClient(deviceStore, nil)
	whatsappClient.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("Received message from: " + v.Info.PushName)

			if userNumber == "" && userName == "" {
				fmt.Println("No user number or name specified, returning")
				return
			}
			if !strings.Contains(v.Info.Sender.User, userNumber) || userNumber == "" {
				if v.Info.PushName != userName && userName != "" {
					return
				}
			}
			time.Sleep(5 * time.Second)

			newChatContent := chatContent
			if _, ok := userChatHistory[v.Info.Sender.User]; ok {
				for _, msg := range userChatHistory[v.Info.Sender.User] {
					newChatContent.HistoryMessages = append(newChatContent.HistoryMessages, msg)
				}
			}
			answer, err := utils.SendChatRequest(v.Message.GetConversation(), newChatContent, client)
			if err != nil {
				log.Fatal(err)
			}

			_, err = whatsappClient.SendMessage(context.Background(), v.Info.Sender.ToNonAD(), &waProto.Message{
				Conversation: proto.String(answer),
			})
			if err != nil {
				fmt.Println("Error seding message: " + err.Error())
			}
			fmt.Printf("Received a message:%v, received from:%v\n", v.Message.GetConversation(), v.Info.Sender)
			//check if user exists in map
			if _, ok := userChatHistory[v.Info.Sender.User]; !ok {
				userChatHistory[v.Info.Sender.User] = make([]openai.ChatCompletionMessage, 0)
			}
			userChatHistory[v.Info.Sender.User] = append(userChatHistory[v.Info.Sender.User], openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: v.Message.GetConversation(),
			})
			userChatHistory[v.Info.Sender.User] = append(userChatHistory[v.Info.Sender.User], openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: answer,
			})
		}
	})

	if whatsappClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := whatsappClient.GetQRChannel(context.Background())
		err = whatsappClient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				//e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = whatsappClient.Connect()
		if err != nil {
			panic(err)
		}
	}
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsappClient.Disconnect()
	return nil
}

func InitWhatsappBot(client *openai.Client, properties map[string]string) *model.ChatContent {
	whatsappSystemMsg := properties["whatsappSystemMsg"]
	whatsappPromptMsg := properties["whatsappPromptMsg"]
	whatsappChatContent := model.ChatContent{
		SystemMsg:       whatsappSystemMsg,
		PromptMsg:       whatsappPromptMsg,
		HistoryMessages: []openai.ChatCompletionMessage{},
	}

	go func() {
		err := StartWhatsappBot(whatsappChatContent, properties["whatsappNumber"], properties["whatsappName"], client)
		if err != nil {
			fmt.Printf("Stopped WhatsappBot %v", err.Error())
			return
		}
	}()
	return nil
}
