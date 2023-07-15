package streamdeck

type IStreamdeckHandler interface {
	AddOnPressHandler(page Page, buttonId Index, handler func() error)
	AddOnReleaseHandler(page Page, buttonId Index, handler func() error)
	GetOnPressHandler(page Page, buttonId Index) (func() error, bool)
	GetOnReleaseHandler(page Page, buttonId Index) (func() error, bool)
	GetDevice() DeviceWrapper
	StartListenAsync() error

	AddButtonText(page Page, buttonId Index, text string) error
	GetButtonIndexToText() map[Page]map[Index]string
	GetPage() Page
	SetPage(page Page)
}
