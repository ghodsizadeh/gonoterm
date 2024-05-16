package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

const fileName = "textblocks.json"

// Save textBlocks to a file
func saveToFile(textBlocks []TextBlock) error {
	data, err := json.MarshalIndent(textBlocks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

// Load textBlocks from a file
func loadFromFile() ([]TextBlock, error) {
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
