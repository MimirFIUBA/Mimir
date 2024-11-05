package utils

import (
	"encoding/json"
	"io"
	"io/fs"
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
		return nil
	}

	var files []string
	for _, v := range mdFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}

func DecodeJsonToMap(reader io.Reader, result any) error {
	decoder := json.NewDecoder(reader)
	for {
		err := decoder.Decode(result)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
