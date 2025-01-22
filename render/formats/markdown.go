package formats

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/nathan-hello/personal-site/render/customs"
)

func mdCodeHighlighter(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		high, err := customs.CodeHighlighter(string(code.Info), string(code.Literal))
		if err != nil {
			w.Write(code.Literal)
		}
		w.Write([]byte(high))
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func MarkdownRenderer() *mdhtml.Renderer {
	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: mdCodeHighlighter,
	}
	return mdhtml.NewRenderer(opts)
}

func () {
}
