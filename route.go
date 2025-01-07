package reitti

import (
	"strings"
)

type route struct {
	template []component
	handler  Handler
}

// match returns true if the route matches the name. A route can contain template parameters.
//
// Example:
//
//	route{template: "repos/{owner}/{repo}/issues", handler: listIssues}
//
// The handler will be called with the context and the owner and repo as arguments.
func (r route) match(name string) ([]string, bool) {
	parts := strings.Split(name, "/")
	if len(parts) != len(r.template) {
		return nil, false
	}
	var params []string
	for i, t := range r.template {
		param, ok := t.match(parts[i])
		if !ok {
			return nil, false
		}
		if param != "" {
			params = append(params, param)
		}
	}
	return params, true
}

type component interface {
	match(string) (string, bool)
}

type literal string

var _ component = literal("repos")

func (l literal) match(s string) (string, bool) {
	return "", s == string(l)
}

type parameter string

var _ component = parameter("{owner}")

func (p parameter) match(s string) (string, bool) {
	return s, true
}
