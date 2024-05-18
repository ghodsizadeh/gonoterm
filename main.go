package main

import (
	"flag"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const logFileName = "app.log"

func main() {
	app := tview.NewApplication()

	enableLogging := flag.Bool("log", false, "Enable logging")
	flag.Parse()
	var logFile *os.File
	if *enableLogging {
		logFile, _ = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// Set log output to the file
		log.SetOutput(logFile)
		defer logFile.Close()

	}
	log.Println("Starting gonoterm")
	textBlocks, err := loadFromFile()
	if err != nil {
		log.Fatalf("Error loading text blocks: %v", err)
	}

	if len(textBlocks) == 0 {
		textBlocks = append(textBlocks, TextBlock{Type: Text, Text: "Hello Gonoterm!"})
	}

	var inputFields []*tview.TextArea
	textAreaGrid := tview.NewGrid().
		SetRows(0).
		SetColumns(0)

	updateTextBlocks(textAreaGrid, textBlocks, &inputFields)

	helpInfo := tview.NewTextView().
		SetText("Press Ctrl-N to add block, Ctrl-R to remove block, Ctrl-S to save, Ctrl-D/H to switch between inputs")
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()

	mainView := tview.NewGrid().
		SetRows(0, 1, 1).
		AddItem(textAreaGrid, 0, 0, 1, 1, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 2, 0, 1, 1, 0, 0, false)

	pages.AddAndSwitchToPage("main", mainView, true)

	if len(inputFields) > 0 {
		app.SetFocus(inputFields[0])
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			textBlocks = append(textBlocks, TextBlock{Type: Text, Text: ""})
			updateTextBlocks(textAreaGrid, textBlocks, &inputFields)
			if len(inputFields) > 0 {
				app.SetFocus(inputFields[len(inputFields)-1])
			}
			return nil

		case tcell.KeyCtrlR:
			if len(textBlocks) > 0 {
				textBlocks = textBlocks[:len(textBlocks)-1]
				updateTextBlocks(textAreaGrid, textBlocks, &inputFields)
				if len(inputFields) > 0 {
					app.SetFocus(inputFields[len(inputFields)-1])
				}
			}
			return nil
		case tcell.KeyCtrlS:
			if err := saveToFile(textBlocks); err != nil {
				log.Fatalf("Error saving text blocks: %v", err)
			}
			return nil
		case tcell.KeyCtrlV:
			handleTypeSelection(app, pages, textBlocks, &inputFields, textAreaGrid)

			return nil

		case tcell.KeyCtrlD:
			currentFocus := app.GetFocus()
			for i, inputField := range inputFields {
				if currentFocus == inputField && i > 0 {
					app.SetFocus(inputFields[i-1])
					break
				}
			}
			return nil
		case tcell.KeyCtrlH:
			currentFocus := app.GetFocus()
			for i, inputField := range inputFields {
				if currentFocus == inputField && i < len(inputFields)-1 {
					app.SetFocus(inputFields[i+1])
					break
				}
			}
			return nil
		}
		saveToFile(textBlocks)
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
