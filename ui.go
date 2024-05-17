package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func updateTextBlocks(grid *tview.Grid, textBlocks []TextBlock, inputFields *[]*tview.TextArea) {
	grid.Clear()
	*inputFields = nil
	for i := range textBlocks {
		//
		inputField := tview.NewTextArea().
			SetText(textBlocks[i].Text, true).
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
		inputField.SetChangedFunc(func() {
			textBlocks[index].Text = inputField.GetText()
		})

		*inputFields = append(*inputFields, inputField)

	}
}
