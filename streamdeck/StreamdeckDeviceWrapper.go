package streamdeck

import (
	"fmt"
	"github.com/muesli/streamdeck"
	"image"
)

type StreamdeckDeviceWrapper struct {
	DeviceWrapper
	device *streamdeck.Device
}

func NewStreamdeckDeviceWrapper(device *streamdeck.Device) *StreamdeckDeviceWrapper {
	return &StreamdeckDeviceWrapper{
		device: device,
	}
}

func (s *StreamdeckDeviceWrapper) GetRows() int {
	return int(s.device.Rows)
}

func (s *StreamdeckDeviceWrapper) GetColumns() int {
	return int(s.device.Columns)
}

func (s *StreamdeckDeviceWrapper) Clear() error {
	return s.device.Clear()
}

func (s *StreamdeckDeviceWrapper) Close() interface{} {
	return s.device.Close()
}

func (s *StreamdeckDeviceWrapper) SetImage(index uint8, img image.Image) error {
	return s.device.SetImage(index, img)
}

func (dw *StreamdeckDeviceWrapper) SetText(index int, text string) error {
	img, err := CreateTextImage(dw, text)
	if err != nil {
		return err
	}
	err = dw.SetImage(uint8(index), img)
	if err != nil {
		fmt.Println("setting image: ", err)
		return err
	}
	return nil
}

func (s *StreamdeckDeviceWrapper) GetPixels() int {
	return int(s.device.Pixels)
}

func (s *StreamdeckDeviceWrapper) DPI() uint {
	return s.device.DPI
}
