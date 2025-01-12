package router

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type MetaData struct {
	ASCII       string
	Title       string
	Description string
}

func ParseHTMLHeader(filename string) (*MetaData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line + "\n")
		if strings.Contains(line, "-->") {
			break
		}
	}

	asciiRe := regexp.MustCompile(`(?s)<ascii>(.*?)</ascii>`)
	titleRe := regexp.MustCompile(`<title>(.*?)</title>`)
	descRe := regexp.MustCompile(`<description>(.*?)</description>`)

	content := sb.String()

	meta := &MetaData{}
	if m := asciiRe.FindStringSubmatch(content); len(m) > 1 {
		meta.ASCII = m[1]
	}
	if m := titleRe.FindStringSubmatch(content); len(m) > 1 {
		meta.Title = m[1]
	}
	if m := descRe.FindStringSubmatch(content); len(m) > 1 {
		meta.Description = m[1]
	}

	return meta, scanner.Err()
}

