// Package reitti provides a simple router for matching routes with handlers.
// It allows adding routes with template parameters and matching them against
// incoming requests. The router supports context and can handle functions with
// varying signatures as handlers.
package reitti

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

var ErrNoRoute = fmt.Errorf("no route")

type handler func(ctx context.Context, args ...string) (any, error)

type Router struct {
	routes []route
}

func (r *Router) Add(template string, handler any) {
	r.routes = append(r.routes, route{
		template: compile(template),
		handler:  wrap(handler),
	})
}

// wrap converts a function into a route.handler.
func wrap(f any) func(context.Context, ...string) (any, error) {
	ty := reflect.TypeOf(f)
	validateHandler(ty)
	val := reflect.ValueOf(f)
	h, isHandler := f.(func(context.Context, ...string) (any, error))
	return func(ctx context.Context, args ...string) (any, error) {
		if isHandler {
			return h(ctx, args...)
		}
		return reflectCall(ctx, val, ty, args)
	}
}

func reflectCall(ctx context.Context, val reflect.Value, ty reflect.Type, args []string) (any, error) {
	n := ty.NumIn()
	if n != len(args)+1 {
		return nil, fmt.Errorf("expected %d args, got %d", n-1, len(args))
	}
	in := make([]reflect.Value, n)
	in[0] = reflect.ValueOf(ctx)
	for i := 1; i < n; i++ {
		in[i] = reflect.ValueOf(args[i-1])
	}
	values := val.Call(in)
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 results, got %d", len(values))
	}
	res, err := values[0], values[1]
	if !err.Type().Implements(reflect.TypeFor[error]()) {
		return nil, fmt.Errorf("expected an error, got %v", err.Type())
	}
	e := err.Interface()
	if e != nil {
		return nil, e.(error)
	}
	return res.Interface(), nil
}

func validateHandler(ty reflect.Type) {
	if ty.Kind() != reflect.Func {
		panic(fmt.Errorf("expected a function, got %v", ty))
	}
	if ty.NumOut() != 2 {
		panic(fmt.Errorf("expected 2 results, got %d", ty.NumOut()))
	}
	secondOut := ty.Out(1)
	if secondOut != reflect.TypeFor[error]() {
		panic(fmt.Errorf("expected an error as the second result, got %v", secondOut))
	}
	if ty.NumIn() < 1 {
		panic(fmt.Errorf("expected at least 1 argument, got %d", ty.NumIn()))
	}
	firstArg := ty.In(0)
	if firstArg != reflect.TypeFor[context.Context]() {
		panic(fmt.Errorf("expected a context.Context as the first argument, got %v", firstArg))
	}
}

func (r *Router) Match(name string) (func(ctx context.Context) (any, error), error) {
	for _, route := range r.routes {
		if params, ok := route.match(name); ok {
			return func(ctx context.Context) (any, error) {
				return route.handler(ctx, params...)
			}, nil
		}
	}
	return nil, ErrNoRoute
}

func compile(template string) []component {
	var components []component
	for _, part := range strings.Split(template, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			components = append(components, parameter(part[1:len(part)-1]))
		} else {
			components = append(components, literal(part))
		}
	}
	return components
}
