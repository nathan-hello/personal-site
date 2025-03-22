package render

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/nathan-hello/personal-site/utils"
)

func Rss(blogs []utils.Blog, output string) error {
    slices.SortFunc(blogs, func(a, b utils.Blog) int {
		return b.Frnt.Date.Compare(a.Frnt.Date)
	})

    rssPath := fmt.Sprintf("%s/rss", output)
    title := "<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\"> <channel> <title> Nat/e </title> <link> https://reluekiss.com </link>"
   	f, err := os.OpenFile(rssPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
    defer f.Close()
    _, err = f.Write([]byte(title))
  	if err != nil {
		return err
	}

    for _, blog := range(blogs)  {
        var image utils.Image
        if len(blog.Frnt.Images) > 0 {
            image = blog.Frnt.Images[0] 
        }
        
        full := strings.Join([]string{
            "<item>",
                "<link>",
                    fmt.Sprintf("https://reluekiss.com/%s", blog.Url),
                "</link>",
                "<title>",
                    blog.Frnt.Title,
                "</title>",
                "<author>",
                    blog.Frnt.Author,
                "</author>",
                "<pubDate>",
                    blog.Frnt.Date.String(),
                "</pubDate>",
                fmt.Sprintf("<enclosure url=\"%s\" length=\"%d\" type=\"image/%s\" />", image.Url, image.BytesCount, strings.TrimPrefix(image.Ext, ".")),
                "<description>",
                    fmt.Sprintf("<![CDATA[alt:%s\n<br/>\n%s...]]>", image.Alt, mdHead(blog.Markdown)),
                "</description>",
            "</item>",
        }, "\n")

        _, err = f.Write([]byte(full))
  	    if err != nil {
	    	return err
	    }
    }
    end := "</channel>\n</rss>"
    _, err = f.Write([]byte(end))
  	if err != nil {
		return err
	}
    return nil
}

func mdHead(md string) string {
	md = strings.TrimLeft(md, "\n")
	if idx := strings.Index(md, "\n"); idx != -1 {
		return md[:idx]
	}
	return md
}
