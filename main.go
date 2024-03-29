//
// Copyright (c) 2021-2024 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown/parser"
)

var (
	program     = path.Base(os.Args[0])
	extensions  = parser.CommonExtensions | parser.AutoHeadingIDs
	tmpl        *Template
	flagVerbose bool
	flagDraft   bool
	flagRTF     bool
	flagSite    bool
	flagLibrary string
)

func main() {
	log.SetFlags(0)

	template := flag.String("t", "templates/mtr", "blog template")
	out := flag.String("o", "", "output directory")

	flag.BoolVar(&flagVerbose, "v", false, "verbose output")
	flag.BoolVar(&flagDraft, "draft", false, "process draft articles")
	flag.BoolVar(&flagRTF, "rtf", false, "generate RTF output")
	flag.BoolVar(&flagSite, "site", false, "site mode")
	flag.StringVar(&flagLibrary, "lib", ".", "asset library path")

	flag.Parse()

	if len(*out) == 0 {
		log.Fatalf("%s: output directory not specified", program)
	}

	var err error

	tmpl, err = loadTemplate(path.Join(flagLibrary, *template))
	if err != nil {
		log.Fatalf("failed to load template: %s", err)
	}

	for _, arg := range flag.Args() {
		if flagSite {
			assets := NewAssets(arg)
			siteAssets = append(siteAssets, assets)
			err = traverseSite(assets, arg, arg)
		} else {
			err = traverse(arg)
		}
		if err != nil {
			log.Fatalf("process failed: %s\n", err)
		}
	}
	if flagSite {
		for _, a := range siteArticles {
			articles = append(articles, a)
		}
	}

	err = makeOutput(*out)
	if err != nil {
		log.Fatalf("failed to create output: %s\n", err)
	}

	if flagRTF {
		err = makeRTF(*out)
		if err != nil {
			log.Fatalf("failed to create RTF output: %s\n", err)
		}
	}
}

var siteArticles = make(map[string]*Article)
var siteAssets []*Assets

func traverseSite(assets *Assets, root, dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	entries, err := d.ReadDir(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = traverseSite(assets, root, path.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), "~") {
			// Skip Emacs backup files.
		} else if strings.HasSuffix(entry.Name(), ".toml") {
			name := entry.Name()
			name = name[0 : len(name)-5]

			article := getSiteArticle(path.Join(dir, name)[len(root):], name)
			err = article.ParseSiteFileSettings(dir, entry.Name())
			if err != nil {
				return err
			}

		} else if strings.HasSuffix(entry.Name(), ".md") {
			name := entry.Name()
			name = name[0 : len(name)-3]

			parts := strings.Split(name, ",")
			if len(parts) != 2 {
				return fmt.Errorf("invalid input file '%s', expected 2 parts",
					name)
			}
			article := getSiteArticle(path.Join(dir, parts[0])[len(root):],
				parts[0])
			err = article.ParseSiteFile(path.Join(dir, entry.Name()), parts[1])
			if err != nil {
				return err
			}
		} else {
			assets.Add(path.Join(dir, entry.Name()), entry)
		}
	}
	return nil
}

func getSiteArticle(fullName, name string) *Article {
	article, ok := siteArticles[name]
	if !ok {
		article = NewSiteArticle(extensions, fullName)
		siteArticles[name] = article
	}
	return article
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
		return fmt.Errorf("%s: %s", dir, err.Error())
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

	for _, asset := range siteAssets {
		err = asset.Copy(out)
		if err != nil {
			return err
		}
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Timestamp.After(articles[j].Timestamp)
	})

	// Assign indices for articles, published on the same day.
	byDay := make(map[string][]*Article)
	for _, article := range articles {
		folder := article.OutputFolder()
		byDay[folder] = append(byDay[folder], article)
	}
	for _, group := range byDay {
		if len(group) > 1 {
			for idx, article := range group {
				article.SetFolderSuffix(fmt.Sprintf("-%d", idx))
			}
		}
	}

	var indexLinks string

	Verbose("Generate\n")
	for idx, article := range articles {
		Verbose(" - %s\n", article.OutputName())
		if err := article.Generate(out, tmpl); err != nil {
			return err
		}
		if idx > 0 {
			indexLinks += "</br>"
		}
		indexLinks += article.Link() + "\n"
	}
	if flagSite {
		return nil
	}
	if index == nil {
		return fmt.Errorf("no index")
	}
	index.Values.SetRaw(ValLinks, indexLinks)
	index.Values.SetRaw(ValTags, tags.HTML(""))

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
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Timestamp.After(articles[j].Timestamp)
	})

	f, err := os.Create(path.Join(out, TagOutputName(tag)))
	if err != nil {
		return err
	}
	defer f.Close()

	values := NewValues()
	value := "<ul>"
	for _, article := range articles {
		value += "\n  <li>"
		value += article.Link()
	}
	value += "\n</ul>\n"

	h1 := fmt.Sprintf("Tag Category '%s'", tag)

	values.Set(ValTitle, fmt.Sprintf("%s - Tag Category", tag))
	values.Set(ValH1, h1)
	values.SetRaw(ValTags, tags.HTML(""))
	values.SetRaw(ValTagLinks, value)

	values.Set(ValMetaTitle, h1)
	values.Set(ValMetaDescription, fmt.Sprintf("Articles in category '%s'",
		tag))

	return tmpl.Templates[TmplTag].Execute(f, values)
}

func makeRTF(out string) error {
	for _, article := range articles {
		err := article.GenerateRTF(out)
		if err != nil {
			return err
		}
	}
	return nil
}
