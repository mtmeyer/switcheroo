package utils

import (
	"os"
	"path"
)

func GetMetadataForList(directories []Directory) ([]string, error) {
	dirs, err := os.ReadDir(path.Join(ConfigDirectory, "plugins"))

	if err != nil {
		return nil, err
	}

	itemMetadata := [][]string{}

	for _, item := range dirs {
		pluginMetadata := GetMetadataFromPlugin(item.Name(), directories)
		itemMetadata = append(itemMetadata, pluginMetadata)
	}

	mergedItemMetadata := mergeMetadata(itemMetadata)

	return mergedItemMetadata, nil
}

func mergeMetadata(metadata [][]string) []string {
	mergedMetadata := []string{}

	for _, pluginMetadata := range metadata {
		for i, item := range pluginMetadata {
			if len(mergedMetadata) == i {
				mergedMetadata = append(mergedMetadata, item+"\n")
			} else {
				mergedMetadata[i] = mergedMetadata[i] + item + "\n"
			}
		}
	}

	return mergedMetadata
}
