package app

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/render"
)

type Context struct {
	context.Context
	data     *sync.Map
	flash    *Flash
	request  *http.Request
	response http.ResponseWriter
}

func (c *Context) Data() map[string]interface{} {
	m := map[string]interface{}{}

	if c.data == nil {
		return m
	}

	c.data.Range(func(k, v interface{}) bool {
		s, ok := k.(string)
		if !ok {
			return false
		}
		m[s] = v
		return true
	})

	return m
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.Response(), c.Request(), url, http.StatusSeeOther)
}

func (c *Context) Render(status int, rr render.Renderer) error {
	var err error

	if rr == nil {
		c.Response().WriteHeader(status)
		return nil
	}

	data := c.Data()
	data["flash"] = c.flash.data

	output := &bytes.Buffer{}

	err = rr.Render(output, data)
	if err != nil {
		return err
	}

	c.Response().WriteHeader(status)
	_, err = io.Copy(c.Response(), output)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() http.ResponseWriter {
	return c.response
}

func (a *App) newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Context: r.Context(),
		data:    &sync.Map{},
		flash: &Flash{
			data: map[string]string{},
		},
		request:  r,
		response: w,
	}
}
