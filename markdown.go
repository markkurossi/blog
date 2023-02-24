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
		opts.RenderNodeHook = article.renderPresentation
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

func className(pagenum int) string {
	switch pagenum {
	case 0:
		return "current"
	case 1:
		return "next"
	default:
		return "far-next"
	}
}

// I returns an indent string for the specified level.
func I(count int) string {
	result := "      "
	for i := 0; i < count; i++ {
		result += "  "
	}
	return result
}

func (article *Article) renderPresentation(w io.Writer, node ast.Node,
	entering bool) (ast.WalkStatus, bool) {

	switch n := node.(type) {
	case *ast.Heading:
		if n.Level < 3 {
			// Slide.

			var tag string
			if n.Level == 1 {
				tag = "h1"
			} else {
				tag = "h3"
			}

			if entering {
				if article.Pagenum > 0 {
					if article.Pagenum == 1 {
						fmt.Fprintf(w, "%s</div>\n", I(1))
					} else {
						fmt.Fprintf(w,
							"%s<span class=\"pagenumber\">%d</span>\n",
							I(1), article.Pagenum)
					}
					fmt.Fprintf(w, "%s</article>\n\n", I(0))
				}

				fmt.Fprintf(w, "%s<article class=\"%s\">\n",
					I(0), className(article.Pagenum))
				fmt.Fprintf(w, "%s<%s>", I(1), tag)
			} else {
				fmt.Fprintf(w, "</%s>\n", tag)
				if article.Pagenum == 0 {
					fmt.Fprintf(w, "%s<div class=\"presenter\">\n", I(1))
				}
				article.Pagenum++
			}
		} else {
			// Slide column.
		}
		return ast.GoToNext, true

	default:
		return ast.GoToNext, false
	}
}
