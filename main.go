package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

func main() {

	var word_amount int = 400

	file, err := os.Open("nno-words.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var nynorsk_words []string

	for len(nynorsk_words) <= word_amount {
		for i := 0; i < len(records); i++ {
			var word_list_on_rank []string = strings.Split(records[i][3], ",")

			nynorsk_words = append(nynorsk_words, word_list_on_rank...)
		}
	}
	if len(nynorsk_words) > word_amount {
		nynorsk_words = nynorsk_words[:word_amount]
	}
	fmt.Println(nynorsk_words)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var dbCSV [][]string

	var word_class string
	var priority string

	askMode := false
	var askModeInput string

	fmt.Println("Vil du spørres om ordene? Skriv 'Ja' dersom du vil.")
	fmt.Scanln(&askModeInput)
	if askModeInput == "Ja" {
		askMode = true
	}
	for i := 0; i < len(nynorsk_words); i++ {
		nynorsk_word := nynorsk_words[i]
		resp, err := http.Get("https://beta.apertium.org/apy/translate?langpair=nno|nob&q=" + nynorsk_word)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		translated_word := strings.Split(string(body), "\"")[5]
		if translated_word != nynorsk_word {

			if askMode {
				fmt.Println("Hva er ordklassen til dette ordet:", translated_word+"/"+nynorsk_word)
				fmt.Scanln(&word_class)
				fmt.Println("Prioriteten?")
				fmt.Scanln(&priority)

				dbCSV = append(dbCSV, []string{nynorsk_word, translated_word, word_class, priority, "translate"})
			} else {
				priority = fmt.Sprint((int(float64(word_amount) / math.Max(float64(len(dbCSV)), 1))))
				dbCSV = append(dbCSV, []string{nynorsk_word, translated_word, "N/A", priority, "translate"})
			}
			fmt.Println(dbCSV)
		}
	}

	outfile, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}

	defer outfile.Close()
	writer = csv.NewWriter(outfile)
	defer writer.Flush()

	_ = writer.Write([]string{"answer", "prompt", "word_class", "priority", "task_type"})
	for _, value := range dbCSV {
		err := writer.Write(value)
		if err != nil {
			fmt.Println("Cannot write to file", err)
		}
	}
}
