package streamdeck

import (
	"fmt"
	"log"
)

type IStreamdeckHandler interface {
	AddOnPressHandler(buttonId int, handler func() error)
	AddOnReleaseHandler(buttonId int, handler func() error)
	GetOnPressHandler(buttonId int) (func() error, bool)
	GetOnReleaseHandler(buttonId int) (func() error, bool)
	GetDevice() DeviceWrapper
	StartListenAsync() error
}

type StreamdeckHandler struct {
	IStreamdeckHandler
	streamDeckButtonIdToOnPressHandler   map[int]func() error
	streamDeckButtonIdToOnReleaseHandler map[int]func() error
	device                               *StreamdeckDeviceWrapper
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
	}, nil
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
	go func() {
		keys, err := sh.device.device.ReadKeys()
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case key := <-keys:
				fmt.Printf("Key pressed index %v, is pressed %v\n", key.Index, key.Pressed)
				if key.Pressed {
					if handler, ok := sh.GetOnPressHandler(int(key.Index)); ok {
						err := handler()
						if err != nil {
							log.Fatal(err)
						}
					}
				} else {
					if handler, ok := sh.GetOnReleaseHandler(int(key.Index)); ok {
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
