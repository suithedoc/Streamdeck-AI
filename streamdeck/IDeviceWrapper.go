package streamdeck

import (
	"image"
)

type DeviceWrapper interface {
	Clear() error
	Close() interface{}
	SetImage(index uint8, img image.Image) error
	GetPixels() int
	DPI() uint
}
