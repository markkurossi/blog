//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (article *Article) format(data []byte) []byte {
	parser := parser.NewWithExtensions(article.Extensions)

	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderCodeBlock,
	}
	renderer := html.NewRenderer(opts)

	return markdown.ToHTML(data, parser, renderer)
}

func renderCodeBlock(w io.Writer, node ast.Node, entering bool) (
	ast.WalkStatus, bool) {
	code, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.GoToNext, false
	}
	io.WriteString(w, "<pre>\n")
	if len(code.Info) > 00 {
		io.WriteString(w, fmt.Sprintf("%s:\n%s",
			code.Info, code.Literal))
	} else {
		w.Write(code.Literal)
	}
	io.WriteString(w, "</pre>\n")
	return ast.GoToNext, true
}
