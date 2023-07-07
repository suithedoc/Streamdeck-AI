package streamdeck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	"strconv"
)

type ButtonEvent struct {
	ButtonId   int
	IsPressed  bool
	IsReleased bool
}

type UiDeviceWrapper struct {
	DeviceWrapper
	window      fyne.Window
	app         fyne.App
	textView    *widget.Entry
	buttons     []*CustomButton
	buttonEvent chan ButtonEvent
}

func (dw *UiDeviceWrapper) ReadEvents() (chan ButtonEvent, error) {
	return dw.buttonEvent, nil
}

func (dw *UiDeviceWrapper) Clear() error {
	for _, button := range dw.buttons {
		button.SetText("")
	}

	return nil
}

func (dw *UiDeviceWrapper) SetText(index int, text string) error {
	fmt.Printf("SetText: %d, %s\n", index, text)
	dw.buttons[index].SetText(text)
	return nil
}

func (dw *UiDeviceWrapper) Close() interface{} {
	return nil
}

func (dw *UiDeviceWrapper) SetImage(index uint8, img image.Image) error {
	return nil
}

func (dw *UiDeviceWrapper) GetPixels() int {
	return int(dw.buttons[0].Size().Height * dw.buttons[0].Size().Width)
}

func (dw *UiDeviceWrapper) DPI() uint {
	return 0
}

//var textView *widget.Entry

type CustomButton struct {
	widget.Button
	isActive bool
}

func InitUiStreamdeck() (*UiDeviceWrapper, error) {
	uiWrapper := &UiDeviceWrapper{}
	uiWrapper.buttonEvent = make(chan ButtonEvent)
	uiWrapper.app = app.New()
	//myApp := app.New()
	uiWrapper.window = uiWrapper.app.NewWindow("Numpad")

	uiWrapper.textView = widget.NewMultiLineEntry()
	uiWrapper.textView.SetPlaceHolder("Enter numbers here")

	grid := container.NewGridWithColumns(5,
		uiWrapper.makeButton(0), uiWrapper.makeButton(1), uiWrapper.makeButton(2), uiWrapper.makeButton(3), uiWrapper.makeButton(4),
		uiWrapper.makeButton(5), uiWrapper.makeButton(6), uiWrapper.makeButton(7), uiWrapper.makeButton(8), uiWrapper.makeButton(9),
		uiWrapper.makeButton(10), uiWrapper.makeButton(11), uiWrapper.makeButton(12), uiWrapper.makeButton(13), uiWrapper.makeButton(14),
	)

	content := container.NewVBox(grid, uiWrapper.textView)
	uiWrapper.window.SetContent(content)
	//uiWrapper.window.ShowAndRun()
	return uiWrapper, nil
}

func (dw *UiDeviceWrapper) makeButton(buttonId int) *CustomButton {
	button := &CustomButton{}
	button.Text = strconv.Itoa(buttonId)
	button.ExtendBaseWidget(button)

	button.OnTapped = func() {
		if button.isActive {
			dw.textView.SetText(removeLastChar(dw.textView.Text))
			button.isActive = false
			dw.buttonEvent <- ButtonEvent{ButtonId: buttonId, IsPressed: false, IsReleased: true}
		} else {
			dw.textView.SetText(dw.textView.Text + strconv.Itoa(buttonId))
			button.isActive = true
			dw.buttonEvent <- ButtonEvent{ButtonId: buttonId, IsPressed: true, IsReleased: false}
		}
	}
	dw.buttons = append(dw.buttons, button)
	return button
}

func (dw *UiDeviceWrapper) GetRows() int {
	return 3
}

func (dw *UiDeviceWrapper) GetColumns() int {
	return 5
}

func removeLastChar(str string) string {
	if len(str) > 0 {
		return str[:len(str)-1]
	}
	return ""
}
