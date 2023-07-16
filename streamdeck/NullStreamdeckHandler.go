package streamdeck

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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
	return NullDeviceWrapper{}
}
func (nsh NullStreamdeckHandler) StartListenAsync() error {
	// Create a channel to receive OS signals.
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal. This will block until a signal is received.
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	fmt.Println("Exiting due to received signal.")

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
