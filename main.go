package main

import (
	"fmt"
	"log"
	"switcheroo/utils"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	config, err := utils.ParseConfig()

	if err != nil {
		log.Fatal(err)
	}

	directories, err := utils.GetAllDirectoryContents(config.Directories)

	if err != nil {
		log.Fatal(err)
	}

	// Initialise fuzzy finder
	selectedIndex, err := fuzzyfinder.Find(
		directories,
		func(i int) string {
			return fmt.Sprintf("%s", directories[i].Name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			return fmt.Sprintf("TODO: metadata")
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	output := utils.DetermineOutput(config, selectedIndex, directories)

	fmt.Printf("%s", output)
}
