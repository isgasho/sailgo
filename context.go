package sailgo

import (
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
}

func (ctx *Context) Getparam(name string) string {
	if value, ok := ctx.Params[name]; ok {
		return value
	}
	return ""
}
