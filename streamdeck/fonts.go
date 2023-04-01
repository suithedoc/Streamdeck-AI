package streamdeck

import (
	"github.com/flopp/go-findfont"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
)

func maxPointSize(text string, c *freetype.Context, dpi uint, width, height int) (float64, int) {
	fontsize := float64(height<<6) / float64(dpi) / (64.0 / 72.0)
	fontsize++
	var actwidth int
	for actwidth = width + 1; actwidth > width; fontsize-- {
		c.SetFontSize(fontsize)

		textExtent, err := c.DrawString(text, freetype.Pt(0, 0))
		if err != nil {
			return 0, 0
		}

		actwidth = textExtent.X.Round()
	}

	return fontsize, actwidth
}

func LoadFont(name string) (*truetype.Font, error) {
	fontPath, err := findfont.Find(name)
	if err != nil {
		return nil, err
	}

	ttf, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(ttf)
}
