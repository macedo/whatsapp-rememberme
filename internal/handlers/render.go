package handlers

import (
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/render"
)

var rr *render.Engine

func init() {
	rr = render.New()
}
