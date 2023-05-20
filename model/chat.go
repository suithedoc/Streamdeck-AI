package model

import (
	sd "OpenAITest/streamdeck"
	"fmt"
	"github.com/muesli/streamdeck"
	"github.com/sashabaranov/go-openai"
	"log"
	"strconv"
)

type ChatContent struct {
	SystemMsg       string
	PromptMsg       string
	HistoryMessages []openai.ChatCompletionMessage
}

type StreamdeckHandler struct {
	device                               *streamdeck.Device
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
	page                                 int
	buttonItToText                       map[int]string
	numOfButtons                         int
}

func NewStreamdeckHandler(device *streamdeck.Device) *StreamdeckHandler {
	streamdeckHandler := &StreamdeckHandler{
		device:                               device,
		streamDeckButtonIdToOnPressHandler:   make(map[int]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[int]func() error),
		page:                                 1,
		buttonItToText:                       make(map[int]string),
		numOfButtons:                         int(device.Rows * (device.Columns - 1)),
	}
	streamdeckHandler.AddOnPressHandler(streamdeckHandler.ReverseTraverseButtonId(int(device.Columns-1)), func() error {
		streamdeckHandler.SwitchPage(streamdeckHandler.page + 1)
		return nil
	})
	streamdeckHandler.AddOnPressHandler(streamdeckHandler.ReverseTraverseButtonId(int(device.Columns*3-1)), func() error {
		streamdeckHandler.SwitchPage(streamdeckHandler.page - 1)
		return nil
	})
	return streamdeckHandler
}

func (sh *StreamdeckHandler) AddButtonText(buttonId int, text string) error {
	sh.buttonItToText[sh.TraverseButtonId(buttonId)] = text
	return nil
}

func (sh *StreamdeckHandler) AddOnPressHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnPressHandler[sh.TraverseButtonId(buttonId)] = handler
}

func (sh *StreamdeckHandler) AddOnReleaseHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnReleaseHandler[sh.TraverseButtonId(buttonId)] = handler
}

func (sh *StreamdeckHandler) GetOnPressHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnPressHandler[buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetOnReleaseHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[buttonId]
	return handler, ok
}

//func (sh *StreamdeckHandler) TraverseCoslAndRowByButtonId(buttonId int) (col int, row int) {
//	row = buttonId % int(sh.device.Rows)
//	col = buttonId / int(sh.device.Rows)
//	return col, row
//}

func (sh *StreamdeckHandler) SwitchPage(page int) {
	sh.device.Clear()
	sh.page = page
	err := sd.SetStreamdeckButtonText(sh.device, sh.device.Columns-1, ">")
	if err != nil {
		log.Fatal(err)
	}
	err = sd.SetStreamdeckButtonText(sh.device, sh.device.Columns*2-1, strconv.Itoa(sh.page))
	if err != nil {
		log.Fatal(err)
	}
	err = sd.SetStreamdeckButtonText(sh.device, sh.device.Columns*3-1, "<")
	if err != nil {
		log.Fatal(err)
	}
	for buttonId, text := range sh.buttonItToText {
		reverseButtonId := sh.ReverseTraverseButtonId(buttonId)
		buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
		if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {
			//newId := sh.TraverseButtonId(buttonIdPerPage)
			err := sd.SetStreamdeckButtonText(sh.device, uint8(sh.TraverseButtonId(buttonIdPerPage)), text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (sh *StreamdeckHandler) StartAsync() {
	sh.SwitchPage(1)
	//for buttonId, text := range sh.buttonItToText {
	//	err := sd.SetStreamdeckButtonText(sh.device, uint8(buttonId), text)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	go func() {
		keys, err := sh.device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				if key.Pressed {
					reverseButtonId := sh.ReverseTraverseButtonId(int(key.Index))
					buttonIdPerPage := reverseButtonId - (sh.page-1)*int(sh.device.Rows*sh.device.Columns)
					if reverseButtonId >= 0 && buttonIdPerPage < int(sh.device.Rows*sh.device.Columns) {
						if handler, ok := sh.GetOnPressHandler(int(key.Index)); ok {
							err := handler()
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				} else {
					reverseButtonId := sh.ReverseTraverseButtonId(int(key.Index))
					buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
					if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {
						if handler, ok := sh.GetOnReleaseHandler(int(key.Index)); ok {
							err := handler()
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}
			}
		}
	}()
}

func (sh *StreamdeckHandler) TraverseButtonId(buttonId int) int {
	rows := int(sh.device.Rows)
	cols := int(sh.device.Columns)

	// Convert the vertical-first index into a row and column.
	row := buttonId % rows
	col := buttonId / rows

	// Convert the row and column into a horizontal-first index.
	newButtonId := row*cols + col

	return newButtonId
}

func (sh *StreamdeckHandler) ReverseTraverseButtonId(buttonId int) int {
	rows := int(sh.device.Rows)
	cols := int(sh.device.Columns)

	// Convert the horizontal-first index into a row and column.
	row := buttonId / cols
	col := buttonId % cols

	// Convert the row and column into a vertical-first index.
	newButtonId := col*rows + row

	return newButtonId
}
