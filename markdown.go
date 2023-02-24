//
// Copyright (c) 2021-2023 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"html"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (article *Article) format(data []byte) []byte {
	parser := parser.NewWithExtensions(article.Extensions)

	opts := mdhtml.RendererOptions{
		Flags: mdhtml.CommonFlags,
	}
	switch article.Type() {
	case TmplPresentation:
		opts.RenderNodeHook = renderPresentation
	case TmplArticle:
		opts.RenderNodeHook = renderArticle
	}

	renderer := mdhtml.NewRenderer(opts)

	return markdown.ToHTML(data, parser, renderer)
}

func renderArticle(w io.Writer, node ast.Node, entering bool) (
	ast.WalkStatus, bool) {
	code, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.GoToNext, false
	}
	io.WriteString(w, "<pre>\n")
	if len(code.Info) > 0 {
		io.WriteString(w, fmt.Sprintf("%s:\n%s",
			code.Info, html.EscapeString(string(code.Literal))))
	} else {
		io.WriteString(w, html.EscapeString(string(code.Literal)))
	}
	io.WriteString(w, "</pre>\n")
	return ast.GoToNext, true
}

func renderPresentation(w io.Writer, node ast.Node, entering bool) (
	ast.WalkStatus, bool) {
	code, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.GoToNext, false
	}
	io.WriteString(w, "<pre>\n")
	if len(code.Info) > 0 {
		io.WriteString(w, fmt.Sprintf("%s:\n%s",
			code.Info, html.EscapeString(string(code.Literal))))
	} else {
		io.WriteString(w, html.EscapeString(string(code.Literal)))
	}
	io.WriteString(w, "</pre>\n")
	return ast.GoToNext, true
}
