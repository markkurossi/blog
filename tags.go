//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"html"
	"sort"
)

type Tags map[string]string

func NewTags() Tags {
	return make(map[string]string)
}

func (tags Tags) Add(tag string) {
	tags[tag] = tag
}

func (tags Tags) Merge(t Tags) {
	for k := range t {
		tags.Add(k)
	}
}

func (tags Tags) Tags() []string {
	var values []string
	for tag := range tags {
		values = append(values, tag)
	}
	sort.Strings(values)

	return values
}

func (tags Tags) HTML() string {
	var result string

	values := tags.Tags()
	for idx, v := range values {
		if idx > 0 {
			result += " "
		}
		result += fmt.Sprintf(`<div class="tag">%s</div>`, html.EscapeString(v))
	}
	return result
}
