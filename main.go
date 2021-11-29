//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"log"

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

	for _, arg := range flag.Args() {
		article := NewArticle(extensions)
		err = article.Parse(arg)
		if err != nil {
			log.Fatalf("process failed: %s\n", err)
		}

		err = article.Generate("out/index.html", tmpl)
		if err != nil {
			log.Fatalf("generation failed: %s\n", err)
		}
	}
}
