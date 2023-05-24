package streamdeck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
)

type UiDeviceWrapper struct {
	DeviceWrapper
	window   fyne.Window
	app      fyne.App
	textView *widget.Entry
}

func (dw *UiDeviceWrapper) Clear() error {
	return nil
}

func (dw *UiDeviceWrapper) Close() interface{} {
	return nil
}

func (dw *UiDeviceWrapper) SetImage(index uint8, img image.Image) error {
	return nil
}

func (dw *UiDeviceWrapper) GetPixels() int {
	return 0
}

func (dw *UiDeviceWrapper) DPI() uint {
	return 0
}

type UiStreamdeckHandler struct {
	IStreamdeckHandler
	device                               *UiDeviceWrapper
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
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
	}, nil
}

func (sh *UiStreamdeckHandler) GetDevice() DeviceWrapper {
	return sh.device
}

func (sh *UiStreamdeckHandler) StartListenAsync() error {
	sh.device.window.ShowAndRun()
	return nil
}

func (sh *UiStreamdeckHandler) AddOnPressHandler(buttonId int, handler func() error) {

}

func (sh *UiStreamdeckHandler) AddOnReleaseHandler(buttonId int, handler func() error) {

}

func (sh *UiStreamdeckHandler) GetOnPressHandler(buttonId int) (func() error, bool) {
	return nil, false
}

func (sh *UiStreamdeckHandler) GetOnReleaseHandler(buttonId int) (func() error, bool) {
	return nil, false
}

//var textView *widget.Entry

type customButton struct {
	widget.Button
	isActive bool
}

func InitUiStreamdeck() (*UiDeviceWrapper, error) {
	uiWrapper := &UiDeviceWrapper{}
	uiWrapper.app = app.New()
	//myApp := app.New()
	uiWrapper.window = uiWrapper.app.NewWindow("Numpad")

	uiWrapper.textView = widget.NewMultiLineEntry()
	uiWrapper.textView.SetPlaceHolder("Enter numbers here")

	grid := container.NewGridWithColumns(5,
		uiWrapper.makeButton("0"), uiWrapper.makeButton("3"), uiWrapper.makeButton("6"), uiWrapper.makeButton("9"), uiWrapper.makeButton("12"),
		uiWrapper.makeButton("1"), uiWrapper.makeButton("4"), uiWrapper.makeButton("7"), uiWrapper.makeButton("10"), uiWrapper.makeButton("13"),
		uiWrapper.makeButton("2"), uiWrapper.makeButton("5"), uiWrapper.makeButton("8"), uiWrapper.makeButton("11"), uiWrapper.makeButton("14"),
	)

	content := container.NewVBox(grid, uiWrapper.textView)
	uiWrapper.window.SetContent(content)
	uiWrapper.window.ShowAndRun()
	return uiWrapper, nil
}

func (dw *UiDeviceWrapper) makeButton(label string) *customButton {
	button := &customButton{}
	button.Text = label
	button.ExtendBaseWidget(button)

	button.OnTapped = func() {
		if button.isActive {
			dw.textView.SetText(removeLastChar(dw.textView.Text))
			button.isActive = false
		} else {
			dw.textView.SetText(dw.textView.Text + label)
			button.isActive = true
		}
	}

	return button
}

func removeLastChar(str string) string {
	if len(str) > 0 {
		return str[:len(str)-1]
	}
	return ""
}
