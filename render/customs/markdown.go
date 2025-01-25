package customs

import (
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MarkdownRender(md []byte) []byte {

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: mdCodeHighlighter,
	}
	renderer := mdhtml.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func mdCodeHighlighter(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		high, err := CodeHighlighter(string(code.Info), string(code.Literal))
                if err != nil {
			w.Write(code.Literal)
		}
		w.Write([]byte(high))
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
