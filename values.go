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

type Values map[string]string

func NewValues() Values {
	return map[string]string{
		"Year": strconv.Itoa(time.Now().Year()),
	}
}

func (values Values) Set(k, v string) {
	values[k] = html.EscapeString(v)
}

func (values Values) SetRaw(k, v string) {
	values[k] = v
}
