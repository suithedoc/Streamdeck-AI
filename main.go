package main

import (
	"OpenAITest/bots"
	"OpenAITest/model"
	"OpenAITest/utils"
	"OpenAITest/wakeword"
	"bufio"
	"bytes"
	"fmt"
	"git.tcp.direct/kayos/sendkeys"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/charmbracelet/lipgloss"
	fcolor "github.com/fatih/color"
	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/muesli/streamdeck"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tinyzimmer/go-gst/gst"
	"go.mau.fi/whatsmeow/types/events"
	"golang.design/x/clipboard"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	sd "OpenAITest/streamdeck"
	_ "github.com/mattn/go-sqlite3"
	"github.com/micmonay/keybd_event"
)

var (
	style  lipgloss.Style
	c      *fcolor.Color
	client *openai.Client
	speech *htgotts.Speech
	kbw    *sendkeys.KBWrap
)

type EchoWithColor struct {
	c   *fcolor.Color
	msg string
}

var (
	echoWithColorQueue chan EchoWithColor
	streamdeckHandler  *model.StreamdeckHandler
	kb                 keybd_event.KeyBonding
)

func init() {
	echoWithColorQueue = make(chan EchoWithColor, 100)
	streamdeckHandler = model.NewStreamdeckHandler()
	utils.OpenaiModel = "gpt3.5"
}

func AddConsoleOutputWithColorToQueue(c *fcolor.Color, msg string) {
	echoWithColorQueue <- EchoWithColor{c: c, msg: msg}
}

func AsyncConsoleOutput() {
	go func() {
		for {
			select {
			//case output := <-printerQueue:
			//	fmt.Print(output)
			case echoWithColor := <-echoWithColorQueue:
				if echoWithColor.c == nil {
					fmt.Print(echoWithColor.msg)
					continue
				} else {
					echoWithColor.c.Print(echoWithColor.msg)
				}
			}
		}
	}()
}

func StartListenStreamDeckAsync(device *streamdeck.Device) error {

	go func() {
		keys, err := device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				if key.Pressed {
					if handler, ok := streamdeckHandler.GetOnPressHandler(int(key.Index)); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					if handler, ok := streamdeckHandler.GetOnReleaseHandler(int(key.Index)); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}()

	/*
		ver, err := d.FirmwareVersion()
		if err != nil {
			return fmt.Errorf("can't retrieve device info: %s", err)
		}
		fmt.Printf("Found device with serial %s (firmware %s)\n",
			d.Serial, ver)Hello, how are zou_ Hallo, wie geht es dir_ mir geht das eigentlich gany gut
	*/

	return nil
}

func fetch(text string) (io.Reader, error) {
	data := []byte(text)

	chunkSize := len(data)
	if len(data) > 32 {
		chunkSize = 32
	}

	urls := make([]string, 0)
	for prev, i := 0, 0; i < len(data); i++ {
		if i%chunkSize == 0 && i != 0 {
			chunk := string(data[prev:i])
			url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=%d&client=tw-ob&q=%s&tl=%s", chunkSize, url.QueryEscape(chunk), speech.Language)
			urls = append(urls, url)
			prev = i
		} else if i == len(data)-1 {
			chunk := string(data[prev:])
			url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=%d&client=tw-ob&q=%s&tl=%s", chunkSize, url.QueryEscape(chunk), speech.Language)
			urls = append(urls, url)
			prev = i
		}
	}

	buf := new(bytes.Buffer)
	for _, url := range urls {
		r, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
	}
	return buf, nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		//message, err := client.SendMessage(context.Background(), toJid, &waProto.Message{
		//	Conversation: proto.String("Hello, World!"),
		//})
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func InitWakeWordCommander(properties map[string]string, commanderChatContent *model.ChatContent, kb *keybd_event.KeyBonding, client *openai.Client) (err error) {
	porcupineAccessKey := properties["PorcupineAccessKey"]
	wakeWordChannel := make(chan bool)
	go func() {
		err = wakeword.StartListeningToWakeword("jarvis", 0.3, porcupineAccessKey, wakeWordChannel)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		finishedChannel := make(chan bool, 2)
		for {
			select {
			case <-wakeWordChannel:
				fmt.Println("Wake word detected")
				err := speech.Speak("Yes master?")
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(time.Millisecond * 200)
				utils.RecordAndSaveAudioAsWav("test.wav", time.Millisecond*800, finishedChannel)
				//case <-finishedChannel:
				fmt.Println("Finished recording")
				transcription, err := utils.ParseMp3ToText("test.wav", client)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
				}
				bots.EvaluateCommanderGptResponseStrings([]string{transcription}, false, nil, *commanderChatContent, client)
				//err = bots.TypeWhisperSTT(transcription, kb)
				//if err != nil {
				//	fmt.Printf("Error typing whisper: %s\n", err)
				//}
				finishedChannel = make(chan bool, 2)
			}
		}
	}()
	return nil
}

func InitWakeWord(properties map[string]string, kb *keybd_event.KeyBonding) (err error) {
	porcupineAccessKey := properties["PorcupineAccessKey"]
	wakeWordChannel := make(chan bool)
	go func() {
		err = wakeword.StartListeningToWakeword("computer", 0.3, porcupineAccessKey, wakeWordChannel)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		finishedChannel := make(chan bool, 2)
		for {
			select {
			case <-wakeWordChannel:
				fmt.Println("Wake word detected")
				err := speech.Speak("Yes master?")
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(time.Millisecond * 200)
				utils.RecordAndSaveAudioAsWav("test.wav", time.Millisecond*800, finishedChannel)
				//case <-finishedChannel:
				fmt.Println("Finished recording")
				transcription, err := utils.ParseMp3ToText("test.wav", client)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text 2: %s\n", err)
				}
				err = bots.TypeWhisperSTT(transcription, kb)
				if err != nil {
					fmt.Printf("Error typing whisper: %s\n", err)
				}
				finishedChannel = make(chan bool, 2)
			}
		}
	}()
	return nil
}

