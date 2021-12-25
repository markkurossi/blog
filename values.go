//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"html"
	"strconv"
	"time"
)

// Template variables.
const (
	ValTitle           = "Title"
	ValH1              = "H1"
	ValTags            = "Tags"
	ValTagLinks        = "TagLinks"
	ValDraft           = "Draft"
	ValPublished       = "Published"
	ValLinks           = "Links"
	ValYear            = "Year"
	ValMetaTitle       = "MetaTitle"
	ValMetaDescription = "MetaDescription"
)

// Values define template variables and their values.
type Values map[string]string

// NewValues creates a new values object.
func NewValues() Values {
	return map[string]string{
		"Year": strconv.Itoa(time.Now().Year()),
	}
}

// Set sets value for a key. The value is HTML escaped.
func (values Values) Set(k, v string) {
	values[k] = html.EscapeString(v)
}

// SetRaw sets the raw value for a key.
func (values Values) SetRaw(k, v string) {
	values[k] = v
}
