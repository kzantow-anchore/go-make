package gomake

import (
	"github.com/anchore/go-make/lang"
	"github.com/anchore/go-make/log"
)

func List[T any](items ...T) []T {
	return lang.List(items...)
}

var (
	Log = log.Log
)
