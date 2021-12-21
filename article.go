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
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown/parser"
)

// Article implements a blog article.
type Article struct {
	Values     map[string]string
	Extensions parser.Extensions
	Settings   Settings
	Name       string
	Tags       Tags
	Timestamp  time.Time
}

// Settings define the article settings.
type Settings struct {
	Article struct {
		Title string
		Tags  []string
	} `toml:"article"`
}

// NewArticle creates a new article with the Markdown extensions.
func NewArticle(extensions parser.Extensions) *Article {
	return &Article{
		Values:     make(map[string]string),
		Extensions: extensions,
		Tags:       NewTags(),
	}
}

// Parse parses article data from the argument directory.
func (article *Article) Parse(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf(" - %s\n", dir)

	// Tags from the path.
	parts := strings.Split(path.Clean(dir), "/")
	for i := 1; i < len(parts)-1; i++ {
		article.Tags.Add(parts[i])
	}

	files, err := f.ReadDir(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "~") {
			continue
		}
		fmt.Printf("   - %s\n", file.Name())
		fi, err := file.Info()
		if err != nil {
			return err
		}
		if fi.ModTime().After(article.Timestamp) {
			article.Timestamp = fi.ModTime()
		}

		if strings.HasSuffix(file.Name(), ".md") {
			err = article.processFile(dir, file.Name())
			if err != nil {
				return err
			}
		} else if file.Name() == "settings.toml" {
			err = article.readSettings(dir, file.Name())
			if err != nil {
				return err
			}
		}
	}
	article.Name = path.Base(dir)
	switch article.Name {
	case ".", "/":
		return fmt.Errorf("invalid input name: %s", dir)
	}

	// Create tags value.
	for _, tag := range article.Settings.Article.Tags {
		article.Tags.Add(tag)
	}
	article.Values["Tags"] = article.Tags.HTML()

	// XXX Published timestamp from settings.

	article.Values["Links"] = ""
	article.Values["Year"] = strconv.Itoa(time.Now().Year())

	return nil
}

// IsIndex tests if this article is the blog main index article.
func (article *Article) IsIndex() bool {
	return article.Name == "index"
}

// Title returns the article title.
func (article *Article) Title() string {
	title, ok := article.Values["Title"]
	if ok {
		return title
	}
	return article.Name
}

func (article *Article) processFile(dir, file string) error {
	f, err := os.Open(path.Join(dir, file))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	parts := strings.Split(file[:len(file)-3], "-")
	for idx, part := range parts {
		parts[idx] = strings.Title(part)
	}
	sectionName := strings.Join(parts, "")
	sectionData := string(article.format(data))

	article.Values[sectionName] = sectionData
	return nil
}

func (article *Article) readSettings(dir, file string) error {
	_, err := toml.DecodeFile(path.Join(dir, file), &article.Settings)
	if err != nil {
		return err
	}
	article.Values["Title"] = article.Settings.Article.Title
	return nil
}

// Generate generates article HTML to the argument directory, using
// the specified output template.
func (article *Article) Generate(dir string, tmpl *Template) error {

	file := path.Join(dir, article.Name+".html")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if article.Name == "index" {
		return tmpl.Templates[TmplIndex].Execute(f, article.Values)
	}
	return tmpl.Templates[TmplArticle].Execute(f, article.Values)
}
