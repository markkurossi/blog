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
	"strconv"
	"time"

	"github.com/gomarkdown/markdown/parser"
)

var (
	program     = path.Base(os.Args[0])
	extensions  = parser.CommonExtensions | parser.AutoHeadingIDs
	tmpl        *Template
	flagVerbose bool
	flagDraft   bool
)

func main() {
	log.SetFlags(0)

	template := flag.String("t", "templates/mtr", "blog template")
	out := flag.String("o", "", "output directory")

	flag.BoolVar(&flagVerbose, "v", false, "verbose output")
	flag.BoolVar(&flagDraft, "draft", false, "process draft articles")

	flag.Parse()

	if len(*out) == 0 {
		log.Fatalf("%s: output directory not specified", program)
	}

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
	err = makeOutput(*out)
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
		if article.Published || flagDraft {
			articles = append(articles, article)
			tags.Merge(article.Tags)
		}
	}

	return nil
}

func makeOutput(out string) error {
	err := os.MkdirAll(out, 0777)
	if err != nil {
		return err
	}

	err = tmpl.CopyAssets(out)
	if err != nil {
		return err
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Timestamp.After(articles[j].Timestamp)
	})

	var indexLinks string

	Verbose("Generate")
	for idx, article := range articles {
		Verbose(" - %s\n", article.Name)
		if err := article.Generate(out, tmpl); err != nil {
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

	Verbose(" - %s\n", index.Name)
	if err := index.Generate(out, tmpl); err != nil {
		return err
	}

	// Tag indices.
	if tmpl.Templates[TmplTag] == nil {
		return fmt.Errorf("tag template %s not defined", TmplTag)
	}
	for _, tag := range tags.Tags() {
		fmt.Printf(" - %s\n", tag)
		if err := makeTagOutput(out, tag, tags[tag]); err != nil {
			return err
		}
	}

	return nil
}

func makeTagOutput(out, tag string, articles []*Article) error {
	f, err := os.Create(path.Join(out, TagOutputName(tag)))
	if err != nil {
		return err
	}
	defer f.Close()

	values := make(map[string]string)
	value := "<ul>"
	for _, article := range articles {
		value += "\n  <li>"
		value += article.Link()
	}
	value += "\n</ul>\n"

	values["Title"] = fmt.Sprintf("%s - Blog Category", tag)
	values["H1"] = fmt.Sprintf("Blog Category '%s'", tag)
	values["Tags"] = tags.HTML()
	values["TagLinks"] = value
	values["Year"] = strconv.Itoa(time.Now().Year())

	return tmpl.Templates[TmplTag].Execute(f, values)
}
