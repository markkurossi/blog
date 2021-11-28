//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	template := flag.String("t", "templates/mtr", "blog template")
	flag.Parse()

	tmpl, err := loadTemplate(*template)
	if err != nil {
		log.Fatalf("failed to load template: %s", err)
	}

	values := map[string]string{
		"Title":         "Title",
		"ColumnLeft":    "Left column",
		"ColumnArticle": "Article column",
	}

	err = tmpl.Tmpl.Execute(os.Stdout, values)
	if err != nil {
		log.Fatalf("failed to execute tempate: %s", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	for _, arg := range flag.Args() {
		err = processFile(arg, parser)
		if err != nil {
			log.Fatalf("process failed: %s\n", err)
		}
	}

}

func processFile(name string, parser *parser.Parser) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	html := markdown.ToHTML(data, parser, nil)

	fmt.Printf("HTML: %s\n", html)
	return nil
}
