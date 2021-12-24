//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"io"
	"log"
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
	Assets    map[string]os.DirEntry
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
		Assets:    make(map[string]os.DirEntry),
	}

	for _, file := range files {
		fn := file.Name()

		if file.IsDir() {
			err = tmpl.addAssets(path.Join(dir, fn))
			if err != nil {
				return nil, err
			}
		} else if strings.HasSuffix(fn, ".html") {
			t, err := template.ParseFiles(path.Join(dir, fn))
			if err != nil {
				return nil, err
			}
			t.Option("missingkey=error")

			tmpl.Templates[fn] = t
		} else {
			tmpl.Assets[path.Join(dir, fn)] = file
		}
	}

	return
}

func (tmpl *Template) addAssets(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.ReadDir(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := path.Join(dir, file.Name())
		if file.IsDir() {
			err = tmpl.addAssets(name)
			if err != nil {
				return err
			}
		} else {
			tmpl.Assets[name] = file
		}
	}

	return nil
}

// CopyAssets copies the template assets to the argument directory.
func (tmpl *Template) CopyAssets(dir string) error {
	dir = path.Clean(dir)
	for asset, assetEntry := range tmpl.Assets {

		assetInfo, err := assetEntry.Info()
		if err != nil {
			return err
		}
		output := path.Join(dir, asset[len(tmpl.Dir):])

		if isValid(assetInfo, output) {
			continue
		}

		err = os.MkdirAll(path.Dir(output), 0777)
		if err != nil {
			return err
		}

		src, err := os.Open(asset)
		if err != nil {
			return err
		}
		dst, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			assetInfo.Mode())
		if err != nil {
			src.Close()
			return err
		}
		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()
		if err != nil {
			return err
		}

		log.Printf("%s\t=> %s\n", asset, output)
	}
	return nil
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
