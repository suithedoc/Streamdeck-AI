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
	sh.buttonItToText[TraverseButtonId(buttonId, sh.GetDevice())] = text
	return nil
}

func (sh *StreamdeckHandler) SwitchPage(page int) {
	sh.device.Clear()
	sh.page = page
	err := SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns-1, ">")
	if err != nil {
		log.Fatal(err)
	}
	err = SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns*2-1, strconv.Itoa(sh.page))
	if err != nil {
		log.Fatal(err)
	}
	err = SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns*3-1, "<")
	if err != nil {
		log.Fatal(err)
	}
	for buttonId, text := range sh.buttonItToText {
		reverseButtonId := ReverseTraverseButtonId(buttonId, sh.GetDevice())
		buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
		if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {
			//newId := sh.TraverseButtonId(buttonIdPerPage)
			err := SetStreamdeckButtonText(sh.device, uint8(TraverseButtonId(buttonIdPerPage, sh.GetDevice())), text)
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
				fmt.Printf("Key pressed index reverseTransverse %v, is pressed %v\n", ReverseTraverseButtonId(int(key.Index), sh.GetDevice()), key.Pressed)
				if key.Pressed {
					if handler, ok := sh.GetOnPressHandler(ReverseTraverseButtonId(int(key.Index), sh.GetDevice())); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					if handler, ok := sh.GetOnReleaseHandler(ReverseTraverseButtonId(int(key.Index), sh.GetDevice())); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}()
	return nil
}
