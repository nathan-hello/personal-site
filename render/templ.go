package render

import (
	"bytes"
	"context"
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
		v.Templ.Render(context.Background(), &bits)
		parts := strings.Split(v.Route, "/")
		if len(parts) > 1 {
			parts = parts[:len(parts)-1]
		}
		folder := output + strings.Join(parts, "/")
		os.MkdirAll(folder, 0777)
		os.WriteFile("dist"+v.Route, bits.Bytes(), 0777)
	}

	return nil
}
