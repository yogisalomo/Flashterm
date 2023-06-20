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

type Card struct {
	Key string
	Value string
	Weight float32
}

func commandInputValidation(input string) (appState string) {
	if input == "/menu" {
		return "menu"
	} else if input == "/quit" {
		return "quit"
	}

	return ""
}

func WeightedRandom(vocabulary []Card) int {
    n := len(vocabulary)
    if n == 0 {
        return 0
    }
    cdf := make([]float32, n)
    var sum float32 = 0.0
    for idx, card := range vocabulary {
        if idx > 0 {
            cdf[idx] = cdf[idx-1] + card.Weight
        } else {
            cdf[idx] = card.Weight
        }
        sum += card.Weight
    }
    r := rand.Float32() * sum
    var l, h int = 0, n - 1
    for l <= h {
        m := l + (h-l)/2
        if r <= cdf[m] {
            if m == 0 || (m > 0 && r > cdf[m-1]) {
                return m
            }
            h = m - 1
        } else {
            l = m + 1
        }
    }
    return -1
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
	var vocabulary []Card

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

		vocabulary = append(vocabulary, Card{Key: dataLine[0], Value: dataLine[1], Weight: 100.0})
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
			vocabulary = append(vocabulary, Card{Key: key, Value: value})

			_, err := io.WriteString(dataWriter, fmt.Sprintf("%s,%s\n", key, value))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			dataWriter.Flush()

			pterm.Success.Printfln(fmt.Sprintf("Stored %s : %s", key, value))
			pterm.Info.Printfln("Input /menu to go to menu, /quit to quit the program.\n\n")
		} else if appState == "test" {
			if len(vocabulary) == 0 {
				pterm.Warning.Printfln("No record on the data, cannot test yet. Please record some data first")
				appState = "record"
				continue
			}
			
			testedIndex := WeightedRandom(vocabulary)
			testedKey := vocabulary[testedIndex].Key
			testedValue := vocabulary[testedIndex].Value

			guessedValue, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(testedKey).Show()
			updatedState := commandInputValidation(guessedValue)
			if len(updatedState) > 0 {
				appState = updatedState
				continue
			}

			if guessedValue == testedValue {
				pterm.Success.Printfln(fmt.Sprintf("You guessed correctly! (%s : %s)", testedKey, guessedValue))
				currentStreak+=1
				if vocabulary[testedIndex].Weight > 10 {
					vocabulary[testedIndex].Weight -= 10
				}
				pterm.Success.Printfln(fmt.Sprintf("Current streak: %d", currentStreak))
			} else {
				pterm.Error.Printfln(fmt.Sprintf("You guessed wrong! :( (%s : %s) should be %s", testedKey, guessedValue, testedValue))
				currentStreak = 0
				if vocabulary[testedIndex].Weight < 100 {
					vocabulary[testedIndex].Weight += 10
				}
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