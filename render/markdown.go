package render

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func latexRender(b []byte) []byte {
	inFile, err := os.CreateTemp("", "pandoc_input_*.md")
	if err != nil {
		return b
	}
	defer os.Remove(inFile.Name())
	if _, err := inFile.Write(b); err != nil {
		inFile.Close()
		return b
	}
	inFile.Close()

	outFile, err := os.CreateTemp("", "pandoc_output_*.html")
	if err != nil {
		return b
	}
	outFile.Close()
	defer os.Remove(outFile.Name())

	cmd := exec.Command("pandoc", "-s", inFile.Name(), "-o", outFile.Name(), "--mathml")
	if err := cmd.Run(); err != nil {
		return b
	}

	outBytes, err := os.ReadFile(outFile.Name())
	if err != nil {
		return b
	}

   	outStr := string(outBytes)
    start := strings.Index(outStr, "<p>")
    end := strings.Index(outStr, "</body>")

	if start != -1 && end != -1 {
        return []byte(outStr[start:end])
	}
    return []byte(outStr)
}

func hasInlineLatex(b []byte) bool {
    for i := 0; i < len(b); i++ {
        if b[i] == '$' && (i == 0 || b[i-1] != '\\') {
            for j := i + 1; j < len(b); j++ {
                if b[j] == '$' && b[j-1] != '\\' {
                    if j > i+1 {
                        return true
                    }
                    break
                }
            }
        }
    }
    return false
}

func MarkdownRender(md []byte) []byte {
    if hasInlineLatex(md) {
        md = latexRender(md)
    }

    extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	doc := p.Parse(md)

	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: mdRenderHooks,
	}
	renderer := mdhtml.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func mdRenderHooks(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		high, err := CodeHighlighter(string(code.Info), string(code.Literal))
		if err != nil {
			w.Write(code.Literal)
		}
		w.Write([]byte(high))
                
		return ast.GoToNext, true
	}
    if v, ok := node.(*ast.Paragraph); ok {
        if entering && v.Parent != nil {
            w.Write([]byte("<p class=\"pb-2\">"))
        } else {
            w.Write([]byte("</p>"))
        }
        return ast.GoToNext, true
    }

        if v, ok := node.(*ast.HTMLBlock); ok {
                w.Write(v.Literal)
		return ast.GoToNext, true
        }

	return ast.GoToNext, false
}


