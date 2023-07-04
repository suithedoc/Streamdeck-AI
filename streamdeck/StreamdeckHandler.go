package streamdeck

import (
	"fmt"
	"log"
	"strconv"
)

type StreamdeckHandler struct {
	IStreamdeckHandler
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
	device                               *StreamdeckDeviceWrapper
	page                                 int
	buttonItToText                       map[int]string
	numOfButtons                         int
}

func NewStreamdeckHandler() (IStreamdeckHandler, error) {
	device, err := InitStreamdeckDevice()
	if err != nil {
		return nil, err
	}
	return &StreamdeckHandler{
		streamDeckButtonIdToOnPressHandler:   make(map[int]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[int]func() error),
		device:                               device,
		page:                                 1,
		buttonItToText:                       make(map[int]string),
		numOfButtons:                         int(device.device.Rows * (device.device.Columns - 1)),
	}, nil
}

func (sh *StreamdeckHandler) AddButtonText(buttonId int, text string) error {
	sh.buttonItToText[sh.TraverseButtonId(buttonId)] = text
	return nil
}

func (sh *StreamdeckHandler) SwitchPage(page int) {
	sh.device.Clear()
	sh.page = page
	err := SetStreamdeckButtonText(sh.device, sh.device.device.Columns-1, ">")
	if err != nil {
		log.Fatal(err)
	}
	err = SetStreamdeckButtonText(sh.device, sh.device.device.Columns*2-1, strconv.Itoa(sh.page))
	if err != nil {
		log.Fatal(err)
	}
	err = SetStreamdeckButtonText(sh.device, sh.device.device.Columns*3-1, "<")
	if err != nil {
		log.Fatal(err)
	}
	for buttonId, text := range sh.buttonItToText {
		reverseButtonId := sh.ReverseTraverseButtonId(buttonId)
		buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
		if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {
			//newId := sh.TraverseButtonId(buttonIdPerPage)
			err := SetStreamdeckButtonText(sh.device, uint8(sh.TraverseButtonId(buttonIdPerPage)), text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (sh *StreamdeckHandler) AddOnPressHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnPressHandler[buttonId] = handler
}

func (sh *StreamdeckHandler) AddOnReleaseHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnReleaseHandler[buttonId] = handler
}

func (sh *StreamdeckHandler) GetOnPressHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnPressHandler[buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetOnReleaseHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetDevice() DeviceWrapper {
	return sh.device
}

func (sh *StreamdeckHandler) StartListenAsync() error {
	sh.SwitchPage(1)
	go func() {
		keys, err := sh.device.device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				if key.Pressed {
					reverseButtonId := sh.ReverseTraverseButtonId(int(key.Index))
					buttonIdPerPage := reverseButtonId - (sh.page-1)*int(sh.device.device.Rows*sh.device.device.Columns)
					if reverseButtonId >= 0 && buttonIdPerPage < int(sh.device.device.Rows*sh.device.device.Columns) {
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
	return nil
}

func (sh *StreamdeckHandler) TraverseButtonId(buttonId int) int {
	rows := int(sh.device.device.Rows)
	cols := int(sh.device.device.Columns)

	// Convert the vertical-first index into a row and column.
	row := buttonId % rows
	col := buttonId / rows

	// Convert the row and column into a horizontal-first index.
	newButtonId := row*cols + col

	return newButtonId
}

func (sh *StreamdeckHandler) ReverseTraverseButtonId(buttonId int) int {
	rows := int(sh.device.device.Rows)
	cols := int(sh.device.device.Columns)

	// Convert the horizontal-first index into a row and column.
	row := buttonId / cols
	col := buttonId % cols

	// Convert the row and column into a vertical-first index.
	newButtonId := col*rows + row

	return newButtonId
}
