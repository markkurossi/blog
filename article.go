//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"html"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown/parser"
)

type Article struct {
	Values     map[string]string
	Extensions parser.Extensions
	Settings   Settings
	Name       string
	Tags       map[string]string
}

type Settings struct {
	Article SettingsArticle `toml:"article"`
}

type SettingsArticle struct {
	Title string
	Tags  []string
}

func NewArticle(extensions parser.Extensions) *Article {
	return &Article{
		Values:     make(map[string]string),
		Extensions: extensions,
		Tags:       make(map[string]string),
	}
}

func (article *Article) Parse(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	// Tags from the path.
	parts := strings.Split(path.Clean(dir), "/")
	for i := 1; i < len(parts)-1; i++ {
		article.Tags[parts[i]] = parts[i]
	}

	files, err := f.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file, "~") {
			continue
		}
		fmt.Printf(" - %s\n", file)
		if strings.HasSuffix(file, ".md") {
			err = article.processFile(dir, file)
			if err != nil {
				return err
			}
		} else if file == "settings.toml" {
			err = article.readSettings(dir, file)
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
		article.Tags[tag] = tag
	}
	var tags []string
	for tag := range article.Tags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	var tagsValue string
	for idx, tag := range tags {
		if idx > 0 {
			tagsValue += " "
		}
		tagsValue += fmt.Sprintf(`<div class="tag">%s</div>`,
			html.EscapeString(tag))
	}
	article.Values["Tags"] = tagsValue

	return nil
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

func (article *Article) Generate(dir string, tmpl *Template) error {

	file := path.Join(dir, article.Name+".html")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Tmpl.Execute(f, article.Values)
}
