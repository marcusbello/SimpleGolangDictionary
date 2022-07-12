package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Result struct {
	Definition string `json:"definition"`
	Word       string `json:"word"`
	Valid      bool   `json:"valid"`
}

func main() {

	//Set New fyne instance
	a := app.New()
	win := a.NewWindow("Simple Dictionary")
	win.Resize(fyne.NewSize(550, 250))

	//Header
	title := canvas.NewText("Simple Dictionary", color.White)
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 20

	//Main Section
	bindWord := binding.NewString()
	bindWord.Set("")
	searchBar := widget.NewEntryWithData(bindWord)
	//searchBar.Resize(fyne.NewSize(240, 50))
	searchBar.SetPlaceHolder("Enter word here")

	//searchBtn.Resize(fyne.NewSize(200, 50))

	word := canvas.NewText("Word : ", color.White)
	word.TextStyle = fyne.TextStyle{
		Italic: true,
	}

	wordBox := widget.NewLabel("")
	wordBox.Wrapping = fyne.TextWrapWord

	definition := canvas.NewText("Meaning : ", color.White)
	definition.TextStyle = fyne.TextStyle{
		Italic: true,
	}

	definitionBox := widget.NewLabel("")
	definitionBox.Wrapping = fyne.TextWrapWord

	searchBtn := widget.NewButton("Search", func() {
		wordToSearch, _ := bindWord.Get()
		wordBox.SetText(wordToSearch)
		fmt.Println(wordToSearch)
		getMeaning, err := GetDefinition(wordToSearch)
		if err != nil {
			dialog.ShowError(err, win)
		}
		definitionBox.SetText(getMeaning)
		definitionBox.Resize(fyne.NewSize(100, 100))
		fmt.Println(getMeaning)
		fmt.Println("Button Pressed")
		bindWord.Set("")
	})

	inputContainer := container.New(layout.NewMaxLayout(), searchBar)

	btnContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), searchBtn, layout.NewSpacer())

	wordContainer := container.New(layout.NewGridLayout(2), word, wordBox)

	definitionContainer := container.New(layout.NewGridLayout(2), definition, definitionBox)

	mainContainer := container.New(
		layout.NewVBoxLayout(),
		title,
		layout.NewSpacer(),
		inputContainer,
		btnContainer,
		layout.NewSpacer(),
		wordContainer,
		widget.NewSeparator(),
		definitionContainer,
		layout.NewSpacer(),
	)

	mainContainer.Refresh()
	win.SetContent(mainContainer)
	win.ShowAndRun()
}

func GetDefinition(word string) (string, error) {
	url := "https://api.api-ninjas.com/v1/dictionary?word=" + word
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("first error", err)
		//log.Fatal(err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "KazBImUElDR7LgzPZmppLg==Y6WqA4T2rseizOMz")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("first error")
		//log.Fatal("", err)
		return "", err
	}

	data, _ := ioutil.ReadAll(response.Body)

	var jsonResult Result
	err = json.Unmarshal([]byte(data), &jsonResult)
	cleaner, err := ResultCleaner(jsonResult.Definition)
	if err != nil {
		fmt.Println("first error")
		//log.Fatal("", err)
		return "", err
	}

	return cleaner, nil
}

func ResultCleaner(meaning string) (string, error) {
	regex := regexp.MustCompile("\\d+\\.+\\s")
	if regex.MatchString(meaning) {
		return regex.Split(meaning, -1)[1], nil
	}

	return meaning, nil
}
