package agentics

import (
	"context"
	"database/sql"
	"net/http"
)

type Kind int8

const (
	PreHook Kind = iota
	PostHook
)

type Context struct {
	Bag    *Bag[any]
	Memory Memory
	HTTP   *http.Client
	DB     *sql.DB
}

type Func func(ctx context.Context, bag *Bag[any], mem Memory) error

var hookRegistry = map[string]Func{}

func RegisterHook(name string, fn Func) {
	hookRegistry[name] = fn
}

func getHook(name string) (Func, bool) {
	fn, ok := hookRegistry[name]
	return fn, ok
}
