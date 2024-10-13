package utils

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

func GetValueFromJSON(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	var value interface{} = data

	for _, key := range keys {
		// Verificamos si el valor actual es un mapa
		if mapVal, ok := value.(map[string]interface{}); ok {
			value = mapVal[key]
		} else {
			return nil, false
		}
	}

	return value, true
}

func ListFilesWithSuffix(dir, suffix string) []string {
	root := os.DirFS(dir)

	mdFiles, err := fs.Glob(root, suffix)

	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, v := range mdFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}
