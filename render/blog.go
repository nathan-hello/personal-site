package render

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/nathan-hello/personal-site/render/customs"
	"github.com/nathan-hello/personal-site/utils"
	"gopkg.in/yaml.v3"
)

type image struct {
	alt string `yaml:"alt"`
}

type frontmatter struct {
	Author         string   `yaml:"author"`
	Title          string   `yaml:"title"`
	Date           string   `yaml:"date"`
	Images         []image  `yaml:"images,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	OverrideHref   string   `yaml:"overrideHref,omitempty"`
	OverrideLayout string   `yaml:"overrideLayout,omitempty"`
	Description    string   `yaml:"description,omitempty"`
	Hidden         bool     `yaml:"hidden,omitempty"`
}

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

type Blog struct {
	Id   int
	Frnt Frontmatter
	Html string
}

func Blogs() error {

	blogs, err := gatherRenderedHtmls()
	if err != nil {
		return err
	}

	slices.SortFunc(blogs, func(a, b Blog) int {
		return a.Frnt.Date.Compare(b.Frnt.Date)
	})

	for _, _ = range blogs {
		//fmt.Printf("%d: %#v\n", i, v)

	}

	return nil
}

func gatherRenderedHtmls() ([]Blog, error) {
	blogs := []Blog{}
	err := filepath.Walk(utils.DIR_BLOG, func(path string, info fs.FileInfo, err error) error {
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

		fm, err := parseBlogFrontmatter(f)
		if err != nil {
			return err
		}
		b := Blog{
			Frnt: fm,
		}

		parseYaml(f)

		if filepath.Ext(info.Name()) == ".html" {

			rendered, err := customs.RenderCustomComponents(f)

			if err != nil {
				return err
			}
			b.Html = rendered
		}

		if filepath.Ext(info.Name()) == ".md" || filepath.Ext(info.Name()) == ".mdx" {

			contents, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			rendered := customs.MarkdownRender(contents)
			b.Html = string(rendered)

		}

		blogs = append(blogs, b)

		return nil
	})
	return blogs, err
}

func parseYaml(f *os.File) (string, error) {
	// Read entire file at once
	raw, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(f.Name())
	var delims [2]string
	switch ext {
	case ".mdx":
		delims = [2]string{"---", "---"}
	case ".html":
		delims = [2]string{"<!--", "-->"}
	}

	// Build a safe regex that matches anything (including newlines) between delimiters
	re := fmt.Sprintf(`(?s)%s(.*?)%s`, regexp.QuoteMeta(delims[0]), regexp.QuoteMeta(delims[1]))
	rx := regexp.MustCompile(re)

	// We want just the content inside the delimiters
	match := rx.FindStringSubmatch(string(raw))
	if len(match) < 2 {
		return "", nil
	}
	fmt.Printf("%#v\n", match)

	return match[1], nil

}

func parseBlogFrontmatter(f *os.File) (Frontmatter, error) {
	scanner := bufio.NewScanner(f)
	var begin, end bool
	frontmatter := Frontmatter{}
	for scanner.Scan() {

		line := scanner.Text()
		begin, end = testBeginningEnd(line, filepath.Ext(f.Name()), begin)

		if end {
			break
		}

		kv := strings.Split(line, ":")

		key := kv[0]
		value := strings.Join(kv[1:], ":")

		err := parseKey(key, value, &frontmatter)
		if err != nil {
			return Frontmatter{}, err
		}
		if key == "images" || key == "image" {
			parseImageFrontmatter([]byte(value), &frontmatter)
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

func parseImageFrontmatter(yamlData []byte, frontmatter *Frontmatter) error {
	var fm struct {
		Images imageFrontmatter `yaml:"images"`
	}

	err := yaml.Unmarshal(yamlData, &fm)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	ifm := fm.Images

	images := []Image{}

	folder := utils.GetBlogImageDir(frontmatter.Date)

	for k, v := range ifm {
		path := fmt.Sprintf("%s/%s", folder, k) // Use 'k' instead of 's'

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", path, err)
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

func testBeginningEnd(line, ext string, begin bool) (bool, bool) {
	if ext == ".md" || ext == ".mdx" {
		if strings.TrimSpace(line) == "---" {
			if !begin {
				return true, false
			}
			return true, true
		}
	}

	if ext == ".html" {
		if strings.Contains(line, "<!--") {
			return true, false
		}
		if strings.Contains(line, "-->") {
			return true, true
		}
	}
	return begin, false
}
