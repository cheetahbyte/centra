package helper

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func ParseQueryParams(r *http.Request) map[string][]string {
	q := r.URL.Query()
	result := make(map[string][]string)

	for key, values := range q {
		for _, v := range values {
			split := strings.Split(v, ",")
			for _, s := range split {
				result[key] = append(result[key], strings.TrimSpace(s))
			}
		}
	}

	return result
}

func matchesValue(metaValue any, allowed []string) bool {
	switch v := metaValue.(type) {
	case string:
		ms := strings.ToLower(v)
		for _, a := range allowed {
			if strings.Contains(ms, strings.ToLower(a)) {
				return true
			}
		}
		return false

	case []string:
		return slices.ContainsFunc(v, func(tag string) bool {
			ts := strings.ToLower(tag)
			for _, a := range allowed {
				if strings.Contains(ts, strings.ToLower(a)) {
					return true
				}
			}
			return false
		})

	case bool, int, int64, float64:
		s := fmt.Sprintf("%v", v)
		ss := strings.ToLower(s)
		for _, a := range allowed {
			if strings.Contains(ss, strings.ToLower(a)) {
				return true
			}
		}
		return false

	default:
		s := fmt.Sprint(v)
		ss := strings.ToLower(s)
		for _, a := range allowed {
			if strings.Contains(ss, strings.ToLower(a)) {
				return true
			}
		}
		return false
	}
}

func MatchesQuery(meta map[string]any, query map[string][]string) bool {
	for key, allowed := range query {
		val, ok := meta[key]
		if !ok || !matchesValue(val, allowed) {
			return false
		}
	}
	return true
}
