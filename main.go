package main

import (
	"OpenAITest/bots"
	"OpenAITest/model"
	sd "OpenAITest/streamdeck"
	"OpenAITest/utils"
	"OpenAITest/wakeword"
	"bufio"
	"bytes"
	"fmt"
	"git.tcp.direct/kayos/sendkeys"
	markdown "github.com/MichaelMure/go-term-markdown"
	fcolor "github.com/fatih/color"
	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	_ "github.com/mattn/go-sqlite3"
	"github.com/micmonay/keybd_event"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tinyzimmer/go-gst/gst"
	"go.mau.fi/whatsmeow/types/events"
	"golang.design/x/clipboard"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	//style  lipgloss.Style
	color  *fcolor.Color
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
	streamdeckHandler  sd.IStreamdeckHandler
	kb                 keybd_event.KeyBonding
)

func init() {
	echoWithColorQueue = make(chan EchoWithColor, 100)
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

//
//func StartListenStreamDeckAsync(device *streamdeck.Device) error {
//
//	go func() {
//		keys, err := device.ReadKeys()
//		if err != nil {
//			log.Fatal(err)
//		}
//		for {
//			select {
//			case key := <-keys:
//				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
//				if key.Pressed {
//					if handler, ok := streamdeckHandler.GetOnPressHandler(int(key.Index)); ok {
//						err := handler()
//						if err != nil {
//							log.Fatal(err)
//						}
//					}
//				} else {
//					if handler, ok := streamdeckHandler.GetOnReleaseHandler(int(key.Index)); ok {
//						err := handler()
//						if err != nil {
//							log.Fatal(err)
//						}
//					}
//				}
//			}
//		}
//	}()
//
//	/*
//		ver, err := d.FirmwareVersion()
//		if err != nil {
//			return fmt.Errorf("can't retrieve device info: %s", err)
//		}
//		fmt.Printf("Found device with serial %s (firmware %s)\n",
//			d.Serial, ver)Hello, how are zou_ Hallo, wie geht es dir_ mir geht das eigentlich gany gut
//	*/
//
//	return nil
//}

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

func downloadUrlContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func main() {
	var err error

	//device, err := sd.InitStreamdeckDevice()
	//if err != nil {
	//	log.Fatal(err)
	//}
	streamdeckHandler, err = sd.NewStreamdeckHandler()
	if err != nil || streamdeckHandler == nil {
		streamdeckHandler, err = sd.NewUiStreamdeckHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	device := streamdeckHandler.GetDevice()
	if device == nil {
		log.Fatal("No device found")
	}
	err = device.Clear()
	if err != nil {
		log.Fatal(err)
	}
	defer func(device sd.DeviceWrapper) {
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

	botFactory := bots.NewBotFactory(streamdeckHandler, client, device, &kb)

	assistantButtonConfig := sd.StreamdeckButtonConfig{
		ButtonIndex:               0,
		ButtonIndexHistory:        1,
		ButtonIndexHistoryAndCopy: 2,
		Page:                      0,
	}
	assistantBot := botFactory.CreateBot("Assistant", properties["assistantSystemMsg"], properties["assistantPromptMsg"], assistantButtonConfig, voices.German)
	assistantBot.AddResponseListener(func(bot *bots.AiBot, s string) error {
		fmt.Printf("Assistant: %s\n", s)
		return nil
	})
	assistantBot.AddResponseListener(bots.SpeakResultFunc)

	commanderButtonConfig := sd.StreamdeckButtonConfig{
		ButtonIndex:               3,
		ButtonIndexHistory:        4,
		ButtonIndexHistoryAndCopy: 5,
		Page:                      0,
	}
	commanderBot := botFactory.CreateBot("Commander", properties["commanderSystemMsg"], properties["commanderPromptMsg"], commanderButtonConfig, voices.German)
	commanderBot.AddResponseListener(bots.ExecuteCommandResultFunc)

	//labelBot := botFactory.CreateBotWithHistory("Label",
	//	"Analyse the Input that I give you and only respond with one word. A label that best describes the input."+
	//		"If it is necessary to search the internet for resolving the request, respond with 'search'."+
	//		"If the request contains a development question, respond with 'code'."+
	//		"If the request could be resolved by a linux command, respond with 'linux'.",
	//	"",
	//	15, 16, voices.German)
	//labelBot.AddResponseListener(func(bot *bots.AiBot, s string) error {
	//
	//	fmt.Printf("Label: %s\n", s)
	//	if strings.Contains(s, "search") {
	//		lastMessage := bot.CompletionHistory[0]
	//		fmt.Println("Last message: ", lastMessage.Content)
	//		fmt.Println("Last message: ", lastMessage.Content)
	//		err := commanderBot.EvaluateGptResponseStrings([]string{"To solve the following request, search with googler --noprompt and use" +
	//			"wget to download the content of the first url found with googler.", lastMessage.Content})
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//	}
	//	return nil
	//})

	codeButtonConfig := sd.StreamdeckButtonConfig{
		ButtonIndex:               8,
		ButtonIndexHistory:        11,
		ButtonIndexHistoryAndCopy: 10,
		Page:                      0,
	}
	codeBot := botFactory.CreateBot("Cpp", properties["cppSystemMsg"], properties["cppPromptMsg"], codeButtonConfig, voices.German)
	codeBot.AddResponseListener(func(bot *bots.AiBot, answer string) error {
		codeBlocks := utils.ExtractCodeBlockFromMarkdown(answer)
		if len(codeBlocks) == 0 {
			codeBlocks = utils.ExtractCodeBlockFromMarkdownWithOneBacktick(answer)
		}
		if len(codeBlocks) == 0 {
			log.Println("No code block found")
			fmt.Println("No Code blocks found in this: ", answer)
		}
		fmt.Println("Code blocks found: ", len(codeBlocks))
		allCodeBlocks := strings.Join(codeBlocks, "\n")
		err = utils.TypeCodeCommand(allCodeBlocks, &kb)
		if err != nil {
			fmt.Printf("Error typing command: %v\n", err)
		}
		return nil
	})

	whisperButtonConfig := sd.StreamdeckButtonConfig{
		ButtonIndex:               6,
		ButtonIndexHistory:        -1,
		ButtonIndexHistoryAndCopy: -1,
		Page:                      0,
	}
	whisperBot := botFactory.CreateBot("STT", "", "", whisperButtonConfig, voices.German)
	whisperBot.DisableAi = true
	whisperBot.AddResponseListener(func(bot *bots.AiBot, s string) error {
		fmt.Printf("Assistant: %s\n", s)
		return nil
	})
	whisperBot.AddTranscriptionListener(func(bot *bots.AiBot, s string) error {
		fmt.Printf("Assistant: %s\n", s)
		err = bots.TypeWhisperSTT(s, &kb)
		if err != nil {
			fmt.Printf("Error typing command: %v\n", err)
			return err
		}
		return nil
	})

	//err = TypeWhisperSTT(transcription, kb)
	//if err != nil {
	//	return err
	//}
	//whisperBot.AddResponseListener(bots.SpeakResultFunc)

	//englishTeacherBot := botFactory.CreateBotWithHistory("English",
	//	properties["englishTeacherSystemMsg"],
	//	properties["englishTeacherPromptMsg"],
	//	1,
	//	2, voices.English)
	//englishTeacherBot.AddResponseListener(bots.SpeakResultFunc)
	//bots.InitWhisperBot(streamdeckHandler, device, &kb, client, 2)
	//bots.InitMinecraftGPTBot(client, device, properties, streamdeckHandler, speech, &kb, 9)
	//assistantBot.AddResponseListener(bots.SpeakResultFunc)
	//bots.InitCodeGPTBot(client, device, properties, streamdeckHandler, speech, &kb, 10)

	//err = InitWakeWordCommander(properties, commanderChatContent, &kb, client)
	//if err != nil {
	//	log.Printf("initializing wakeword: %v", err)
	//	os.Exit(1)
	//}

	//err = streamdeckHandler.StartListenAsync()
	//if err != nil {
	//	log.Fatal(err)
	//}

	gst.Init(nil)
	AsyncConsoleOutput()

	//style = lipgloss.NewStyle().
	//	Bold(true).
	//	Background(lipgloss.Color("#7D56F4")).
	//	PaddingTop(2).
	//	PaddingLeft(4).
	//	Width(22)

	color = fcolor.New(fcolor.FgCyan).Add(fcolor.BgWhite)
	// Scanner

	err = InitWakeWord(properties, &kb)
	if err != nil {
		log.Fatal(fmt.Errorf("initializing wakeword: %v", err))
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
	err = streamdeckHandler.StartListenAsync()
	if err != nil {
		log.Fatal(err)
	}

	startCommandlineInput(assistantBot, commanderBot)
}

func startCommandlineInput(assistantBot *bots.AiBot, commanderBot *bots.AiBot) {
	isRecording := false
	quitChannel := make(chan bool)
	finished := make(chan bool)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !isRecording {
			fmt.Print("Enter Multiline Input: \n")
		} else {
			fmt.Print("Stop recording with 'r' \n")
		}
		color.Print(">")
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
			assistantBot.EvaluateGptResponseStrings(input)
			continue
		} else if strings.HasPrefix(input[0], "ah:") {
			assistantBot.EvaluateGptResponseStringsWithHistory(input)
		} else if strings.HasPrefix(input[0], "h:") {
			commanderBot.EvaluateGptResponseStringsWithHistory(input)
		} else {
			commanderBot.EvaluateGptResponseStrings(input)
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
		color.Print(">")
	}

	return input
}
