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
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Article struct {
	Values     map[string]string
	Extensions parser.Extensions
	Settings   Settings
}

type Settings struct {
	Article SettingsArticle `toml:"article"`
}

type SettingsArticle struct {
	Title string
}

func NewArticle(extensions parser.Extensions) *Article {
	return &Article{
		Values:     make(map[string]string),
		Extensions: extensions,
	}
}

func (article *Article) Parse(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

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

	parser := parser.NewWithExtensions(article.Extensions)

	parts := strings.Split(file[:len(file)-3], "-")
	for idx, part := range parts {
		parts[idx] = strings.Title(part)
	}
	sectionName := strings.Join(parts, "")
	sectionData := string(markdown.ToHTML(data, parser, nil))

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

func (article *Article) Generate(file string,
	tmpl *Template) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Tmpl.Execute(f, article.Values)
}
