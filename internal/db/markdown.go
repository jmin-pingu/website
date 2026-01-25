package db

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/a-h/templ"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Unsafe renders raw HTML as a templ component
func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

// mdToHTML converts markdown to HTML using the custom goldmark parser
func mdToHTML(md []byte) string {
	var buf bytes.Buffer
	customParser := goldmark.New(
		goldmark.WithExtensions(
			NewGenericBlockExtension(),
			extension.Footnote,
			meta.Meta,
			extension.Strikethrough,
			extension.Table,
			mathjax.MathJax,
		),
	)
	if err := customParser.Convert(md, &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	return buf.String()
}

// getMetadata extracts YAML frontmatter metadata from markdown
func getMetadata(md []byte) map[string]interface{} {
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	context := parser.NewContext()
	if err := markdown.Convert(md, &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to extract metadata: %v", err)
	}
	metadata := meta.Get(context)
	return metadata
}

// ----- Generic Block Extension for :::name{props} syntax -----

// GenericBlock represents a custom block with name and properties
type GenericBlock struct {
	ast.BaseBlock
	Name  string
	Props map[string]string
}

func (n *GenericBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindGenericBlock = ast.NewNodeKind("GenericBlock")

func (n *GenericBlock) Kind() ast.NodeKind {
	return KindGenericBlock
}

// genericBlockParser parses :::name{props} blocks
type genericBlockParser struct{}

func NewGenericBlockParser() parser.BlockParser {
	return &genericBlockParser{}
}

func (b *genericBlockParser) Trigger() []byte {
	return []byte{':'}
}

func (b *genericBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()

	// Check for :::
	if len(line) < 3 || line[0] != ':' || line[1] != ':' || line[2] != ':' {
		return nil, parser.NoChildren
	}

	content := strings.TrimSpace(string(line[3:]))

	// Parse :::name{prop1:val1,prop2:val2}
	nameEnd := strings.IndexAny(content, "{ \n")
	if nameEnd == -1 {
		nameEnd = len(content)
	}

	name := content[:nameEnd]
	props := make(map[string]string)

	// Parse properties
	if nameEnd < len(content) && content[nameEnd] == '{' {
		propsEnd := strings.Index(content, "}")
		if propsEnd > nameEnd {
			propsStr := content[nameEnd+1 : propsEnd]
			props = parseProps(propsStr)
		}
	}

	reader.Advance(len(line))

	return &GenericBlock{
		Name:  name,
		Props: props,
	}, parser.HasChildren
}

func (b *genericBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()

	// Check for closing :::
	if len(line) >= 3 && line[0] == ':' && line[1] == ':' && line[2] == ':' {
		reader.Advance(len(line))
		return parser.Close
	}

	return parser.Continue | parser.HasChildren
}

func (b *genericBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (b *genericBlockParser) CanInterruptParagraph() bool                                { return true }
func (b *genericBlockParser) CanAcceptIndentedLine() bool                                { return false }

// parseProps parses "prop1:val1,prop2:val2" into a map
func parseProps(propsStr string) map[string]string {
	props := make(map[string]string)
	pairs := strings.Split(propsStr, ",")

	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			props[key] = val
		}
	}

	return props
}

// GenericBlockHTMLRenderer renders GenericBlock nodes to HTML
type GenericBlockHTMLRenderer struct{}

func NewGenericBlockHTMLRenderer() renderer.NodeRenderer {
	return &GenericBlockHTMLRenderer{}
}

func (r *GenericBlockHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindGenericBlock, r.renderGenericBlock)
}

func (r *GenericBlockHTMLRenderer) renderGenericBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*GenericBlock)

	if entering {
		// Write div with data attributes
		w.WriteString(fmt.Sprintf(`<div data-component="%s"`, n.Name))

		// Add each property as data-{name}-{prop}="{value}"
		for key, val := range n.Props {
			w.WriteString(fmt.Sprintf(` data-%s-%s="%s"`, n.Name, key, val))
		}

		w.WriteString(`>`)
	} else {
		w.WriteString(`</div>`)
	}

	return ast.WalkContinue, nil
}

// GenericBlockExtension is a goldmark extension for :::name{props} blocks
type GenericBlockExtension struct{}

func (e *GenericBlockExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewGenericBlockParser(), 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewGenericBlockHTMLRenderer(), 100),
		),
	)
}

func NewGenericBlockExtension() goldmark.Extender {
	return &GenericBlockExtension{}
}
