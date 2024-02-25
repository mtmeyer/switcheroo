package main

import (
	"fmt"
	"switcheroo/src"
)

func main() {
	config, _ := src.ParseConfig()
	fmt.Println(config.Directories)
}
