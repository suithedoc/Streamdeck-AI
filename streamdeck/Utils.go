package streamdeck

// Convert the convinient numbering, from top to bottom, left to right, to the actual button id.
func TraverseButtonId(buttonId int, device DeviceWrapper) int {

	// Convert the vertical-first index into a row and column.
	row := buttonId % device.GetRows()
	col := buttonId / device.GetRows()

	// Convert the row and column into a horizontal-first index.
	newButtonId := row*device.GetColumns() + col

	return newButtonId
}

// Convert the actual button id which is from left to right, top to bottom, to the convinient numbering.
func ReverseTraverseButtonId(buttonId int, device DeviceWrapper) int {

	// Convert the horizontal-first index into a row and column.
	row := buttonId / device.GetColumns()
	col := buttonId % device.GetColumns()

	// Convert the row and column into a vertical-first index.
	newButtonId := col*device.GetRows() + row

	return newButtonId
}
