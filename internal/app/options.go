package app

import (
	"context"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/macedo/whatsapp-rememberme/pkg/defaults"
	"github.com/macedo/whatsapp-rememberme/pkg/env"
)

type Options struct {
	Addr         string
	Env          string
	SessionStore *sessions.CookieStore

	Context context.Context
	cancel  context.CancelFunc
}

func options_with_default(opts Options) Options {
	opts.Env = defaults.String(opts.Env, env.Get("APP_ENV", "development"))

	addr := "0.0.0.0"
	if opts.Env == "development" {
		addr = "127.0.0.1"
	}
	opts.Addr = defaults.String(
		opts.Addr,
		fmt.Sprintf("%s:%s", env.Get("ADDR", addr), env.Get("PORT", "3000")),
	)

	if opts.SessionStore == nil {
		secret := env.Get("SECRET", "")
		if secret == "" && (opts.Env == "development" || opts.Env == "test") {
			secret = "secret"
		}

		if secret == "" {
			fmt.Println("Unless you set SECRET env variable, your session storage is not protected!")
		}

		cookieStore := sessions.NewCookieStore([]byte(secret))
		cookieStore.Options.HttpOnly = true
		if opts.Env == "production" {
			cookieStore.Options.Secure = true
		}
		opts.SessionStore = cookieStore
	}

	opts.Context, opts.cancel = context.WithCancel(context.Background())

	return opts
}
