package streamdeck

import (
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

func (s *StreamdeckDeviceWrapper) Clear() error {
	return s.device.Clear()
}

func (s *StreamdeckDeviceWrapper) Close() interface{} {
	return s.device.Close()
}

func (s *StreamdeckDeviceWrapper) SetImage(index uint8, img image.Image) error {
	return s.device.SetImage(index, img)
}

func (s *StreamdeckDeviceWrapper) GetPixels() int {
	return int(s.device.Pixels)
}

func (s *StreamdeckDeviceWrapper) DPI() uint {
	return s.device.DPI
}
