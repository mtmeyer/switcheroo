package main

import (
	"fmt"
	"log"
	"switcheroo/utils"
)

func main() {
	config, _ := utils.ParseConfig()

	directories, err := utils.GetAllDirectoryContents(config.Directories)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(directories)
}
