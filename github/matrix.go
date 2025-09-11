package github

import (
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/anchore/go-make/config"
)

var (
	MatrixSuffix = matrixSuffix()
)

func matrixSuffix() string {
	matrixJSON := config.Env("MATRIX_JSON", "")
	if matrixJSON == "" {
		return ""
	}
	values := map[string]any{}
	err := json.Unmarshal([]byte(matrixJSON), &values)
	if err != nil {
		return "-unknown"
	}
	out := ""
	// skip top-level keys
	for _, k := range slices.Sorted(maps.Keys(values)) {
		out += stringify(values[k])
	}
	return out
}

func stringify(value any) string {
	out := ""
	pat := regexp.MustCompile("[^-a-z0-9_]+")
	switch value := value.(type) {
	case string, int, float64, bool:
		out = fmt.Sprintf("-%v", value)
	case []any:
		for _, v := range value {
			out += stringify(v)
		}
	case map[string]any:
		for _, k := range slices.Sorted(maps.Keys(value)) {
			out += pat.ReplaceAllString(strings.ToLower(fmt.Sprintf("-%v%v", k, stringify(value[k]))), "-")
		}
	}
	return out
}
