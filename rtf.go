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

// RtfRenderer implements Markdown RTF renderer.
type RtfRenderer struct {
	InText    bool
	ListLevel int
}

// RenderNode renders Markdown node to RTF.
func (rtf *RtfRenderer) RenderNode(w io.Writer, node ast.Node,
	entering bool) ast.WalkStatus {

	var nextInText bool

	switch n := node.(type) {
	case *ast.Document, *ast.Link:
		nextInText = true

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

	case *ast.Text, *ast.Code, *ast.CodeBlock:
		preserveNewlines := true
		_, ok := node.(*ast.Text)
		if ok {
			preserveNewlines = false
		}
		_, codeBlock := node.(*ast.CodeBlock)
		if codeBlock {
			fmt.Fprintf(w, "\n\\par\\par\n")
		}

		leaf := node.AsLeaf()
		if leaf != nil {
			var data []byte
			for _, b := range leaf.Literal {
				switch b {
				case '\\':
					data = append(data, '\\')
					data = append(data, '\\')
				case '\n':
					if preserveNewlines {
						data = append(data, []byte("\n\\par")...)
					}
					data = append(data, ' ')
				default:
					data = append(data, b)
				}
			}
			w.Write(data)
		}
		if codeBlock {
			fmt.Fprintf(w, "\n\\par\n")
		}

	case *ast.Emph:
		if entering {
			fmt.Fprintf(w, "\\i ")
		} else {
			fmt.Fprintf(w, "\\i0 ")
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

	case *ast.Table:
		if !entering {
			nextInText = true
		}

	case *ast.TableHeader:
		if entering {
			fmt.Fprintf(w, "\\b\n")
		} else {
			fmt.Fprintf(w, "\\b0\n")
		}

	case *ast.TableBody:

	case *ast.TableRow:
		if entering {
			fmt.Fprintf(w, "\\line ")
		}
	case *ast.TableCell:
		if entering {
			fmt.Fprintf(w, "\\tab ")
		}

	default:
		fmt.Printf(" - %T %v\n", node, entering)
		nextInText = true
	}

	rtf.InText = nextInText

	return ast.GoToNext
}

// RenderHeader creates the RTF document header.
func (rtf *RtfRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	fmt.Fprintf(w, `{\rtf1\ansi\ansicpg1252\deff0\deflang1033{\fonttbl{\f0\fswiss\fcharset0 %s;}}\viewkind4\uc1\pard\ql\f0`,
		"Arial")
}

// RenderFooter creates the RTF document trailer.
func (rtf *RtfRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	fmt.Fprint(w, "}")
}

// GenerateRTF generates RTF representation of the article.
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

// ToRTF converts the argument markdown file to RTF.
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
