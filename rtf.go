//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

var headingFonts = []int{
	32, 26, 24, 24, 24,
}

type RtfRenderer struct {
	InText    bool
	ListLevel int
}

func (rtf *RtfRenderer) RenderNode(w io.Writer, node ast.Node,
	entering bool) ast.WalkStatus {

	var nextInText bool

	switch n := node.(type) {
	case *ast.Heading:
		if entering {
			if rtf.InText {
				fmt.Fprintf(w, "\n\\par\n")
			}
			fmt.Fprintf(w, "\n\\par\\fs%d\\b ", headingFonts[n.Level])
		} else {
			fmt.Fprintf(w, "\\b0\n")
		}
		nextInText = true

	case *ast.Paragraph:
		if entering && rtf.InText {
			fmt.Fprintf(w, "\n\\par\\par\n")
		}
		nextInText = true

	case *ast.Text, *ast.Code:
		leaf := node.AsLeaf()
		if leaf != nil {
			var data []byte
			for _, b := range leaf.Literal {
				switch b {
				case '\n':
					data = append(data, ' ')
				case '\\':
					data = append(data, '\\')
					data = append(data, '\\')
				default:
					data = append(data, b)
				}
			}
			w.Write(data)
		}

	case *ast.List:
		if entering {
			rtf.ListLevel++
		} else {
			rtf.ListLevel--
			nextInText = true
		}

	case *ast.ListItem:
		if entering {
			fmt.Fprintf(w, "\\line")
			for i := 1; i < rtf.ListLevel; i++ {
				fmt.Fprintf(w, "\\tab")
			}
			fmt.Fprintf(w, `\~\bullet\~`)
		}

	default:
		fmt.Printf(" - %T %v\n", node, entering)
		nextInText = true
	}

	rtf.InText = nextInText

	return ast.GoToNext
}

func (rtf *RtfRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	fmt.Fprintf(w, `{\rtf1\ansi\ansicpg1252\deff0\deflang1033{\fonttbl{\f0\fswiss\fcharset0 %s;}}\viewkind4\uc1\pard\ql\f0`,
		"Arial")
}

func (rtf *RtfRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	fmt.Fprint(w, "}")
}

func (article *Article) GenerateRTF(dir string) error {
	input := path.Join(article.Dir, "column-article.md")
	output := path.Join(dir, article.RTFOutputName())

	rtf, err := article.ToRTF(input)
	if err != nil {
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(rtf)
	return err
}

func (article *Article) ToRTF(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	parser := parser.NewWithExtensions(article.Extensions)
	doc := markdown.Parse(data, parser)

	return markdown.Render(doc, new(RtfRenderer)), nil
}
