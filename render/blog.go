package render

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/layouts"
	"github.com/nathan-hello/personal-site/utils"
	"gopkg.in/yaml.v3"
)

func Blogs(input, output string, write bool) ([]utils.Blog, error) {

	blogs, err := gatherRenderedHtmls(input)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(blogs, func(a, b utils.Blog) int {
		return a.Frnt.Date.Compare(b.Frnt.Date)
	})

	for i, v := range blogs {
		v.Id = 100_000 + i + 1

		dist := fmt.Sprintf("%s/%s/p/%d.html", output, strings.ToLower(v.Frnt.Author), v.Id)
		v.Url = strings.TrimPrefix(dist, "./dist")

		v.Comments = getComments(int64(v.Id))

		blogs[i] = v

		if write {
			err = writeBlogPost(v, dist)
			if err != nil {
				return nil, err
			}
		}

	}

	return blogs, nil
}

func writeBlogPost(v utils.Blog, dist string) error {
	var bits bytes.Buffer

	comp := chooseBlogLayout(v)
	err := comp.Render(context.Background(), &bits)
	if err != nil {
		return err
	}

	parts := strings.Split(dist, "/")
	folder := strings.Join(parts[:len(parts)-1], "/")

	fmt.Printf("INFO: writing file %s in folder %s\n", dist, folder)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		return err
	}
	err = os.WriteFile(dist, bits.Bytes(), 0777)
	if err != nil {
		return err
	}
	return nil
}

func gatherRenderedHtmls(input string) ([]utils.Blog, error) {
	blogs := []utils.Blog{}
	err := filepath.Walk(input, func(path string, info fs.FileInfo, err error) error {
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

		b := utils.Blog{
			Frnt: utils.Frontmatter{},
		}

		yml, content, err := yoinkFrontmatter(f)
		if err != nil {
			return err
		}

		b.Frnt, err = parseFrontmatter(yml)
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".html" {
			rendered, err := RenderCustomComponents(content)

			if err != nil {
				return err
			}
			b.Html = rendered
		}

		if filepath.Ext(info.Name()) == ".md" || filepath.Ext(info.Name()) == ".mdx" {
			rendered := MarkdownRender([]byte(content))
			b.Html = string(rendered)
		}

		blogs = append(blogs, b)

		return nil
	})
	return blogs, err
}

func yoinkFrontmatter(f *os.File) (string, string, error) {

	ext := filepath.Ext(f.Name())
	asdf, err := io.ReadAll(f)
	if err != nil {
		return "", "", err
	}
	text := string(asdf)
	lines := strings.Split(text, "\n")

	var delims [2]string
	if ext == ".mdx" || ext == ".md" {
		delims = [2]string{"---", "---"}
	}
	if ext == ".html" {
		delims = [2]string{"<!--", "-->"}
	}
	var idxs = [2]int{-1, -1}

	var a = 0
	for i, v := range lines {
		if strings.Contains(v, delims[a]) {
			idxs[a] = i
			a++
		}
		if idxs[0] > -1 && idxs[1] > -1 {
			break
		}
	}

	if idxs[0] == -1 || idxs[1] == -1 {
		return "", "", fmt.Errorf("could not get frontmatter for file %s", f.Name())
	}
	frntString := strings.Join(lines[idxs[0]+1:idxs[1]], "\n")
	content := strings.Join(lines[idxs[1]+1:], "\n")

	return frntString, content, nil
}

type ymlImage struct {
	Alt string `yaml:"alt,omitempty"`
}
type ymlFrontmatter struct {
	Author         string              `yaml:"author"`
	Title          string              `yaml:"title"`
	Date           string              `yaml:"date"`
	Images         map[string]ymlImage `yaml:"images,omitempty"`
	Image          map[string]ymlImage `yaml:"image,omitempty"`
	Tags           []string            `yaml:"tags,omitempty"`
	OverrideHref   string              `yaml:"overrideHref,omitempty"`
	OverrideLayout string              `yaml:"overrideLayout,omitempty"`
	Description    string              `yaml:"description,omitempty"`
	Hidden         bool                `yaml:"hidden,omitempty"`
}

func parseFrontmatter(s string) (utils.Frontmatter, error) {

	yml := ymlFrontmatter{}

	err := yaml.Unmarshal([]byte(s), &yml)
	if err != nil {
		log.Println(s)
		return utils.Frontmatter{}, err
	}

	fm := utils.Frontmatter{}
	fm.Author = yml.Author
	fm.Title = yml.Title
	fm.Date, err = utils.DateStringToObject(yml.Date)
	if err != nil {
		return utils.Frontmatter{}, err
	}

	getImages(yml.Image, &fm)  // key "image:"
	getImages(yml.Images, &fm) // key "images:"

	fm.Tags = yml.Tags
	fm.OverrideHref = yml.OverrideHref
	fm.OverrideLayout = yml.OverrideLayout
	fm.Description = yml.Description
	fm.Hidden = yml.Hidden

	return fm, nil

}

func blogImageLocation(name string, d time.Time) string {
	year := d.Year()
	return fmt.Sprintf("./public/images/covers/%d/%s", year, name)
}

func getImages(yml map[string]ymlImage, fm *utils.Frontmatter) error {
	for k, v := range yml {
		publicDir := blogImageLocation(k, fm.Date)
		url := strings.TrimPrefix(publicDir, "./public")
		f, err := os.Open(publicDir)
		if err != nil {
			return err
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			return err
		}

		fm.Images = append(fm.Images,
			utils.Image{
				Name: k,
				Size: utils.FormatSize(stat.Size()),
				Ext:  filepath.Ext(k),
				Url:  url,
				Alt:  v.Alt,
			},
		)

	}

	return nil
}

type BlogComponent = func(utils.Blog) templ.Component

var registeredBlogLayouts = map[string]BlogComponent{
	"natalie": layouts.NataliePost,
	"nathan":  layouts.NathanPost,
}

func chooseBlogLayout(blog utils.Blog) templ.Component {
	comp, ok := registeredBlogLayouts[blog.Frnt.Author]
	if !ok {
		return registeredBlogLayouts["nathan"](blog)
	}
	return comp(blog)
}

func getComments(blogId int64) []utils.Comment {
	cmts, err := db.Conn.SelectCommentsMany(context.Background(), blogId)
	if err != nil {
		return nil
	}
	c := []utils.Comment{}
	for _, v := range cmts {
		c = append(c, v.NewBlogComment())
	}

	return c
}
