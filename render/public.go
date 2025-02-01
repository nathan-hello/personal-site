package render

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Public(input, output string) error {
	atLeastOne := false
	err := filepath.Walk(input, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		route := strings.TrimPrefix(path, "public") // keep "/" in beginning

		dist := output + route

		bits, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		folder := strings.TrimSuffix(dist, info.Name())
		os.MkdirAll(folder, 0777)
		os.WriteFile(dist, bits, 0777)

		atLeastOne = true

		return err
	})

	if !atLeastOne {
		return errors.New("WARN: public/ had no files in it")
	}

	if err != nil {
		return err
	}

	return nil
}
