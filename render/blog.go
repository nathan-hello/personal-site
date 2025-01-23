package render

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nathan-hello/personal-site/utils"
)

type Frontmatter struct {
	Author         string
	Title          string
	Date           time.Time
	Images         []Image
	Tags           []string
	OverrideHref   string
	OverrideLayout string
	Description    string
	Hidden         bool
}

func Blogs() {
	id := 100001
	filepath.Walk(utils.DIR_BLOG, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if filepath.Ext(info.Name()) == ".html" {
			err := writeHtmlFile(f)

		}

		return nil
	})
}

func frontmatter(f *os.File) (Frontmatter, error) {
	scanner := bufio.NewScanner(f)
	begin := false
	frontmatter := Frontmatter{}
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
		value := strings.Join(kv[1:], ":")

		err := parseKey(key, value, &frontmatter)
		if err != nil {
			return Frontmatter{}, err
		}
		if key == "images" || key == "image" {
			parseImageFrontmatter(value, &frontmatter)
		}
	}
	if !begin {
		return Frontmatter{}, fmt.Errorf("file %s reached EOL without frontmatter", f.Name())
	}

	return frontmatter, nil
}

func parseKey(key string, value string, f *Frontmatter) error {

	switch strings.ToLower(key) {
	case "title":
		f.Title = value
	case "author":
		f.Author = value
	case "date":
		f.Date = utils.DateStringToObject(value)
	case "tags":
		arr := []string{}
		err := json.Unmarshal([]byte(value), &arr)
		if err != nil {
			return err
		}
		f.Tags = arr
	case "overrideHref":
		f.OverrideHref = value
	case "overrideLayout":
		f.OverrideLayout = value
	case "description":
		f.Description = value
	case "hidden":
		if value == "true" {
			f.Hidden = true
		}
	}

	return nil

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

func parseImageFrontmatter(s string, frontmatter *Frontmatter) error {

	ifm := imageFrontmatter{}
	err := json.Unmarshal([]byte(s), &ifm)
	if err != nil {
		panic(err)
	}

	images := []Image{}

	folder := utils.GetBlogImageDir(frontmatter.Date)
	path := fmt.Sprintf("%s/%s", folder, s)

	for k, v := range ifm {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		size := utils.FormatSize(stat.Size())
		images = append(images, Image{
			Name:     k,
			Url:      path,
			AltText:  v.Alt,
			Filesize: size,
		})
	}
	frontmatter.Images = images
	return nil

}
