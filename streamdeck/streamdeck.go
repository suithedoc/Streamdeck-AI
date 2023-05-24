package streamdeck

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/muesli/streamdeck"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
)

func ftContext(img *image.RGBA, ttf *truetype.Font, dpi uint, fontsize float64) *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(float64(dpi))
	c.SetFont(ttf)
	c.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 0}))
	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fontsize)

	return c
}

type Layout struct {
	frames []image.Rectangle
	size   int
	margin int
	height int
}

// FormatLayout returns a layout that is formatted according to frameReps.
func (l *Layout) FormatLayout(frameReps []string, frameCount int) []image.Rectangle {
	if frameCount < 1 {
		frameCount = 1
	}
	for i := 0; i < frameCount; i++ {
		if len(frameReps) < i+1 {
			frame := l.defaultFrame(frameCount, i)
			l.frames = append(l.frames, frame)
			continue
		}

		frame, err := formatFrame(frameReps[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, "using default frame:", err)
			frame = l.defaultFrame(frameCount, i)
		}
		l.frames = append(l.frames, frame)
	}
	return l.frames
}

// Returns the Rectangle representing the index-th horizontal cell.
func (l *Layout) defaultFrame(cells int, index int) image.Rectangle {
	lower := l.margin + (l.height/cells)*index
	upper := l.margin + (l.height/cells)*(index+1)
	return image.Rect(0, lower, l.size, upper)
}

// Converts the string representation of a rectangle into a image.Rectangle.
func formatFrame(layout string) (image.Rectangle, error) {
	split := strings.Split(layout, "+")
	if len(split) < 2 {
		return image.Rectangle{}, fmt.Errorf("invalid rectangle format")
	}
	position, errP := formatCoord(split[0])
	if errP != nil {
		return image.Rectangle{}, errP
	}
	extent, errE := formatCoord(split[1])
	if errE != nil {
		return image.Rectangle{}, errE
	}

	return image.Rectangle{position, position.Add(extent)}, nil
}

// Converts the string representation of a point into a image.Point.
func formatCoord(coords string) (image.Point, error) {
	split := strings.Split(coords, "x")
	if len(split) < 2 {
		return image.Point{}, fmt.Errorf("invalid point format")
	}
	posX, errX := strconv.Atoi(split[0])
	posY, errY := strconv.Atoi(split[1])
	if errX != nil || errY != nil {
		return image.Point{}, fmt.Errorf("invalid point format")
	}
	return image.Pt(posX, posY), nil
}

func NewLayout(size int) *Layout {
	margin := size / 18
	height := size - (margin * 2)

	return &Layout{
		size:   size,
		margin: margin,
		height: height,
	}
}

func InitStreamdeckDevice() (*StreamdeckDeviceWrapper, error) {
	devs, err := streamdeck.Devices()
	if err != nil {
		return nil, fmt.Errorf("no Stream Deck devices found: %s", err)
	}
	if len(devs) == 0 {
		return nil, fmt.Errorf("no Stream Deck devices found")
	}
	d := devs[0]
	if err := d.Open(); err != nil {
		return nil, fmt.Errorf("can't open device: %s", err)
	}
	return NewStreamdeckDeviceWrapper(&d), nil
}

func SetStreamdeckButtonText(device DeviceWrapper, index uint8, text string) error {
	img, err := CreateTextImage(device, text)
	if err != nil {
		return err
	}
	return device.SetImage(index, img)
}

func CreateTextImage(device DeviceWrapper, text string) (image.Image, error) {
	size := int(device.GetPixels())
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	ttfFont, err := LoadFont("Vera.ttf")
	if err != nil {
		return nil, err
	}
	DrawString(img,
		image.Rect(0, 0, size, size),
		ttfFont,
		text,
		device.DPI(),
		-1,
		color.RGBA{255, 255, 255, 255},
		image.Pt(-1, -1))

	return img, nil
}

func DrawString(img *image.RGBA, bounds image.Rectangle, ttf *truetype.Font, texts string, dpi uint, fontsize float64, color color.Color, pt image.Point) {

	splittedTexts := strings.Split(texts, "\n")
	c := ftContext(img, ttf, dpi, fontsize)

	if fontsize <= 0 {
		// pick biggest available height to fit the string
		fontsize, _ = maxPointSize(texts,
			ftContext(img, ttf, dpi, fontsize), dpi,
			bounds.Dx(), bounds.Dy())
		c.SetFontSize(fontsize)
	}

	if pt.X < 0 {
		// center horizontally
		extent, _ := ftContext(img, ttf, dpi, fontsize).DrawString(splittedTexts[0], freetype.Pt(0, 0))
		actwidth := extent.X.Floor()
		xcenter := float64(bounds.Dx())/2.0 - (float64(actwidth) / 2.0)
		pt = image.Pt(bounds.Min.X+int(xcenter), pt.Y)
	}
	if pt.Y < 0 {
		// center vertically
		actheight := c.PointToFixed(fontsize).Round()
		ycenter := float64(bounds.Dy()/2.0) + (float64(actheight) / 2.6)
		pt = image.Pt(pt.X, bounds.Min.Y+int(ycenter))
	}

	c.SetSrc(image.NewUniform(color))
	for _, text := range splittedTexts {
		if _, err := c.DrawString(text, freetype.Pt(pt.X, pt.Y)); err != nil {
			fmt.Fprintf(os.Stderr, "Can't render string: %s\n", err)
			return
		}
		pt.Y += int(fontsize * 1.2) // add some spacing between lines
	}
}
