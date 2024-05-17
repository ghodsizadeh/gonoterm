package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/rivo/tview"
)

var mu sync.Mutex

// Update the text blocks in the grid
func updateTextBlocks(grid *tview.Grid, textBlocks []TextBlock, inputFields *[]*tview.TextArea) {
	grid.Clear()
	*inputFields = nil
	for i := range textBlocks {
		if textBlocks[i].Type == Code {
			inputField := tview.NewTextArea().
				SetText(textBlocks[i].Text, true).
				SetPlaceholder("Enter code here...")
			inputField.SetBorder(true)
			inputField.SetTitle(fmt.Sprintf("[green]Code Block %d", i+1))

			resultView := tview.NewTextView().
				SetDynamicColors(true).
				SetWrap(true)
			resultView.SetBorder(true)
			resultView.SetTitle("Results")

			grid.AddItem(inputField, i, 0, 1, 7, 0, 0, true)
			grid.AddItem(resultView, i, 7, 1, 3, 0, 0, false)

			index := i

			// Fill the result view with the first result
			first_result, _ := interpretCode(textBlocks[i].Text)
			log.Println("First result: ", first_result)
			resultView.SetText(strings.Join(first_result, ""))

			// Update the result view with the new results
			inputField.SetChangedFunc(func() {
				mu.Lock()
				defer mu.Unlock()

				text := inputField.GetText()
				textBlocks[index].Text = text
				go func(text string, resultView *tview.TextView) {
					result, err := interpretCode(text)
					if err != nil {
						log.Println(err)
						resultView.SetText(err.Error())
					} else {
						mu.Lock()
						defer mu.Unlock()
						resultView.SetText(strings.Join(result, ""))
					}
				}(text, resultView)
			})

			*inputFields = append(*inputFields, inputField)
		} else {
			inputField := tview.NewTextArea().
				SetText(textBlocks[i].Text, true).
				SetPlaceholder("Enter text here...")
			inputField.SetBorder(true)
			inputField.SetTitle(fmt.Sprintf("[green]Block %d %s", i+1, textBlocks[i].Type.String()))

			grid.AddItem(inputField, i, 0, 1, 10, 0, 0, true)

			index := i

			inputField.SetChangedFunc(func() {
				mu.Lock()
				defer mu.Unlock()

				text := inputField.GetText()
				textBlocks[index].Text = text
			})

			*inputFields = append(*inputFields, inputField)
		}
	}
}

// Update lines with results in a text area
func updateLinesWithResults(oldLines []string, results []string, ta *tview.TextArea) {
	newLines := []string{}
	for i, line := range oldLines {
		newLines = append(newLines, line)
		if i < len(results) {
			newLines = append(newLines, results[i])
		}
	}
	ta.SetText(strings.Join(newLines, "\n"), true)
}

// Interpret code to execute and return results
func interpretCode(input string) ([]string, error) {
	lines := strings.Split(input, "\n")
	results := []string{}

	for _, line := range lines {
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			results = append(results, "\n")
			continue
		}
		code := strings.Join(tokens[:], " ")
		result, _ := executePythonCode(code)

		results = append(results, result)
	}

	return results, nil
}

// executeGoCode function to run Go code and return the output
func executeGoCode(code string) (string, error) {
	cmd := exec.Command("go", "run", "-")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("package main\nimport \"fmt\"\nfunc main() { %s }", code))
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if errOut.Len() > 0 {
		return "", fmt.Errorf(errOut.String())
	}

	return out.String(), nil
}

// executePythonCode function to run Python code and return the output
func executePythonCode(code string) (string, error) {
	// cmd := exec.Command("python3", "-c", code)
	// print(code)
	cmd := exec.Command("python3", "-c", "from math import *;print("+code+")")
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	log.Println("Python", err, out.String(), errOut.String(), code)
	if err != nil {
		return " ", err
	}

	if errOut.Len() > 0 {
		return "", fmt.Errorf(errOut.String())
	}

	return out.String(), nil
}
