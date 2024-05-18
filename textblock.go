package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

type BlockType int

const (
	Math BlockType = iota
	Text
	Code
)

type TextBlock struct {
	Type BlockType
	Text string
}

func (b BlockType) String() string {
	switch b {
	case Math:
		return "Math"
	case Text:
		return "Text"
	case Code:
		return "Code"
	default:
		return "Unknown"
	}
}

// const fileName = "textblocks.json"
var (
	namespace = flag.String("namespace", "default-gonoterm", "Namespace for the textblocks")
	local     = flag.Bool("local", false, "Use local storage")
)

// Save textBlocks to a file
func saveToFile(textBlocks []TextBlock) error {
	data, err := json.MarshalIndent(textBlocks, "", "  ")
	if err != nil {
		return err
	}

	fileName := getFileName()
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	err = os.WriteFile(fileName, data, 0644)
	log.Println("Saved to file: ", fileName)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	return nil
}

// Load textBlocks from a file
func loadFromFile() ([]TextBlock, error) {

	fileName := getFileName()
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []TextBlock{}, nil
		}
		return nil, err
	}
	var textBlocks []TextBlock
	err = json.Unmarshal(data, &textBlocks)
	return textBlocks, err
}

func getFileName() string {
	if *local {
		return "gonoterm.json"
	}
	// Use a universal file path
	// return "/universal/path/" + *namespace + ".json"
	// get  a path in home directory config
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting config directory: %v", err)
	}
	path := filepath.Join(configDir, "gonoterm", *namespace+".json")
	log.Println("Path: ", path)
	return path
}

func handleTypeSelection(app *tview.Application, pages *tview.Pages, textBlocks []TextBlock, inputFields *[]*tview.TextArea,
	textAreaGrid *tview.Grid) {
	currentFocus := app.GetFocus()
	for i, inputField := range *inputFields {
		if currentFocus == inputField {
			modal := tview.NewModal().
				SetText("Select the type of the block").
				AddButtons([]string{"Text", "Code", "Math"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					switch buttonLabel {
					case "Text":
						textBlocks[i].Type = Text
					case "Code":
						textBlocks[i].Type = Code
					case "Math":
						textBlocks[i].Type = Math
					}
					updateTextBlocks(textAreaGrid, textBlocks, inputFields)

					pages.RemovePage("modal")
				})
			pages.AddPage("modal", modal, true, true)
			app.SetFocus(modal)

			break
		}
	}
}
