package render

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/a-h/templ"
)

type TemplStaticPages struct {
	Templ templ.Component
	Route string
}

func PagesTempl(output string, templs []TemplStaticPages) error {

	var bits bytes.Buffer
	for _, v := range templs {
		err := v.Templ.Render(context.Background(), &bits)
		if err != nil {
			return err
		}
		parts := strings.Split(v.Route, "/")
		if len(parts) > 1 {
			parts = parts[:len(parts)-1]
		}
		folder := output + strings.Join(parts, "/")
		dist := "dist" + v.Route

		fmt.Printf("INFO: writing file %s in folder %s\n", dist, folder)
		err = os.MkdirAll(folder, 0777)
		if err != nil {
			return err
		}
		err = os.WriteFile(dist, bits.Bytes(), 0777)
		if err != nil {
			return err
		}

	}

	return nil
}
