//
// Copyright (c) 2021-2022 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"os"
	"path"
	"strings"
	"text/template"
)

// Well-known template file names.
const (
	TmplIndex   = "index.html"
	TmplArticle = "article.html"
	TmplTag     = "tag.html"
)

// Template defines blog output template.
type Template struct {
	Dir       string
	Templates map[string]*template.Template
	Assets    *Assets
}

func loadTemplate(dir string) (tmpl *Template, err error) {
	dir = path.Clean(dir)
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.ReadDir(0)
	if err != nil {
		return nil, err
	}

	tmpl = &Template{
		Dir:       dir,
		Templates: make(map[string]*template.Template),
		Assets:    NewAssets(dir),
	}

	for _, file := range files {
		fn := file.Name()

		if file.IsDir() {
			err = tmpl.Assets.AddDir(path.Join(dir, fn))
			if err != nil {
				return nil, err
			}
		} else if strings.HasSuffix(fn, "~") {
		} else if strings.HasSuffix(fn, ".html") {
			t, err := template.ParseFiles(path.Join(dir, fn))
			if err != nil {
				return nil, err
			}
			t.Option("missingkey=error")

			tmpl.Templates[fn] = t
		} else {
			tmpl.Assets.Add(path.Join(dir, fn), file)
		}
	}

	return
}

// CopyAssets copies the template assets to the argument directory.
func (tmpl *Template) CopyAssets(dir string) error {
	return tmpl.Assets.Copy(dir)
}

func isValid(srcInfo os.FileInfo, file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return false
	}
	if info.Size() != srcInfo.Size() {
		return false
	}
	if info.ModTime().Before(srcInfo.ModTime()) {
		return false
	}

	return true
}
