//
// Copyright (c) 2021-2024 Markku Rossi
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
	Values       Values
	Extensions   parser.Extensions
	Settings     Settings
	FolderSuffix string
	Name         string
	Tags         Tags
	Timestamp    time.Time
	Published    bool
	Site         bool
	Pagenum      int

	Assets *Assets
}

// Settings define the article settings.
type Settings struct {
	Article struct {
		Title     string
		Tags      []string
		Type      string
		Published time.Time
	} `toml:"article"`
	Meta struct {
		Title       string
		Description string
	} `toml:"meta"`
}

// NewArticle creates a new article with the Markdown extensions.
func NewArticle(extensions parser.Extensions) *Article {
	return &Article{
		Values:     NewValues(),
		Extensions: extensions,
		Tags:       NewTags(),
	}
}

// NewSiteArticle creates a new site article with the Markdown
// extensions.
func NewSiteArticle(extensions parser.Extensions, name string) *Article {
	article := &Article{
		Values:     NewValues(),
		Extensions: extensions,
		Name:       name,
		Tags:       NewTags(),
		Site:       true,
	}

	dir := path.Dir(name)
	if dir == "." {
		dir = ""
	}

	article.Values["OutputDir"] = dir

	return article
}

// Type returns the article type as template well-known name string.
func (article *Article) Type() string {
	switch article.Settings.Article.Type {
	case "presentation":
		return TmplPresentation
	case "article", "":
		return TmplArticle
	default:
		panic(fmt.Sprintf("unknown article type: %v",
			article.Settings.Article.Type))
	}
}

// Parse parses article data from the argument directory.
func (article *Article) Parse(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	Verbose(" - %s\n", dir)
	article.Assets = NewAssets(dir)

	// Tags from the path.
	parts := strings.Split(path.Clean(dir), "/")
	for i := 1; i < len(parts)-1; i++ {
		article.Tags.Add(parts[i], article)
	}

	files, err := f.ReadDir(0)
	if err != nil {
		return err
	}
	// Collect sources and process them once the source directory is
	// processed. We must see settings.toml before processing the
	// content.
	var sources []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "~") {
			continue
		}
		Verbose("   - %s\n", file.Name())

		fi, err := file.Info()
		if err != nil {
			return err
		}
		if fi.ModTime().After(article.Timestamp) {
			article.Timestamp = fi.ModTime()
		}

		if strings.HasSuffix(file.Name(), ".md") {
			sources = append(sources, file.Name())
		} else if file.Name() == "settings.toml" {
			err = article.readSettings(dir, file.Name())
			if err != nil {
				return err
			}
		} else {
			// Save asset files.
			err = article.Assets.Add(path.Join(dir, file.Name()), file)
			if err != nil {
				return err
			}
		}
	}
	// Process sources.
	for _, source := range sources {
		err = article.processFile(dir, source)
		if err != nil {
			return err
		}
	}
	article.Name = path.Base(dir)
	switch article.Name {
	case ".", "/":
		return fmt.Errorf("invalid input name: %s", dir)
	}

	if article.IsIndex() {
		article.Values.SetRaw(ValOutputDir, "")
	} else {
		article.Values.SetRaw(ValOutputDir, "../")
	}

	// Create tags value.
	for _, tag := range article.Settings.Article.Tags {
		article.Tags.Add(tag, article)
	}
	article.Values.SetRaw(ValTags,
		article.Tags.HTML(article.Values[ValOutputDir]))

	ts := article.Settings.Article.Published
	if ts.IsZero() {
		ts = time.Now()
		article.Values.Set(ValDraft, "Draft")
		article.Values.Set(ValPublished, "Unpublished Draft")
	} else {
		article.Values.Set(ValDraft, "")
		article.Values.Set(ValPublished, ts.Format(time.UnixDate))
		article.Published = true

		// Published articles get their timestamp from the published
		// time.
		article.Timestamp = ts
	}

	article.Values.Set(ValLinks, "")
	article.Values.Set(ValYear, strconv.Itoa(ts.Year()))

	// Meta.
	return article.createMeta()
}

func (article *Article) createMeta() error {
	metaTitle := article.Settings.Meta.Title
	if len(metaTitle) == 0 {
		metaTitle = article.Title()
	}
	if len(metaTitle) > MaxMetaTitleLen {
		return fmt.Errorf("meta title too long: %d > %d",
			len(metaTitle), MaxMetaTitleLen)
	}
	metaDesc := article.Settings.Meta.Description
	if len(metaDesc) > MaxMetaDescriptionLen {
		return fmt.Errorf("meta description too long: %d > %d",
			len(metaDesc), MaxMetaDescriptionLen)
	}
	article.Values.Set(ValMetaTitle, metaTitle)
	article.Values.Set(ValMetaDescription, metaDesc)

	return nil
}

// ParseSiteFile parses a site input file.
func (article *Article) ParseSiteFile(file, section string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	sectionName := strings.Title(section)
	sectionData := string(article.format(data))

	article.Values.SetRaw(sectionName, sectionData)
	return nil
}

// ParseSiteFileSettings parses site file settings.
func (article *Article) ParseSiteFileSettings(dir, file string) error {
	err := article.readSettings(dir, file)
	if err != nil {
		return err
	}

	// Meta.
	return article.createMeta()
}

// IsIndex tests if this article is the blog main index article.
func (article *Article) IsIndex() bool {
	return article.Name == "index"
}

// Title returns the article title.
func (article *Article) Title() string {
	title, ok := article.Values[ValTitle]
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

	article.Values.SetRaw(sectionName, sectionData)
	return nil
}

func (article *Article) readSettings(dir, file string) error {
	_, err := toml.DecodeFile(path.Join(dir, file), &article.Settings)
	if err != nil {
		return err
	}
	article.Values.Set(ValTitle, article.Settings.Article.Title)
	return nil
}

// Generate generates article HTML to the argument directory, using
// the specified output template.
func (article *Article) Generate(dir string, tmpl *Template) error {
	filename := path.Join(dir, article.OutputName())
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}

	// Copy asset files.
	if article.Assets != nil {
		err = article.Assets.Copy(path.Join(dir, article.OutputFolder()))
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if !article.Site && article.Name == "index" {
		return tmpl.Templates[TmplIndex].Execute(f, article.Values)
	}
	return tmpl.Templates[article.Type()].Execute(f, article.Values)
}

// OutputFolder returns the article output folder name.
func (article *Article) OutputFolder() string {
	if article.Site {
		return ""
	}
	return article.Timestamp.Format("2006-01-02") + article.FolderSuffix
}

// SetFolderSuffix sets the article output directory suffix.
func (article *Article) SetFolderSuffix(suffix string) {
	article.FolderSuffix = suffix
}

// OutputName returns the article HTML output name.
func (article *Article) OutputName() string {
	return article.outputName(".html")
}

// RTFOutputName returns the article RTF output name.
func (article *Article) RTFOutputName() string {
	return article.outputName(".rtf")
}

func (article *Article) outputName(suffix string) string {
	filename := article.Name + suffix
	if article.IsIndex() {
		return filename
	}
	return path.Join(article.OutputFolder(), filename)
}

// Link returns an HTML link to this article.
func (article *Article) Link() string {
	link := fmt.Sprintf(`<a href="%s">%s`, article.OutputName(),
		article.Title())

	switch article.Type() {
	case TmplPresentation:
		link += " &#x1f4fd;"
	}
	if !article.Published {
		link += " [draft]"
	}
	return link + "</a>"
}
