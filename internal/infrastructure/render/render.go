package render

import (
	"log"

	"github.com/kataras/blocks"
)

type Engine struct {
	TemplateEngine *blocks.Blocks
}

func New() *Engine {
	views := blocks.New("./web/views")
	if err := views.Load(); err != nil {
		log.Fatal(err)
	}

	return &Engine{
		TemplateEngine: views,
	}
}
