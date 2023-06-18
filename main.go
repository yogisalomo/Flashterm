package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func commandInputValidation(input string) (appState string) {
	if input == "/menu" {
		return "menu"
	} else if input == "/quit" {
		return "quit"
	}

	return ""
}

func randomPick(vocabulary map[string]string, keys []string) (key string, value string) {
	keyIndex := rand.Int() % len(keys)
	return keys[keyIndex], vocabulary[keys[keyIndex]]
}

func main() {
	var appState string
	dataFilePath := "./data.csv"

	// Check if data file exist & if not, create an empty one
	if _, err := os.Stat(dataFilePath); errors.Is(err, os.ErrNotExist) {
		newDataFile, err := os.Create(dataFilePath)
		if err != nil {
			fmt.Println("Failed to create file to store data")
			os.Exit(1)
		}

		defer newDataFile.Close()
	}

	// Initialize in-app vocabulary struct
	vocabulary := make(map[string]string)
	var keys []string

	// Open the existing data file & populate vocabulary with current data
	dataFile, err := os.OpenFile(dataFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println(err)
		os.Exit(1)
    }
	defer dataFile.Close()

    fileScanner := bufio.NewScanner(dataFile)
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
        dataLine := strings.Split(fileScanner.Text(), ",")
		if len(dataLine) != 2 {
			fmt.Println("corrupted data on following line:", fileScanner.Text())
			os.Exit(1)
		}

		vocabulary[dataLine[0]] = dataLine[1]
		keys = append(keys, dataLine[0])
    }
	
	// Open Menu & Query user for input
	pterm.DefaultHeader.WithFullWidth(true).Println("Flashterm")
	appState, _ = pterm.DefaultInteractiveSelect.
	WithOptions([]string{"record", "test", "quit"}).
	Show()

	dataWriter := bufio.NewWriter(dataFile)

	currentStreak := 0
	
	for {
		if appState == "record" {
			key, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Key").Show()
			updatedState := commandInputValidation(key)
			if len(updatedState) > 0 {
				appState = updatedState
				continue
			}
			
			value, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Value").Show()
			updatedState = commandInputValidation(key)
			if len(updatedState) > 0 {
				appState = updatedState
				continue
			}
			
			// Update the vocabulary, store to persistent file, and print info to the user
			vocabulary[key] = value
			keys = append(keys, key)

			writtenLen, err := io.WriteString(dataWriter, fmt.Sprintf("%s,%s\n", key, value))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("written this lin:", writtenLen)
			dataWriter.Flush()

			pterm.Success.Printfln(fmt.Sprintf("Stored %s : %s", key, value))
			pterm.Info.Printfln("Input /menu to go to menu, /quit to quit the program.\n\n")
		} else if appState == "test" {
			if len(vocabulary) == 0 || len(keys) == 0 {
				pterm.Warning.Printfln("No record on the data, cannot test yet. Please record some data first")
				appState = "record"
				continue
			}
			
			testedKey, testedValue := randomPick(vocabulary, keys)

			guessedValue, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(testedKey).Show()
			updatedState := commandInputValidation(guessedValue)
			if len(updatedState) > 0 {
				appState = updatedState
				continue
			}

			if guessedValue == testedValue {
				pterm.Success.Printfln(fmt.Sprintf("You guessed correctly! (%s : %s)", testedKey, guessedValue))
				currentStreak+=1
				pterm.Success.Printfln(fmt.Sprintf("Current streak: %d", currentStreak))
			} else {
				pterm.Error.Printfln(fmt.Sprintf("You guessed wrong! :( (%s : %s) should be %s", testedKey, guessedValue, testedValue))
				currentStreak = 0
			}
		} else if appState == "menu" {
			appState, _ = pterm.DefaultInteractiveSelect.
			WithOptions([]string{"record", "test", "quit"}).
			Show()
		} else {
			// Quit state or invalid state
			break
		}
	}

	pterm.Info.Printfln("see you around!")
}