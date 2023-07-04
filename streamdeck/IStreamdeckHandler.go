package streamdeck

type IStreamdeckHandler interface {
	AddOnPressHandler(buttonId int, handler func() error)
	AddOnReleaseHandler(buttonId int, handler func() error)
	GetOnPressHandler(buttonId int) (func() error, bool)
	GetOnReleaseHandler(buttonId int) (func() error, bool)
	SwitchPage(page int)
	GetDevice() DeviceWrapper
	StartListenAsync() error

	AddButtonText(buttonId int, text string) error
	TraverseButtonId(buttonId int) int
	ReverseTraverseButtonId(buttonId int) int
	//StartAsync()
}