func main() {
	var err error

	scanner := bufio.NewScanner(os.Stdin)
	quitChannel := make(chan bool)
	finished := make(chan bool)
	isRecording := false
	device, err := sd.InitStreamdeckDevice()
	if err != nil {
		log.Fatal(err)
	}
	err = device.Clear()
	if err != nil {
		log.Fatal(err)
	}
	defer func(device *streamdeck.Device) {
		err := device.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(device)
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	kbw, err = sendkeys.NewKBWrapWithOptions(sendkeys.Noisy)
	if err != nil {
		println(err.Error())
		return
	}

	//properties, err :=LoadPropertiesFromJsonFile("config.json")
	properties, err := LoadPropertiesFromIniFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	if aimodel, ok := properties["model"]; ok {
		fmt.Printf("Using model %v\n\n", aimodel)
		utils.OpenaiModel = aimodel
	}

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	//speech = &htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	speech = &htgotts.Speech{Folder: "audio", Language: voices.German, Handler: &handlers.Native{}}
	client = openai.NewClient(properties["apiKey"])

	if properties["whatsapp"] == "enabled" {
		bots.InitWhatsappBot(client, properties)
	}

	if properties["discord"] == "enabled" {
		bots.InitDiscordBot(client, properties)
	}

	commanderChatContent := bots.InitCommanderGPTBot(client, device, properties, streamdeckHandler, scanner, 0, 5)
	assistantChatContent := bots.InitAssistantGPTBot(client, device, properties, streamdeckHandler, speech, &kb, 1, 6, 11)
	bots.InitWhisperBot(streamdeckHandler, device, &kb, client, 2)
	bots.InitMinecraftGPTBot(client, device, properties, streamdeckHandler, speech, &kb, 3)
	bots.InitCodeGPTBot(client, device, properties, streamdeckHandler, speech, &kb, 4)

	//err = InitWakeWordCommander(properties, commanderChatContent, &kb, client)
	//if err != nil {
	//	log.Printf("initializing wakeword: %v", err)
	//	os.Exit(1)
	//}

	//evaluators.InitGoogooGPTBot(client, device, properties, streamdeckHandler, scanner, 10, 11)
	err = StartListenStreamDeckAsync(device)
	if err != nil {
		println(err.Error())
	}

	gst.Init(nil)
	AsyncConsoleOutput()

	style = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)

	c = fcolor.New(fcolor.FgCyan).Add(fcolor.BgWhite)
	// Scanner

	err = InitWakeWord(properties, &kb)
	if err != nil {
		log.Printf("initializing wakeword: %v", err)

		log.Fatal(err)
	}

	//porcupineAccessKey := properties["PorcupineAccessKey"]
	//wakeWordChannel := make(chan bool)
	//go func() {
	//	err = wakeword.StartListeningToWakeword("computer", 0.3, porcupineAccessKey, wakeWordChannel)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//go func() {
	//	for {
	//		select {
	//		case <-wakeWordChannel:
	//			err := speech.Speak("Yes master?")
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//		}
	//	}
	//}()

	for {
		if !isRecording {
			fmt.Print("Enter Multiline Input: \n")
		} else {
			fmt.Print("Stop recording with 'r' \n")
		}
		c.Print(">")
		input := RunAndWaitForInputMessage(scanner)
		if len(input) == 1 && input[0] == "r" {
			if isRecording {
				log.Println("Stopping recording")
				quitChannel <- true
				<-finished
				isRecording = false
				mp3Text, err := utils.ParseMp3ToText("audio.wav", client)
				if err != nil {
					fmt.Printf("Error parsing mp3 to text3: %s\n", err)
				}
				fmt.Println(mp3Text)
				//clear input
				input = []string{mp3Text}
			} else {
				go func() {
					fmt.Println("Recording audio")
					isRecording = true
					utils.RecordAndSaveAudioAsMp3("audio.wav", quitChannel, finished)
				}()
				continue
			}
		}
		if len(input) == 0 {
			continue
		}

		if strings.HasPrefix(input[0], "a:") {
			bots.EvaluateAssistantGptResponseStrings(input, false, *assistantChatContent, client, speech)
			continue
		} else if strings.HasPrefix(input[0], "ah:") {
			bots.EvaluateAssistantGptResponseStrings(input, true, *assistantChatContent, client, speech)
		} else if strings.HasPrefix(input[0], "h:") {
			bots.EvaluateCommanderGptResponseStrings(input, true, scanner, *commanderChatContent, client)
		} else {
			bots.EvaluateCommanderGptResponseStrings(input, false, scanner, *commanderChatContent, client)
		}
	}
}

