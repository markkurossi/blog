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
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/markkurossi/blog/asciiart"
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

	case *ast.CodeBlock:
		data := string(n.Literal)
		class := "code"
		var err error

		for _, f := range strings.Split(string(n.Info), ",") {
			Verbose(" - filter: %v\n", f)
			data, class, err = filter(f)(data, class)
			if err != nil {
				fmt.Printf("filter %s: %s\n", f, err)
			}
		}

		fmt.Fprintf(w, "<pre class=\"%s\">\n", class)
		io.WriteString(w, html.EscapeString(data))
		io.WriteString(w, "</pre>\n")
		return ast.GoToNext, true

	case *ast.Image:
		if entering {
			fmt.Fprintf(w, `<p align="center"><img src="%s" title="%s"/>`,
				string(n.Destination), html.EscapeString(string(n.Title)))
		} else {
			fmt.Fprintf(w, "</p>\n")
		}
		return ast.SkipChildren, true

	default:
		return ast.GoToNext, false
	}
}

func filter(name string) func(data, class string) (string, string, error) {
	switch name {
	case "ascii-art":
		return filterASCIIArt

	case "linenumbers":
		return filterLinenumbers

	case "plain":
		return filterPlain

	case "center":
		return filterCenter

	default:
		return filterPassthrough
	}
}

func filterASCIIArt(data, class string) (string, string, error) {
	return asciiart.Process(data), "ascii-art", nil
}

func filterLinenumbers(data, class string) (string, string, error) {
	var result string
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var format string
	if len(lines) > 9 {
		format = "%2d %s\n"
	} else {
		format = "%d %s\n"
	}
	for idx, line := range lines {
		result += fmt.Sprintf(format, idx+1, line)
	}
	return result, class, nil
}

func filterCenter(data, class string) (string, string, error) {
	return data, class + " center", nil
}

func filterPassthrough(data, class string) (string, string, error) {
	return data, class, nil
}

func filterPlain(data, class string) (string, string, error) {
	return data, "code-plain", nil
}
