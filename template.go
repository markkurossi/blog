//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"text/template"
)

type Template struct {
	Tmpl *template.Template
}

func loadTemplate(name string) (tmpl *Template, err error) {
	tmpl = new(Template)
	tmpl.Tmpl, err = template.ParseGlob(name + "/*.html")
	if err != nil {
		return
	}
	tmpl.Tmpl.Option("missingkey=error")

	return
}
