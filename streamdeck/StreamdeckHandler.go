package streamdeck

import (
	"fmt"
	"log"
	"math"
)

type StreamdeckHandler struct {
	IStreamdeckHandler
	streamDeckButtonIdToOnPressHandler   map[Page]map[Index]func() error
	streamDeckButtonIdToOnReleaseHandler map[Page]map[Index]func() error
	device                               *StreamdeckDeviceWrapper
	page                                 Page
	buttonIdToText                       map[Page]map[Index]string
	numOfButtons                         int
}

func NewStreamdeckHandler() (IStreamdeckHandler, error) {
	device, err := InitStreamdeckDevice()
	if err != nil {
		return nil, err
	}
	return &StreamdeckHandler{
		streamDeckButtonIdToOnPressHandler:   make(map[Page]map[Index]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[Page]map[Index]func() error),
		device:                               device,
		page:                                 1,
		buttonIdToText:                       make(map[Page]map[Index]string),
		numOfButtons:                         int(device.device.Rows * (device.device.Columns - 1)),
	}, nil
}

func (sh *StreamdeckHandler) GetButtonIndexToText() map[Page]map[Index]string {
	return sh.buttonIdToText
}

func (sh *StreamdeckHandler) GetPage() Page {
	return sh.page
}

func (sh *StreamdeckHandler) SetPage(page Page) {
	sh.page = page
}

func (sh *StreamdeckHandler) AddButtonText(page Page, buttonId Index, text string) error {
	if sh.buttonIdToText[page] == nil {
		sh.buttonIdToText[page] = make(map[Index]string)
	}
	sh.buttonIdToText[page][TraverseButtonId(buttonId, sh.GetDevice())] = text
	return nil
}

//func (sh *StreamdeckHandler) SwitchPage(page Page) {
//	err := sh.device.Clear()
//	if err != nil {
//		log.Fatal(err)
//		return //log.Fatal(err)
//	}
//	sh.page = page
//	//err := SetStreamdeckButtonText(sh.GetDevice(), uint8(sh.device.GetColumns()-1), ">")
//	cols := sh.device.GetColumns()
//	nextPageButtonId := cols - 1
//	pageNumberButtonId := cols*2 - 1
//	prevPageButtonId := cols*3 - 1
//
//	reverseTraverseNextPageButtonId := ToConvinientVerticalId(Index(nextPageButtonId), sh.GetDevice())
//	reverseTraversePageNumberButtonId := ToConvinientVerticalId(Index(pageNumberButtonId), sh.GetDevice())
//	reverseTraversePrevPageButtonId := ToConvinientVerticalId(Index(prevPageButtonId), sh.GetDevice())
//	err = sh.GetDevice().SetText(reverseTraverseNextPageButtonId, ">")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = sh.GetDevice().SetText(reverseTraversePageNumberButtonId, strconv.Itoa(int(sh.page)))
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = sh.GetDevice().SetText(reverseTraversePrevPageButtonId, "<")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	pageIndices := sh.buttonIdToText[Page(page)]
//	for buttonIndex, text := range pageIndices {
//		convenientVerticalId := ToConvinientVerticalId(buttonIndex, sh.GetDevice())
//		err = sh.GetDevice().SetText(convenientVerticalId, text)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}

//func (sh *StreamdeckHandler) SwitchPage(page Page) {
//	sh.device.Clear()
//	sh.page = page
//	err := SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns-1, ">")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns*2-1, strconv.Itoa(int(sh.page)))
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = SetStreamdeckButtonText(sh.GetDevice(), sh.device.device.Columns*3-1, "<")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for buttonId, text := range sh.buttonIdToText {
//		reverseButtonId := ToConvinientVerticalId(buttonId, sh.GetDevice())
//		buttonIdPerPage := reverseButtonId - (sh.page-1)*sh.numOfButtons
//		if reverseButtonId >= 0 && buttonIdPerPage < sh.numOfButtons {
//			//newId := sh.TraverseButtonId(buttonIdPerPage)
//			err := SetStreamdeckButtonText(sh.device, uint8(TraverseButtonId(buttonIdPerPage, sh.GetDevice())), text)
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}
//}

func (sh *StreamdeckHandler) AddOnPressHandler(page Page, buttonId Index, handler func() error) {
	if sh.streamDeckButtonIdToOnPressHandler[page] == nil {
		sh.streamDeckButtonIdToOnPressHandler[page] = make(map[Index]func() error)
	}
	sh.streamDeckButtonIdToOnPressHandler[page][buttonId] = handler
}

func (sh *StreamdeckHandler) AddOnReleaseHandler(page Page, buttonId Index, handler func() error) {
	if sh.streamDeckButtonIdToOnReleaseHandler[page] == nil {
		sh.streamDeckButtonIdToOnReleaseHandler[page] = make(map[Index]func() error)
	}
	sh.streamDeckButtonIdToOnReleaseHandler[page][buttonId] = handler
}

func (sh *StreamdeckHandler) GetOnPressHandler(page Page, buttonId Index) (func() error, bool) {
	if sh.streamDeckButtonIdToOnPressHandler[page] == nil {
		return nil, false
	}

	handler, ok := sh.streamDeckButtonIdToOnPressHandler[page][buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetOnReleaseHandler(page Page, buttonId Index) (func() error, bool) {
	if sh.streamDeckButtonIdToOnReleaseHandler[page] == nil {
		return nil, false
	}
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[page][buttonId]
	return handler, ok
}

func (sh *StreamdeckHandler) GetDevice() DeviceWrapper {
	return sh.device
}

func (sh *StreamdeckHandler) StartListenAsync() error {
	SwitchPage(sh, 0)
	cols := sh.device.GetColumns()
	nextPageButtonId := cols - 1
	prevPageButtonId := cols*3 - 1

	go func() {
		keys, err := sh.device.device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				fmt.Printf("Key pressed index reverseTransverse %v, is pressed %v\n", ToConvinientVerticalId(Index(key.Index), sh.GetDevice()), key.Pressed)
				if Index(key.Index) == Index(nextPageButtonId) && key.Pressed {
					sh.page++
					SwitchPage(sh, sh.page)
					continue
				} else if Index(key.Index) == Index(prevPageButtonId) && key.Pressed {
					sh.page = Page(int(math.Max(float64(0), float64(sh.page-1))))
					SwitchPage(sh, sh.page)
					continue
				}
				if key.Pressed {
					if handler, ok := sh.GetOnPressHandler(sh.page, ToConvinientVerticalId(Index(key.Index), sh.GetDevice())); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					if handler, ok := sh.GetOnReleaseHandler(sh.page, ToConvinientVerticalId(Index(key.Index), sh.GetDevice())); ok {
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
