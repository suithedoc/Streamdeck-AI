package streamdeck

import (
	"fmt"
	"log"
	"strconv"
)

type UiStreamdeckHandler struct {
	IStreamdeckHandler
	device                               *UiDeviceWrapper
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
	page                                 int
	buttonIdToText                       map[int]string
	numOfButtons                         int
}

func NewUiStreamdeckHandler() (IStreamdeckHandler, error) {
	device, err := InitUiStreamdeck()
	if err != nil {
		return nil, err
	}
	return &UiStreamdeckHandler{
		device:                               device,
		streamDeckButtonIdToOnPressHandler:   make(map[int]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[int]func() error),
		page:                                 1,
		buttonIdToText:                       make(map[int]string),
		numOfButtons:                         int(device.GetRows() * (device.GetColumns() - 1)),
	}, nil
}

func (sh *UiStreamdeckHandler) AddButtonText(buttonId int, text string) error {
	//sh.device.buttons[TraverseButtonId(buttonId, sh.GetDevice())].SetText(text)
	sh.buttonIdToText[TraverseButtonId(buttonId, sh.GetDevice())] = text
	return nil
}

func (sh *UiStreamdeckHandler) GetDevice() DeviceWrapper {
	return sh.device
}

func (sh *UiStreamdeckHandler) StartListenAsync() error {
	sh.SwitchPage(1)
	go func() {
		eventChan, err := sh.device.ReadEvents()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			event := <-eventChan
			if event.IsPressed {
				handler, ok := sh.GetOnPressHandler(ReverseTraverseButtonId(event.ButtonId, sh.GetDevice()))
				if ok {
					err := handler()
					if err != nil {
						fmt.Println(err)
					}
				}
			} else if event.IsReleased {
				handler, ok := sh.GetOnReleaseHandler(ReverseTraverseButtonId(event.ButtonId, sh.GetDevice()))
				if ok {
					err := handler()
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()
	sh.device.window.ShowAndRun()
	return nil
}

func (sh *UiStreamdeckHandler) AddOnPressHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnPressHandler[buttonId] = handler
}

func (sh *UiStreamdeckHandler) AddOnReleaseHandler(buttonId int, handler func() error) {
	sh.streamDeckButtonIdToOnReleaseHandler[buttonId] = handler
}

func (sh *UiStreamdeckHandler) GetOnPressHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnPressHandler[buttonId]
	return handler, ok
}

func (sh *UiStreamdeckHandler) GetOnReleaseHandler(buttonId int) (func() error, bool) {
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[buttonId]
	return handler, ok
}

func (sh *UiStreamdeckHandler) SwitchPage(page int) {
	err := sh.device.Clear()
	if err != nil {
		log.Fatal(err)
		return //log.Fatal(err)
	}
	sh.page = page
	//err := SetStreamdeckButtonText(sh.GetDevice(), uint8(sh.device.GetColumns()-1), ">")
	cols := sh.device.GetColumns()
	err = sh.GetDevice().SetText(ReverseTraverseButtonId(cols-1, sh.GetDevice()), ">")
	if err != nil {
		log.Fatal(err)
	}
	err = sh.GetDevice().SetText(ReverseTraverseButtonId(cols*2-1, sh.GetDevice()), strconv.Itoa(sh.page))
	if err != nil {
		log.Fatal(err)
	}
	err = sh.GetDevice().SetText(ReverseTraverseButtonId(cols*3-1, sh.GetDevice()), "<")
	if err != nil {
		log.Fatal(err)
	}
	for buttonId, text := range sh.buttonIdToText {
		reverseButtonId := ReverseTraverseButtonId(buttonId, sh.GetDevice())
		buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
		if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {

			err = sh.GetDevice().SetText(TraverseButtonId(buttonIdPerPage, sh.GetDevice()), text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (sh *UiStreamdeckHandler) TraverseButtonId(buttonId int) int {
	return buttonId
}

func (sh *UiStreamdeckHandler) ReverseTraverseButtonId(buttonId int) int {
	return buttonId
}
