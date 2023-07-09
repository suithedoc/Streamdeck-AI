package streamdeck

import (
	"fmt"
	"math"
)

type Index int
type Page int

type UiStreamdeckHandler struct {
	IStreamdeckHandler
	device                               *UiDeviceWrapper
	streamDeckButtonIdToOnPressHandler   map[Page]map[Index]func() error
	streamDeckButtonIdToOnReleaseHandler map[Page]map[Index]func() error
	page                                 Page
	buttonIdToText                       map[Page]map[Index]string
	numOfButtons                         int
}

func NewUiStreamdeckHandler() (IStreamdeckHandler, error) {
	device, err := InitUiStreamdeck()
	if err != nil {
		return nil, err
	}
	return &UiStreamdeckHandler{
		device:                               device,
		streamDeckButtonIdToOnPressHandler:   make(map[Page]map[Index]func() error),
		streamDeckButtonIdToOnReleaseHandler: make(map[Page]map[Index]func() error),
		page:                                 1,
		buttonIdToText:                       make(map[Page]map[Index]string),
		numOfButtons:                         int(device.GetRows() * (device.GetColumns() - 1)),
	}, nil
}

func (sh *UiStreamdeckHandler) GetButtonIndexToText() map[Page]map[Index]string {
	return sh.buttonIdToText
}

func (sh *UiStreamdeckHandler) GetPage() Page {
	return sh.page
}

func (sh *UiStreamdeckHandler) SetPage(page Page) {
	sh.page = page
}

func (sh *UiStreamdeckHandler) AddButtonText(page Page, buttonId Index, text string) error {
	//sh.device.buttons[TraverseButtonId(buttonId, sh.GetDevice())].SetText(text)
	if sh.buttonIdToText[page] == nil {
		sh.buttonIdToText[page] = make(map[Index]string)
	}
	sh.buttonIdToText[page][TraverseButtonId(buttonId, sh.GetDevice())] = text
	return nil
}

func (sh *UiStreamdeckHandler) GetDevice() DeviceWrapper {
	return sh.device
}

func (sh *UiStreamdeckHandler) StartListenAsync() error {
	SwitchPage(sh, 0)
	cols := sh.device.GetColumns()
	nextPageButtonId := cols - 1
	prevPageButtonId := cols*3 - 1

	//reverseTraverseNextPageButtonId := ToConvinientVerticalId(nextPageButtonId, sh.GetDevice())
	//reverseTraversePrevPageButtonId := ToConvinientVerticalId(prevPageButtonId, sh.GetDevice())
	go func() {
		eventChan, err := sh.device.ReadEvents()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			event := <-eventChan
			if event.ButtonIndex == Index(nextPageButtonId) {
				sh.page++
				SwitchPage(sh, sh.page)
				continue
			} else if event.ButtonIndex == Index(prevPageButtonId) {
				sh.page = Page(int(math.Max(float64(0), float64(sh.page-1))))
				SwitchPage(sh, sh.page)
				continue
			}
			if event.IsPressed {
				handler, ok := sh.GetOnPressHandler(event.Page, ToConvinientVerticalId(event.ButtonIndex, sh.GetDevice()))
				if ok {
					err := handler()
					if err != nil {
						fmt.Println(err)
					}
				}
			} else if event.IsReleased {
				handler, ok := sh.GetOnReleaseHandler(event.Page, ToConvinientVerticalId(event.ButtonIndex, sh.GetDevice()))
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

func (sh *UiStreamdeckHandler) AddOnPressHandler(page Page, buttonId Index, handler func() error) {
	if sh.streamDeckButtonIdToOnPressHandler[page] == nil {
		sh.streamDeckButtonIdToOnPressHandler[page] = make(map[Index]func() error)
	}
	sh.streamDeckButtonIdToOnPressHandler[page][buttonId] = handler
}

func (sh *UiStreamdeckHandler) AddOnReleaseHandler(page Page, buttonId Index, handler func() error) {
	if sh.streamDeckButtonIdToOnReleaseHandler[page] == nil {
		sh.streamDeckButtonIdToOnReleaseHandler[page] = make(map[Index]func() error)
	}
	sh.streamDeckButtonIdToOnReleaseHandler[page][buttonId] = handler
}

func (sh *UiStreamdeckHandler) GetOnPressHandler(page Page, buttonId Index) (func() error, bool) {
	if sh.streamDeckButtonIdToOnPressHandler[page] == nil {
		return nil, false
	}
	handler, ok := sh.streamDeckButtonIdToOnPressHandler[page][buttonId]
	return handler, ok
}

func (sh *UiStreamdeckHandler) GetOnReleaseHandler(page Page, buttonId Index) (func() error, bool) {
	if sh.streamDeckButtonIdToOnReleaseHandler[page] == nil {
		return nil, false
	}
	handler, ok := sh.streamDeckButtonIdToOnReleaseHandler[page][buttonId]
	return handler, ok
}

//func (sh *UiStreamdeckHandler) SwitchPage(page Page) {
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
//	pageIndices := sh.buttonIdToText[page]
//	for buttonIndex, text := range pageIndices {
//		convenientVerticalId := ToConvinientVerticalId(buttonIndex, sh.GetDevice())
//		err = sh.GetDevice().SetText(convenientVerticalId, text)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}

func (sh *UiStreamdeckHandler) TraverseButtonId(buttonId int) int {
	return buttonId
}

func (sh *UiStreamdeckHandler) ReverseTraverseButtonId(buttonId int) int {
	return buttonId
}
