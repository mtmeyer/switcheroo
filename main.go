package main

import (
	"flag"
	"fmt"
	"log"
	"switcheroo/utils"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var directory string
	var configFile string
	var skipPlugins bool

	flag.StringVar(&directory, "directory", "", "Override directory from config")
	flag.StringVar(&configFile, "configFile", "", "Path to custom config file")
	flag.BoolVar(&skipPlugins, "skipPlugins", false, "If plugins should be skipped")

	flag.Parse()

	config, err := utils.ParseConfig(configFile)

	if len(directory) > 0 {
		config.Directories = map[string]string{
			"custom": directory,
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	directories, err := utils.GetAllDirectoryContents(config.Directories)

	if err != nil {
		log.Fatal(err)
	}

	metadata, err := utils.GetMetadataForList(directories, skipPlugins)

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

			if len(metadata) == 0 {
				return ""
			} else {
				return metadata[i]
			}
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	output := utils.DetermineOutput(config, selectedIndex, directories)

	fmt.Println(output)
}
