package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Result struct {
	Definition string `json:"definition"`
	Word       string `json:"word"`
	Valid      bool   `json:"valid"`
}

func main() {
	word := os.Args[1]
	fmt.Println("Dictionary App")
	GetDefinition(word)
}

func GetDefinition(word string) {
	url := "https://api.api-ninjas.com/v1/dictionary?word=" + word
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "KazBImUElDR7LgzPZmppLg==Y6WqA4T2rseizOMz")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("", err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(data))
	var jsonResult Result
	err = json.Unmarshal([]byte(data), &jsonResult)
	//fmt.Println(jsonResult.Definition)
	fmt.Println(ResultCleaner(jsonResult.Definition)[1])

}

func ResultCleaner(meaning string) []string {
	regex := regexp.MustCompile("\\d+\\.+\\s")
	return regex.Split(meaning, -1)
}
