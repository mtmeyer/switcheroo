package utils

import (
	"os"
)

func GetDirectoryContents(directory string) ([]string, error) {
	// TODO: Allow the directory to be defined from $HOME path, not just absolute
	dirContents, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	directories := []string{}

	for _, dir := range dirContents {
		if dir.IsDir() {
			directories = append(directories, dir.Name())
		}
	}

	return directories, nil
}

func GetAllDirectoryContents(directories map[string]string) (map[string][]string, error) {
	allDirectoryContents := map[string][]string{}

	for name, path := range directories {
		dirContents, err := GetDirectoryContents(path)

		if err != nil {
			return nil, err
		}

		allDirectoryContents[name] = dirContents
	}

	return allDirectoryContents, nil
}
