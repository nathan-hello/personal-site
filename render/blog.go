package render

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Frontmatter struct {
	Author         string
	Title          string
	Date           time.Time
	Images         []Image
	tags           []string
	overrideHref   *string
	overrideLayout *string
	description    *string
	hidden         *bool
}

func frontmatter(f *os.File) (map[string]string, error) {
	props := map[string]string{}
	scanner := bufio.NewScanner(f)

	begin := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "---") {
			if begin {
				break
			}
			begin = true
			continue
		}
		kv := strings.Split(line, ":")

		key := kv[0]
		if key {
		}

	}
	if !begin {
		return nil, fmt.Errorf("file %s reached EOL without frontmatter", f.Name())
	}

	return props, nil
}

type imageFrontmatterProps struct {
	Alt string `json:"alt"`
}

type imageFrontmatter map[string]imageFrontmatterProps
type Image struct {
	Name     string
	Filesize string
	Url      string
	AltText  string
}

func parseImagesProp(s []byte, date string) []Image {

	ifm := imageFrontmatter{}
	err := json.Unmarshal(s, &ifm)
	if err != nil {
		panic(err)
	}

	images := []Image{}

	for k, v := range ifm {
		asdf := os.Open()
		images = append(images, Image{
			Name:    k,
			Url:     k,
			AltText: v.Alt,
		})
	}

}
