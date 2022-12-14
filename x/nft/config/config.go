package config

import (
	"regexp"
)

var (
	CollectionDenomRegExp    = "^[a-zA-Z0-9_-]{1,255}$"
	CollectionDenomValidator = regexp.MustCompile(CollectionDenomRegExp)
)
