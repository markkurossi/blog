//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"path"
	"sort"

	"github.com/gomarkdown/markdown/parser"
)

var extensions = parser.CommonExtensions | parser.AutoHeadingIDs

var tmpl *Template

func main() {
	template := flag.String("t", "templates/mtr", "blog template")
	flag.Parse()

	var err error

	tmpl, err = loadTemplate(*template)
	if err != nil {
		log.Fatalf("failed to load template: %s", err)
	}

	for _, arg := range flag.Args() {
		err = traverse(arg)
		if err != nil {
			log.Fatalf("process failed: %s\n", err)
		}
	}
	err = makeOutput()
	if err != nil {
		log.Fatalf("failed to create output: %s\n", err)
	}
}

func traverse(root string) error {
	settings, err := os.Open(path.Join(root, "settings.toml"))
	if err == nil {
		settings.Close()
		return processArticle(root)
	}
	dir, err := os.Open(root)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		err = traverse(path.Join(root, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

var articles []*Article
var tags = NewTags()
var index *Article

func processArticle(dir string) error {
	article := NewArticle(extensions)
	err := article.Parse(dir)
	if err != nil {
		return err
	}
	if article.IsIndex() {
		index = article
	} else {
		articles = append(articles, article)
	}
	tags.Merge(article.Tags)

	return nil
}

func makeOutput() error {
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Timestamp.After(articles[j].Timestamp)
	})

	var indexLinks string

	fmt.Println("Generate")
	for idx, article := range articles {
		fmt.Printf(" - %s\n", article.Name)
		if err := article.Generate("out", tmpl); err != nil {
			return err
		}
		if idx > 0 {
			indexLinks += "</br>"
		}
		indexLinks += fmt.Sprintf(`<a href="%s">%s</a>`,
			article.Name+".html", html.EscapeString(article.Title()))
		indexLinks += "\n"
	}
	if index == nil {
		return fmt.Errorf("no index")
	}
	index.Values["Links"] = indexLinks
	index.Values["Tags"] = tags.HTML()

	fmt.Printf(" - %s\n", index.Name)
	return index.Generate("out", tmpl)
}
