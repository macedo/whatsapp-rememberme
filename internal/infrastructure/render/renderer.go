package render

import "io"

type Data map[string]interface{}

type Renderer interface {
	Render(io.Writer, Data) error
}
