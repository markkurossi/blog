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
type Tags map[string]string

// NewTags creates a new tags object.
func NewTags() Tags {
	return make(map[string]string)
}

// Add adds the argument tag to this tags object.
func (tags Tags) Add(tag string) {
	tags[tag] = tag
}

// Merge adds argument tags to this tags object.
func (tags Tags) Merge(t Tags) {
	for k := range t {
		tags.Add(k)
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
	for idx, v := range values {
		if idx > 0 {
			result += " "
		}
		result += fmt.Sprintf(`<div class="tag">%s</div>`, html.EscapeString(v))
	}
	return result
}
