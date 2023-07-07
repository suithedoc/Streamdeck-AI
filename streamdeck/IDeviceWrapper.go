package streamdeck

import (
	"image"
)

type DeviceWrapper interface {
	Clear() error
	Close() interface{}
	SetImage(index uint8, img image.Image) error
	SetText(index int, text string) error
	GetPixels() int
	DPI() uint
	GetRows() int
	GetColumns() int
}
