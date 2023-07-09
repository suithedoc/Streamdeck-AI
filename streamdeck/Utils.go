package streamdeck

import (
	"log"
	"strconv"
)

// Convert the convinient numbering, from top to bottom, left to right, to the actual button id.
func TraverseButtonId(buttonId Index, device DeviceWrapper) Index {

	// Convert the vertical-first index into a row and column.
	row := int(buttonId) % device.GetRows()
	col := int(buttonId) / device.GetRows()

	// Convert the row and column into a horizontal-first index.
	newButtonId := row*device.GetColumns() + col

	return Index(newButtonId)
}

// Convert the actual button id which is from left to right, top to bottom, to the convinient numbering.
func ToConvinientVerticalId(buttonId Index, device DeviceWrapper) Index {

	// Convert the horizontal-first index into a row and column.
	row := int(buttonId) / device.GetColumns()
	col := int(buttonId) % device.GetColumns()

	// Convert the row and column into a vertical-first index.
	newButtonId := col*device.GetRows() + row

	return Index(newButtonId)
}

func SwitchPage(sh IStreamdeckHandler, page Page) {
	err := sh.GetDevice().Clear()
	if err != nil {
		log.Fatal(err)
		return //log.Fatal(err)
	}
	sh.SetPage(page)
	//err := SetStreamdeckButtonText(sh.GetDevice(), uint8(sh.device.GetColumns()-1), ">")
	cols := sh.GetDevice().GetColumns()
	nextPageButtonId := cols - 1
	pageNumberButtonId := cols*2 - 1
	prevPageButtonId := cols*3 - 1

	reverseTraverseNextPageButtonId := ToConvinientVerticalId(Index(nextPageButtonId), sh.GetDevice())
	reverseTraversePageNumberButtonId := ToConvinientVerticalId(Index(pageNumberButtonId), sh.GetDevice())
	reverseTraversePrevPageButtonId := ToConvinientVerticalId(Index(prevPageButtonId), sh.GetDevice())
	err = sh.GetDevice().SetText(reverseTraverseNextPageButtonId, ">")
	if err != nil {
		log.Fatal(err)
	}
	err = sh.GetDevice().SetText(reverseTraversePageNumberButtonId, strconv.Itoa(int(sh.GetPage())))
	if err != nil {
		log.Fatal(err)
	}
	err = sh.GetDevice().SetText(reverseTraversePrevPageButtonId, "<")
	if err != nil {
		log.Fatal(err)
	}

	pageIndices := sh.GetButtonIndexToText()[page]
	for buttonIndex, text := range pageIndices {
		convenientVerticalId := ToConvinientVerticalId(buttonIndex, sh.GetDevice())
		err = sh.GetDevice().SetText(convenientVerticalId, text)
		if err != nil {
			log.Fatal(err)
		}
	}
}
