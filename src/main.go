package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/jfcarter2358/ceresdb-go/connection"
)

var Version = "1.0.0"

var Suggestions = map[string][]prompt.Suggest{
	"": {
		{Text: "delete", Description: "Delete data"},
		{Text: "get", Description: "Get data"},
		{Text: "patch", Description: "Update part of data"},
		{Text: "post", Description: "Insert new data"},
		{Text: "put", Description: "Update existing data"},
		{Text: "filter", Description: "Filter the results of a GET query"},
		{Text: "orderasc", Description: "Order results of a GET query in ascending order"},
		{Text: "orderdsc", Description: "Order results of a GET query in descending order"},
		{Text: "limit", Description: "Limit the number of results of a GET query"},
		{Text: "count", Description: "Return the number of results of a GET query"},
		{Text: "exit", Description: "Quit the CeresDB client"},
	},
	"|": {
		{Text: "delete", Description: "Delete data"},
		{Text: "get", Description: "Get data"},
		{Text: "patch", Description: "Update part of data"},
		{Text: "post", Description: "Insert new data"},
		{Text: "put", Description: "Update existing data"},
		{Text: "filter", Description: "Filter the results of a GET query"},
		{Text: "orderasc", Description: "Order results of a GET query in ascending order"},
		{Text: "orderdsc", Description: "Order results of a GET query in descending order"},
		{Text: "limit", Description: "Limit the number of results of a GET query"},
		{Text: "count", Description: "Return the number of results of a GET query"},
		{Text: "exit", Description: "Quit the CeresDB client"},
	},
	"delete": {
		{Text: "collection", Description: "Collection resource"},
		{Text: "database", Description: "Database resource"},
		{Text: "permit", Description: "User permission resource"},
		{Text: "record", Description: "Data record resource"},
		{Text: "user", Description: "Instance-wide user resource"},
	},
	"get": {
		{Text: "collection", Description: "Collection resource"},
		{Text: "database", Description: "Database resource"},
		{Text: "permit", Description: "User permission resource"},
		{Text: "record", Description: "Data record resource"},
		{Text: "user", Description: "Instance-wide user resource"},
	},
	"patch": {
		{Text: "record", Description: "Data record resource"},
	},
	"post": {
		{Text: "collection", Description: "Collection resource"},
		{Text: "database", Description: "Database resource"},
		{Text: "permit", Description: "User permission resource"},
		{Text: "record", Description: "Data record resource"},
		{Text: "user", Description: "Instance-wide user resource"},
	},
	"put": {
		{Text: "collection", Description: "Collection resource"},
		{Text: "permit", Description: "User permission resource"},
		{Text: "record", Description: "Data record resource"},
		{Text: "user", Description: "Instance-wide user resource"},
	},
}

func completer(in prompt.Document) []prompt.Suggest {

	text := in.TextBeforeCursor()
	remove := in.GetWordBeforeCursor()
	text = text[:len(text)-len(remove)]
	if len(text) > 0 {
		words := strings.Split(text, " ")
		text = words[len(words)-2]
	}

	s := Suggestions[text]

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	if contains(os.Args[1:], "-h") || contains(os.Args[1:], "--help") {
		printHelp()
		os.Exit(0)
	}
	if len(os.Args) != 5 {
		printHelp()
		os.Exit(1)
	}
	ceresDBUsername := os.Args[1]
	ceresDBPassword := os.Args[2]
	ceresDBHost := os.Args[3]
	ceresDBPort, err := strconv.Atoi(os.Args[4])
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	connection.Initialize(ceresDBUsername, ceresDBPassword, ceresDBHost, ceresDBPort)

	var input string
	history := make([]string, 0)

	fmt.Printf("CeresDB client (%v)\nType \"exit\" to quit\n\n", Version)

	input = prompt.Input(">>> ", completer,
		prompt.OptionTitle("ceresdb-prompt"),
		prompt.OptionHistory(history),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))

	for input != "exit" {
		data, err := connection.Query(input)
		if err != nil {
			printError(err)
		} else {
			if data != nil {
				byteData, _ := json.MarshalIndent(data, "", "    ")
				fmt.Println(string(byteData))
			}
		}

		input = prompt.Input(">>> ", completer,
			prompt.OptionTitle("ceresdb-prompt"),
			prompt.OptionHistory(history),
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))
		history = append(history, input)
	}

}

func printError(err error) {
	fmt.Println(err.Error())
}

func printHelp() {
	help := `usage: %v [option] username password host port
	options:
		-h, --help      Show this help message and exit
	arguments:
		username        Username to use to connect to CeresDB
		password        Password to use to connect to CeresDB
		host            Hostname of CeresDB instance
		port            Port CeresDB instance it running on
`
	fmt.Printf(help, os.Args[0])
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
