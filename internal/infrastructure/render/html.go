package render

import (
	"io"
	"strings"
)

func (e *Engine) HTML(name string) Renderer {
	return &templateRenderer{
		Engine: e,
		name:   name,
	}
}

type templateRenderer struct {
	*Engine
	name string
}

func (r *templateRenderer) Render(w io.Writer, data Data) error {
	s := strings.Split(r.name, "/")
	tmplName := s[0]

	layoutName := ""

	if len(s) > 1 {
		layoutName = s[1]
	}

	return r.TemplateEngine.ExecuteTemplate(w, tmplName, layoutName, data)
}
