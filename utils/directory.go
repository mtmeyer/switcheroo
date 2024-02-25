package utils

import (
	"os"
	"path"
	"sort"
	"strings"
)

type Directory struct {
	Group string
	Path  string
	Name  string
}

func GetSingleDirectoryContents(directory string, group string) ([]Directory, error) {
	// TODO: Allow the directory to be defined from $HOME path, not just absolute
	dirContents, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	directories := []Directory{}

	for _, dir := range dirContents {
		if dir.IsDir() {
			directories = append(directories, Directory{Name: dir.Name(), Group: group, Path: path.Join(directory, dir.Name())})
		}
	}

	return directories, nil
}

func GetAllDirectoryContents(directories map[string]string) ([]Directory, error) {
	allDirectoryContents := []Directory{}

	for group, path := range directories {
		dirContents, err := GetSingleDirectoryContents(path, group)

		if err != nil {
			return nil, err
		}

		allDirectoryContents = append(allDirectoryContents, dirContents...)
	}

	// Sort directories into alphabetical order based on name
	sort.Slice(allDirectoryContents, func(i, j int) bool {
		comparison := strings.Compare(allDirectoryContents[i].Name, allDirectoryContents[j].Name)
		if comparison == -1 {
			return true
		} else {
			return false
		}
	})

	return allDirectoryContents, nil
}
