//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"path"
	"text/template"
)

type Template struct {
	Index   *template.Template
	Article *template.Template
}

func loadTemplate(name string) (tmpl *Template, err error) {
	tmpl = new(Template)
	tmpl.Index, err = template.ParseFiles(path.Join(name, "index.html"))
	if err != nil {
		return
	}
	tmpl.Article, err = template.ParseFiles(path.Join(name, "article.html"))
	if err != nil {
		return
	}

	tmpl.Index.Option("missingkey=error")
	tmpl.Article.Option("missingkey=error")

	return
}
