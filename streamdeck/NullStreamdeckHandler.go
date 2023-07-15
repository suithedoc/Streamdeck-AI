package streamdeck

type NullStreamdeckHandler struct{}

func (nsh NullStreamdeckHandler) AddOnPressHandler(page Page, buttonId Index, handler func() error) {}
func (nsh NullStreamdeckHandler) AddOnReleaseHandler(page Page, buttonId Index, handler func() error) {
}
func (nsh NullStreamdeckHandler) GetOnPressHandler(page Page, buttonId Index) (func() error, bool) {
	return nil, false
}
func (nsh NullStreamdeckHandler) GetOnReleaseHandler(page Page, buttonId Index) (func() error, bool) {
	return nil, false
}
func (nsh NullStreamdeckHandler) GetDevice() DeviceWrapper {
	return nil
}
func (nsh NullStreamdeckHandler) StartListenAsync() error {
	return nil
}
func (nsh NullStreamdeckHandler) AddButtonText(page Page, buttonId Index, text string) error {
	return nil
}
func (nsh NullStreamdeckHandler) GetButtonIndexToText() map[Page]map[Index]string {
	return nil
}
func (nsh NullStreamdeckHandler) GetPage() Page {
	return 0
}
func (nsh NullStreamdeckHandler) SetPage(page Page) {}
