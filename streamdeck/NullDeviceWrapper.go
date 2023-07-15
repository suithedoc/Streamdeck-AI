package streamdeck

import "image"

type NullDeviceWrapper struct{}

func (ndw NullDeviceWrapper) Clear() error {
	return nil
}

func (ndw NullDeviceWrapper) Close() interface{} {
	return nil
}

func (ndw NullDeviceWrapper) SetImage(index uint8, img image.Image) error {
	return nil
}

func (ndw NullDeviceWrapper) SetText(index Index, text string) error {
	return nil
}

func (ndw NullDeviceWrapper) GetPixels() int {
	return 0
}

func (ndw NullDeviceWrapper) DPI() uint {
	return 0
}

func (ndw NullDeviceWrapper) GetRows() int {
	return 0
}

func (ndw NullDeviceWrapper) GetColumns() int {
	return 0
}
