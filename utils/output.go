package utils

import (
	"strconv"
)

func DetermineOutput(config *Config, selectedIndex int, directories []Directory) string {
	switch config.Output {
	case "path":
		return directories[selectedIndex].Path
	case "index":
		return strconv.FormatInt(int64(selectedIndex), 10)
	case "name":
		return directories[selectedIndex].Name
	default:
		return directories[selectedIndex].Name
	}
}
