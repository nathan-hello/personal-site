package render

import (
	"errors"
	"fmt"
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

		fmt.Printf("INFO: writing file %s in folder %s\n", dist, folder)
		err = os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
		err = os.WriteFile(dist, bits, 0777)
		if err != nil {
			return err
		}

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