func typeMultipleMinecraftCommands(commands string) error {
	for _, command := range strings.Split(commands, "\n") {
		err := typeMinecraftCommand(command)
		if err != nil {
			return err
		}
		time.Sleep(20 * time.Millisecond)
	}
	return nil
}

func typeMinecraftCommand(command string) error {
	kb.SetKeys(keybd_event.VK_T)
	err := kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}

	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}
	clipboard.Write(clipboard.FmtText, []byte(command))
	time.Sleep(20 * time.Millisecond)
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	time.Sleep(20 * time.Millisecond)
	kb.HasCTRL(false)
	kb.SetKeys(keybd_event.VK_ENTER)
	err = kb.Launching()
	if err != nil {
		fmt.Printf("Error launching keyboard: %v\n", err)
		return err
	}
	return nil
}

func EvaluateCodeGptResponseStrings(input []string, chatContent model.ChatContent) {
	answer, err := utils.SendChatRequest(strings.Join(input, "\n"), chatContent, client)
	if err != nil {
		log.Fatal(err)
	}
	result := markdown.Render(answer, 80, 6)
	fmt.Println(string(result))
	codeBlocks := utils.ExtractCodeBlockFromMarkdown(answer)
	if len(codeBlocks) == 0 {
		codeBlocks = utils.ExtractCodeBlockFromMarkdownWithOneBacktick(answer)
	}
	if len(codeBlocks) == 0 {
		log.Println("No code block found")
	}
	var commandReturnValue string
	if len(codeBlocks) > 1 {
		err := speech.Speak("More than one command found")
		if err != nil {
			fmt.Printf("Error speaking: %v\n", err)
		}
		fmt.Println("More than one command found")
		for i, codeBlock := range codeBlocks {
			fmt.Printf("%d: %s\n", i, codeBlock)
		}
	} else if len(codeBlocks) == 1 {
		err := typeMultipleMinecraftCommands(codeBlocks[0])
		if err != nil {
			fmt.Printf("Error typing command: %v\n", err)
		}
	}
	commandReturnValueRendered := markdown.Render(commandReturnValue, 80, 6)
	fmt.Println(string(commandReturnValueRendered))
}

func RunAndWaitForInputMessage(scanner *bufio.Scanner) []string {
	input := []string{}

	for {
		// Scans a line from Stdin(Console)
		scanner.Scan()

		// Holds the string that scanned
		text := scanner.Text()
		if text == "r" {
			input = append(input, text)
			return input
		}
		if len(strings.TrimSpace(text)) != 0 {
			input = append(input, text)
		} else {
			break
		}
		c.Print(">")
	}

	return input
}
