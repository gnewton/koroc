package main

import (
	"github.com/microcosm-cc/bluemonday"
)

var policy = bluemonday.StripTagsPolicy()

func strip(xml string) string {
	return policy.Sanitize(xml)
}
