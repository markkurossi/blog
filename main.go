//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/gomarkdown/markdown/parser"
)

func main() {
	template := flag.String("t", "templates/mtr", "blog template")
	flag.Parse()

	tmpl, err := loadTemplate(*template)
	if err != nil {
		log.Fatalf("failed to load template: %s", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	root := os.DirFS("articles")
	fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
		return nil
	})

	for _, arg := range flag.Args() {
		article := NewArticle(extensions)
		err = article.Parse(arg)
		if err != nil {
			log.Fatalf("process failed: %s\n", err)
		}

		err = article.Generate("out", tmpl)
		if err != nil {
			log.Fatalf("generation failed: %s\n", err)
		}
	}
}
