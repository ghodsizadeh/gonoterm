package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func updateTextBlocks(grid *tview.Grid, textBlocks []TextBlock, inputFields *[]*tview.InputField) {
	grid.Clear()
	*inputFields = nil
	for i := range textBlocks {
		//
		inputField := tview.NewInputField().
			SetText(textBlocks[i].Text).
			SetPlaceholder("Enter text here...")
		// SetLabel(fmt.Sprintf("Block %d: ", i+1)).
		// SetFieldWidth(30)
		inputField.SetBorder(true)
		inputField.SetTitle(fmt.Sprintf("[green]Block %d %s", i+1, textBlocks[i].Type.String()))
		// change the title color bas

		grid.AddItem(inputField, i, 0, 1, 1, 0, 0, true)

		// Capture the loop variable
		index := i

		// Update block text on change
		inputField.SetChangedFunc(func(text string) {
			textBlocks[index].Text = text
		})

		*inputFields = append(*inputFields, inputField)

	}
}
