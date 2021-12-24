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

// Tags defines article tags.
type Tags map[string][]*Article

// NewTags creates a new tags object.
func NewTags() Tags {
	return make(map[string][]*Article)
}

// Add adds the argument tag to this tags object.
func (tags Tags) Add(tag string, article *Article) {
	tags[tag] = append(tags[tag], article)
}

// Merge adds argument tags to this tags object.
func (tags Tags) Merge(t Tags) {
	for tag, articles := range t {
		for _, article := range articles {
			tags.Add(tag, article)
		}
	}
}

// Tags returns the tags as an array of strings.
func (tags Tags) Tags() []string {
	var values []string
	for tag := range tags {
		values = append(values, tag)
	}
	sort.Strings(values)

	return values
}

// HTML returns the tags as HTML.
func (tags Tags) HTML() string {
	var result string

	values := tags.Tags()
	for idx, tag := range values {
		if idx > 0 {
			result += " "
		}
		result += fmt.Sprintf(`<a href="%s"><div class="tag">%s</div></a>`,
			TagOutputName(tag), html.EscapeString(tag))
	}
	return result
}

func TagOutputName(tag string) string {
	return fmt.Sprintf("tag-%s.html", tag)
}
